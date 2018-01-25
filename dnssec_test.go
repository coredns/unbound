package unbound

import (
	"testing"

	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

func TestFilter(t *testing.T) {
	m := new(dns.Msg)
	m.Answer = append(m.Answer, test.MX("miek.nl.		30	IN	MX	10 aspmx3.googlemail.com."))
	m.Answer = append(m.Answer, test.RRSIG("miek.nl.		30	IN	RRSIG	MX 8 2 1800 20180224031301 20180125031301 12051 miek.nl. cOuEEHN14S81aqkAdoZEUqJmpp3jX1X7zjtRDN"))

	filter(m, dnssec)

	if len(m.Answer) != 1 {
		t.Errorf("Expecting 1 RR in answer section, got %d", len(m.Answer))
	}
	if x := m.Answer[0].Header().Rrtype; x != dns.TypeMX {
		t.Errorf("Expecting MX in answer section, got %s", dns.TypeToString[x])
	}
}
