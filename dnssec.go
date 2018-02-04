package unbound

import "github.com/miekg/dns"

// filter removes records from m, according to filter function.
func filter(m *dns.Msg, filter func(dns.RR) bool) {
	rrs := []dns.RR{}
	for _, r := range m.Answer {
		if !filter(r) {
			rrs = append(rrs, r)
		}
	}
	m.Answer = rrs

	rrs = []dns.RR{}
	for _, r := range m.Ns {
		if !filter(r) {
			rrs = append(rrs, r)
		}

	}
	m.Ns = rrs

	rrs = []dns.RR{}
	for _, r := range m.Extra {
		if !filter(r) {
			rrs = append(rrs, r)
		}

	}
	m.Extra = rrs
}

// dnssec returns true if rr is an RRSIG, NSEC or NSEC3 record.
func dnssec(rr dns.RR) bool {
	if _, ok := rr.(*dns.RRSIG); ok {
		return true
	}
	if _, ok := rr.(*dns.NSEC); ok {
		return true
	}
	if _, ok := rr.(*dns.NSEC3); ok {
		return true
	}
	return false
}
