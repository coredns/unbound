package unbound

import (
	"testing"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

func TestUnbound(t *testing.T) {
	u := New()

	tests := []struct {
		qname         string
		qtype         uint16
		expectedRcode int
		expectedReply []string // ownernames for the records in the answer section.
		expectedErr   error
	}{
		{
			qname:         "example.org",
			qtype:         dns.TypeA,
			expectedRcode: dns.RcodeSuccess,
			expectedReply: []string{"example.org."},
			expectedErr:   nil,
		},
		{
			qname:         "no-such-record.example.org",
			qtype:         dns.TypeA,
			expectedRcode: dns.RcodeNameError,
			expectedReply: nil,
			expectedErr:   nil,
		},
	}

	ctx := context.TODO()

	for i, tc := range tests {
		req := new(dns.Msg)
		req.SetQuestion(dns.Fqdn(tc.qname), tc.qtype)

		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		_, err := u.ServeDNS(ctx, rec, req)

		if err != tc.expectedErr {
			t.Errorf("Test %d: Expected error %v, but got %v", i, tc.expectedErr, err)
		}
		if rec.Msg.Rcode != int(tc.expectedRcode) {
			t.Errorf("Test %d: Expected status code %d, but got %d", i, tc.expectedRcode, rec.Msg.Rcode)
		}
		if len(tc.expectedReply) != 0 {
			for i, expected := range tc.expectedReply {
				actual := rec.Msg.Answer[i].Header().Name
				if actual != expected {
					t.Errorf("Test %d: Expected answer %s, but got %s", i, expected, actual)
				}
			}
		}
	}
}
