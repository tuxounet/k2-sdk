package routes

import (
	"testing"

	runtimeTypes "github.com/tuxounet/k2-sdk/types"
)

func TestMatchLongestPrefix(t *testing.T) {
	tests := []struct {
		name        string
		requestPath string
		authMap     map[string]runtimeTypes.IAccessPolicy
		wantPolicy  runtimeTypes.IAccessPolicy
	}{
		{
			name:        "empty authMap returns empty",
			requestPath: "/api/hello/sayHello",
			authMap:     map[string]runtimeTypes.IAccessPolicy{},
			wantPolicy:  "",
		},
		{
			name:        "exact match",
			requestPath: "/api/hello/",
			authMap: map[string]runtimeTypes.IAccessPolicy{
				"/api/hello/": runtimeTypes.AccessPolicyPublic,
			},
			wantPolicy: runtimeTypes.AccessPolicyPublic,
		},
		{
			name:        "no matching prefix returns empty",
			requestPath: "/other/path",
			authMap: map[string]runtimeTypes.IAccessPolicy{
				"/api/": runtimeTypes.AccessPolicyPublic,
			},
			wantPolicy: "",
		},
		{
			name:        "single prefix shorter than request",
			requestPath: "/api/hello/sayHello",
			authMap: map[string]runtimeTypes.IAccessPolicy{
				"/api/": runtimeTypes.AccessPolicyAuthenticated,
			},
			wantPolicy: runtimeTypes.AccessPolicyAuthenticated,
		},
		{
			name:        "longest prefix wins — controller public over component authenticated",
			requestPath: "/f-app/v1/download",
			authMap: map[string]runtimeTypes.IAccessPolicy{
				"/":                   runtimeTypes.AccessPolicyPublic,
				"/f-app/v1/":          runtimeTypes.AccessPolicyAuthenticated,
				"/f-app/v1/download/": runtimeTypes.AccessPolicyPublic,
				"/f-app/v1/publish/":  runtimeTypes.AccessPolicyAuthenticated,
			},
			wantPolicy: runtimeTypes.AccessPolicyPublic,
		},
		{
			name:        "longest prefix wins — authenticated controller over public component",
			requestPath: "/f-app/v1/publish/artifact",
			authMap: map[string]runtimeTypes.IAccessPolicy{
				"/":                  runtimeTypes.AccessPolicyPublic,
				"/f-app/v1/":         runtimeTypes.AccessPolicyPublic,
				"/f-app/v1/publish/": runtimeTypes.AccessPolicyAuthenticated,
			},
			wantPolicy: runtimeTypes.AccessPolicyAuthenticated,
		},
		{
			name:        "longest prefix wins — health endpoint (issue #3 regression)",
			requestPath: "/f-app/v1/health",
			authMap: map[string]runtimeTypes.IAccessPolicy{
				"/":                   runtimeTypes.AccessPolicyPublic,
				"/f-app/v1/":          runtimeTypes.AccessPolicyAuthenticated,
				"/f-app/v1/health/":   runtimeTypes.AccessPolicyPublic,
				"/f-app/v1/download/": runtimeTypes.AccessPolicyPublic,
			},
			wantPolicy: runtimeTypes.AccessPolicyPublic,
		},
		{
			name:        "deeply nested path with multiple levels",
			requestPath: "/api/v2/users/admin/profile",
			authMap: map[string]runtimeTypes.IAccessPolicy{
				"/":                    runtimeTypes.AccessPolicyPublic,
				"/api/":                runtimeTypes.AccessPolicyAuthenticated,
				"/api/v2/":             runtimeTypes.AccessPolicyAuthenticated,
				"/api/v2/users/":       runtimeTypes.AccessPolicyAuthenticated,
				"/api/v2/users/admin/": runtimeTypes.AccessPolicyPublic,
			},
			wantPolicy: runtimeTypes.AccessPolicyPublic,
		},
		{
			name:        "prefixes with trailing slash handled correctly",
			requestPath: "/admin/dashboard",
			authMap: map[string]runtimeTypes.IAccessPolicy{
				"/admin":  runtimeTypes.AccessPolicyAuthenticated,
				"/admin/": runtimeTypes.AccessPolicyPublic,
			},
			wantPolicy: runtimeTypes.AccessPolicyPublic,
		},
		{
			name:        "request exactly matches shorter prefix rather than longer",
			requestPath: "/api",
			authMap: map[string]runtimeTypes.IAccessPolicy{
				"/":        runtimeTypes.AccessPolicyPublic,
				"/api":     runtimeTypes.AccessPolicyAuthenticated,
				"/api/v1/": runtimeTypes.AccessPolicyPublic,
			},
			wantPolicy: runtimeTypes.AccessPolicyAuthenticated,
		},
		{
			name:        "deterministic regardless of map iteration order (same length)",
			requestPath: "/path/to/resource",
			authMap: map[string]runtimeTypes.IAccessPolicy{
				"/path/":       runtimeTypes.AccessPolicyPublic,
				"/other/path/": runtimeTypes.AccessPolicyAuthenticated,
			},
			wantPolicy: runtimeTypes.AccessPolicyPublic,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run multiple times to ensure deterministic behavior
			// (map iteration order is randomized in Go)
			for i := 0; i < 50; i++ {
				got := matchLongestPrefix(tt.requestPath, tt.authMap)
				if got != tt.wantPolicy {
					t.Errorf("matchLongestPrefix(%q) = %q, want %q (iteration %d)\n  authMap keys: %v",
						tt.requestPath, got, tt.wantPolicy, i, mapKeys(tt.authMap))
					break
				}
			}
		})
	}
}

