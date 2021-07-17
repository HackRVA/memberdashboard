package api

import (
	"memberserver/config"
	"net/http"
	"strings"
)

type UserRole uint

const (
	admin UserRole = iota + 1
	user
)

func (ur UserRole) ToString() string {
	switch ur {
	case admin:
		return "admin"
	default:
		return "user"
	}
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// rbac is middleware that will restrict access based on the roles you pass in
func (api API) rbac(next http.HandlerFunc, allowedRoles []UserRole) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conf, _ := config.Load()
		if strings.Contains(conf.AlwaysAdmin, "true") {
			next.ServeHTTP(w, r)
			return
		}

		_, user, _ := strategy.AuthenticateRequest(r)

		for _, role := range allowedRoles {
			if contains(user.GetGroups(), role.ToString()) {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}
