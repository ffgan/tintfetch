package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicIPResponseValidation(t *testing.T) {
	for _, tc := range []struct {
		name, body, want string
		status           int
	}{
		{"IPv4", "203.0.113.8\n", "203.0.113.8", 200},
		{"IPv6", "2001:db8::1", "2001:db8::1", 200},
		{"HTML", "<html>203.0.113.8</html>", "", 200},
		{"failure", "203.0.113.8", "", 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()
			if got := publicIP(server.URL, "1"); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
	if got := publicIP("file:///etc/passwd", "1"); got != "" {
		t.Fatal("invalid scheme accepted")
	}
}