// mapKeys returns the keys of a map for diagnostic output.
func mapKeys(m map[string]runtimeTypes.IAccessPolicy) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestMatchLongestPrefix_Determinism(t *testing.T) {
	// This test specifically verifies that the longest-prefix behavior is
	// deterministic across many iterations, solving the issue #3 scenario
	// where map iteration randomness caused intermittent 401s.

	authMap := map[string]runtimeTypes.IAccessPolicy{
		"/":                   runtimeTypes.AccessPolicyPublic,
		"/f-app/v1/":          runtimeTypes.AccessPolicyAuthenticated,
		"/f-app/v1/download/": runtimeTypes.AccessPolicyPublic,
		"/f-app/v1/health/":   runtimeTypes.AccessPolicyPublic,
		"/f-app/v1/publish/":  runtimeTypes.AccessPolicyAuthenticated,
		"/f-app/v1/admin/":    runtimeTypes.AccessPolicyAuthenticated,
	}

	tests := []struct {
		path string
		want runtimeTypes.IAccessPolicy
	}{
		{"/f-app/v1/download", runtimeTypes.AccessPolicyPublic},
		{"/f-app/v1/download/linux/install.sh", runtimeTypes.AccessPolicyPublic},
		{"/f-app/v1/download/", runtimeTypes.AccessPolicyPublic},
		{"/f-app/v1/health", runtimeTypes.AccessPolicyPublic},
		{"/f-app/v1/health/status", runtimeTypes.AccessPolicyPublic},
		{"/f-app/v1/publish", runtimeTypes.AccessPolicyAuthenticated},
		{"/f-app/v1/publish/artifact", runtimeTypes.AccessPolicyAuthenticated},
		{"/f-app/v1/admin", runtimeTypes.AccessPolicyAuthenticated},
		{"/f-app/v1/admin/users", runtimeTypes.AccessPolicyAuthenticated},
		{"/f-app/v1/", runtimeTypes.AccessPolicyAuthenticated},
		{"/f-app/v1", runtimeTypes.AccessPolicyAuthenticated}, // suffix match: "/f-app/v1" + "/" == "/f-app/v1/"
		{"/other", runtimeTypes.AccessPolicyPublic},           // only "/" matches
		{"/unknown/path", runtimeTypes.AccessPolicyPublic},    // only "/" matches
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				got := matchLongestPrefix(tt.path, authMap)
				if got != tt.want {
					t.Errorf("matchLongestPrefix(%q) = %q, want %q (iteration %d)",
						tt.path, got, tt.want, i)
					return
				}
			}
		})
	}
}
