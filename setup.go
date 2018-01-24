package unbound

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("unbound", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	u, err := unboundParse(c)
	if err != nil {
		return plugin.Error("unbound", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		u.Next = next
		return u
	})

	return nil
}

func unboundParse(c *caddy.Controller) (*Unbound, error) {
	u := New()

	for c.Next() {

		u.from = c.RemainingArgs()
		if len(u.from) == 0 {
			u.from = make([]string, len(c.ServerBlockKeys))
			copy(u.from, c.ServerBlockKeys)
		}
		for i, str := range u.from {
			u.from[i] = plugin.Host(str).Normalize()
		}

		for c.NextBlock() {
			switch c.Val() {
			case "except":
				except := c.RemainingArgs()
				if len(except) == 0 {
					return nil, c.ArgErr()
				}
				for i := 0; i < len(except); i++ {
					except[i] = plugin.Host(except[i]).Normalize()
				}
				u.except = except
			default:
				return nil, c.ArgErr()
			}
		}
	}
	return u, nil
}
