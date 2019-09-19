package unbound

import (
	"testing"

	"github.com/caddyserver/caddy"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
	}{
		{`unbound`, false},
		{`unbound .`, false},
		{`unbound a b`, false},
	}

	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		_, err := unboundParse(c)

		if test.shouldErr && err == nil {
			t.Errorf("Test %d: Expected error but found none for input %s", i, test.input)
		}

		if err != nil {
			if !test.shouldErr {
				t.Errorf("Test %d: Expected no error but found one for input %s. Error was: %v", i, test.input, err)
			}
		}
	}
}

func TestSetupExtended(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
	}{
		{`unbound {
			option msg-cache-size 0
			option msg-cache-size 0
		}`, false},
		{`unbound {
			option msg-cache-size 0
			except example.org example.net
		}`, false},

		{`unbound {
			option bla yes
		}`, true},
		{`unbound {
			optoin qname-minimisation yes
		}`, true},
		{`unbound {
			config testdata/unbound.conf.bad
		}`, true},
		{`unbound {
			config testdata/unbound.conf.good
		}`, false},
	}

	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		_, err := unboundParse(c)

		if test.shouldErr && err == nil {
			t.Errorf("Test %d: Expected error but found none for input %s", i, test.input)
		}

		if err != nil {
			if !test.shouldErr {
				t.Errorf("Test %d: Expected no error but found one for input %s. Error was: %v", i, test.input, err)
			}
		}
	}
}
