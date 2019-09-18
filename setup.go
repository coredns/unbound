package unbound

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"

	"github.com/caddyserver/caddy"
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

	c.OnStartup(func() error {
		once.Do(func() {
			m := dnsserver.GetConfig(c).Handler("prometheus")
			if m == nil {
				return
			}
			if x, ok := m.(*metrics.Metrics); ok {
				x.MustRegister(RequestDuration)
				x.MustRegister(RcodeCount)
			}
		})
		return nil
	})
	c.OnShutdown(u.Stop)

	return nil
}

func unboundParse(c *caddy.Controller) (*Unbound, error) {
	u := New()

	i := 0
	for c.Next() {
		// Return an error if unbound block specified more than once
		if i > 0 {
			return nil, plugin.ErrOnce
		}
		i++

		u.from = c.RemainingArgs()
		if len(u.from) == 0 {
			u.from = make([]string, len(c.ServerBlockKeys))
			copy(u.from, c.ServerBlockKeys)
		}
		for i, str := range u.from {
			u.from[i] = plugin.Host(str).Normalize()
		}

		for c.NextBlock() {
			var args []string
			var err error

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
			case "option":
				args = c.RemainingArgs()
				if len(args) != 2 {
					return nil, c.ArgErr()
				}
				if err = u.setOption(args[0], args[1]); err != nil {
					return nil, err
				}
			case "config":
				args = c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				if err = u.config(args[0]); err != nil {
					return nil, err
				}
			default:
				return nil, c.ArgErr()
			}
		}
	}
	return u, nil
}
