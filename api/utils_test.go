package api

import (
	"context"
	"errors"
	"testing"
)

// mockResolver is a mock DNS resolver for testing
type mockResolver struct {
	names []string
	err   error
}

func (m *mockResolver) LookupAddr(ctx context.Context, addr string) ([]string, error) {
	return m.names, m.err
}

func TestResolveHostname_Success(t *testing.T) {
	tests := []struct {
		name      string
		ipAddress string
		dnsNames  []string
		expected  string
	}{
		{
			name:      "single hostname returned",
			ipAddress: "192.168.1.1",
			dnsNames:  []string{"server.example.com."},
			expected:  "server.example.com",
		},
		{
			name:      "multiple hostnames returned",
			ipAddress: "10.0.0.1",
			dnsNames:  []string{"primary.example.com.", "secondary.example.com."},
			expected:  "primary.example.com",
		},
		{
			name:      "hostname without trailing dot",
			ipAddress: "172.16.0.1",
			dnsNames:  []string{"server.local"},
			expected:  "server.local",
		},
		{
			name:      "empty dns names list",
			ipAddress: "8.8.8.8",
			dnsNames:  []string{},
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &mockResolver{names: tt.dnsNames, err: nil}
			result := ResolveHostname(tt.ipAddress, resolver)
			if result != tt.expected {
				t.Errorf("ResolveHostname(%s) = %q, want %q", tt.ipAddress, result, tt.expected)
			}
		})
	}
}

func TestResolveHostname_Error(t *testing.T) {
	tests := []struct {
		name      string
		ipAddress string
		err       error
		expected  string
	}{
		{
			name:      "dns lookup error",
			ipAddress: "192.168.1.1",
			err:       errors.New("no such host"),
			expected:  "",
		},
		{
			name:      "timeout error",
			ipAddress: "10.0.0.1",
			err:       errors.New("i/o timeout"),
			expected:  "",
		},
		{
			name:      "network unreachable",
			ipAddress: "172.16.0.1",
			err:       errors.New("network is unreachable"),
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &mockResolver{names: nil, err: tt.err}
			result := ResolveHostname(tt.ipAddress, resolver)
			if result != tt.expected {
				t.Errorf("ResolveHostname(%s) with error = %q, want %q", tt.ipAddress, result, tt.expected)
			}
		})
	}
}

func TestResolveHostname_NilResolver(t *testing.T) {
	// when resolver is nil, it should use DefaultDNSResolver
	// this test mainly verifies no panic occurs
	// the actual result depends on the test environment
	result := ResolveHostname("127.0.0.1", nil)
	// we don't assert the result since it depends on the environment
	// but it should not panic and should return a string (possibly empty)
	_ = result
}

func TestResolveHostname_TrimsTrailingDot(t *testing.T) {
	resolver := &mockResolver{
		names: []string{"hostname.with.trailing.dot."},
		err:   nil,
	}

	result := ResolveHostname("1.2.3.4", resolver)
	expected := "hostname.with.trailing.dot"

	if result != expected {
		t.Errorf("ResolveHostname() = %q, want %q (trailing dot should be removed)", result, expected)
	}
}

func TestResolveHostname_IPv6Address(t *testing.T) {
	resolver := &mockResolver{
		names: []string{"ipv6host.example.com."},
		err:   nil,
	}

	result := ResolveHostname("2001:db8::1", resolver)
	expected := "ipv6host.example.com"

	if result != expected {
		t.Errorf("ResolveHostname(IPv6) = %q, want %q", result, expected)
	}
}
