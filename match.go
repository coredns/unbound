package unbound

import (
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
)

func (u *Unbound) match(state request.Request) bool {
	for _, f := range u.from {
		if plugin.Name(f).Matches(state.Name()) {
			return true
		}
	}

	if u.isAllowedDomain(state.Name()) {
		return true

	}

	return false
}

func (u *Unbound) isAllowedDomain(name string) bool {
	for _, except := range u.except {
		if plugin.Name(except).Matches(name) {
			return false
		}
	}
	return true
}
