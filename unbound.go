package unbound

import (
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
	"github.com/miekg/unbound"
	"golang.org/x/net/context"
)

// Unbound is a plugin that resolves requests using libunbound.
type Unbound struct {
	u      *unbound.Unbound
	t      *unbound.Unbound
	Next   plugin.Handler
	from   []string
	except []string
}

// New returns a pointer to an initialzed Unbound.
func New() *Unbound {
	udp := unbound.New()
	tcp := unbound.New()
	err := tcp.SetOption("tcp-upstream:", "yes")
	if err != nil {
		println(err.Error())
	}
	return &Unbound{u: udp, t: tcp}
}

// ServeDNS implements the plugin.Handler interface.
func (u *Unbound) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	var (
		res *unbound.Result
		err error
	)
	switch state.Proto() {
	case "tcp":
		println("wheee")
		res, err = u.t.Resolve(state.QName(), state.QType(), state.QClass())
	case "udp":
		res, err = u.u.Resolve(state.QName(), state.QType(), state.QClass())

	}

	if err != nil {
		return dns.RcodeServerFailure, err
	}

	res.AnswerPacket.Id = r.Id
	state.SizeAndDo(res.AnswerPacket)
	w.WriteMsg(res.AnswerPacket)

	return 0, nil
}

// Name implements the Handler interface.
func (u *Unbound) Name() string { return "unbound" }
