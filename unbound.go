package unbound

import (
	"fmt"
	"strconv"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
	"github.com/miekg/unbound"
	"golang.org/x/net/context"
)

var log = clog.NewWithPlugin("unbound")

// Unbound is a plugin that resolves requests using libunbound.
type Unbound struct {
	u *unbound.Unbound
	t *unbound.Unbound

	from   []string
	except []string

	Next plugin.Handler
}

// options for unbound, see unbound.conf(5).
var options = map[string]string{
	"msg-cache-size":   "0",
	"rrset-cache-size": "0",
}

// New returns a pointer to an initialzed Unbound.
func New() *Unbound {
	udp := unbound.New()
	tcp := unbound.New()
	tcp.SetOption("tcp-upstream:", "yes")

	u := &Unbound{u: udp, t: tcp}

	for k, v := range options {
		if err := u.setOption(k, v); err != nil {
			log.Warningf("Could not set option: %s", err)
		}
	}

	return u
}

// Stop stops unbound and cleans up the memory used.
func (u *Unbound) Stop() error {
	u.u.Destroy()
	u.t.Destroy()
	return nil
}

// setOption sets option k to value v in u.
func (u *Unbound) setOption(k, v string) error {
	// Add ":" as unbound expects it
	k += ":"
	// Set for both udp and tcp handlers, return the error from the latter.
	u.u.SetOption(k, v)
	err := u.t.SetOption(k, v)
	if err != nil {
		return fmt.Errorf("failed to set option %q with value %q: %s", k, v, err)
	}
	return nil
}

// ServeDNS implements the plugin.Handler interface.
func (u *Unbound) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	if !u.match(state) {
		return plugin.NextOrFailure(u.Name(), u.Next, ctx, w, r)
	}

	var (
		res *unbound.Result
		err error
	)
	switch state.Proto() {
	case "tcp":
		res, err = u.t.Resolve(state.QName(), state.QType(), state.QClass())
	case "udp":
		res, err = u.u.Resolve(state.QName(), state.QType(), state.QClass())
	}

	rcode := dns.RcodeServerFailure
	if err == nil {
		rcode = res.AnswerPacket.Rcode
	}
	rc, ok := dns.RcodeToString[rcode]
	if !ok {
		rc = strconv.Itoa(rcode)
	}

	server := metrics.WithServer(ctx)
	RcodeCount.WithLabelValues(server, rc).Add(1)
	RequestDuration.WithLabelValues(server).Observe(res.Rtt.Seconds())

	if err != nil {
		return dns.RcodeServerFailure, err
	}

	// If the client *didn't* set the opt record, and specifically not the DO bit,
	// strip this from the reply (unbound default to setting DO).
	if !state.Do() {
		// technically we can still set bufsize and fluff, for now remove the entire OPT record.
		for i := 0; i < len(res.AnswerPacket.Extra); i++ {
			rr := res.AnswerPacket.Extra[i]
			if _, ok := rr.(*dns.OPT); ok {
				res.AnswerPacket.Extra = append(res.AnswerPacket.Extra[:i], res.AnswerPacket.Extra[i+1:]...)
				break // TODO(miek): more than one? Think TSIG?
			}
		}
		filter(res.AnswerPacket, dnssec)
	}

	res.AnswerPacket.Id = r.Id

	// If the advertised size of the client is smaller than we got, unbound either retried with TCP or something else happened.
	if state.Size() < res.AnswerPacket.Len() {
		res.AnswerPacket, _ = state.Scrub(res.AnswerPacket)
		res.AnswerPacket.Truncated = true
		w.WriteMsg(res.AnswerPacket)

		return 0, nil
	}

	state.SizeAndDo(res.AnswerPacket)
	w.WriteMsg(res.AnswerPacket)

	return 0, nil
}

// Name implements the Handler interface.
func (u *Unbound) Name() string { return "unbound" }
