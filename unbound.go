package unbound

import (
	"log"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
	"github.com/miekg/unbound"
	"golang.org/x/net/context"
)

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

	for k, v := range options {
		k += ":" // add :, need for setting options in libunbound
		err := udp.SetOption(k, v)
		if err != nil {
			log.Printf("[WARNING] Could not set option %s with value %s: %s", k, v, err)
		}
		// same failure here, don't repeat log
		tcp.SetOption(k, v)
	}

	return &Unbound{u: udp, t: tcp}
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

	if err != nil {
		return dns.RcodeServerFailure, err
	}

	// If the client *didn't* set the opt record, and specifically not the DO bit,
	// strip this from the reply (unbound default to setting DO).
	// TODO(miek): strip RRSIG/DNSSEC if not directly requested too?
	if !state.Do() {
		// technically we can still set bufsize and fluff, for now remove the entire OPT record.
		for i := 0; i < len(res.AnswerPacket.Extra); i++ {
			rr := res.AnswerPacket.Extra[i]
			if _, ok := rr.(*dns.OPT); ok {
				res.AnswerPacket.Extra = append(res.AnswerPacket.Extra[:i], res.AnswerPacket.Extra[i+1:]...)
				break
			}
		}
	}

	res.AnswerPacket.Id = r.Id
	state.SizeAndDo(res.AnswerPacket)
	w.WriteMsg(res.AnswerPacket)

	return 0, nil
}

// Name implements the Handler interface.
func (u *Unbound) Name() string { return "unbound" }
