package unbound

import "testing"

func TestAllowedDomain(t *testing.T) {
	u := New()
	u.except = []string{"download.miek.nl.", "static.miek.nl."}
	u.from = []string{"miek.nl."}

	tests := []struct {
		name     string
		expected bool
	}{
		{"miek.nl.", true},
		{"download.miek.nl.", false},
		{"static.miek.nl.", false},
		{"blaat.miek.nl.", true},
	}

	for i, test := range tests {
		allowed := u.isAllowedDomain(test.name)
		if test.expected != allowed {
			t.Errorf("Test %d: expected %v found %v for %s", i+1, test.expected, allowed, test.name)
		}
	}
}
