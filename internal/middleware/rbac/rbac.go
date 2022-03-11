package rbac

import (
	"memberserver/internal/services/config"
	"net/http"
	"strings"

	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
)

type UserRole uint

const (
	Admin UserRole = iota + 1
	User
)

type AccessControl interface {
	Restrict(next http.HandlerFunc, allowedRoles []UserRole) http.HandlerFunc
}

type RBAC struct {
	strategy union.Union
}

func New(strategy union.Union) RBAC {
	return RBAC{
		strategy: strategy,
	}
}

func (ur UserRole) ToString() string {
	switch ur {
	case Admin:
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
func (rb RBAC) Restrict(next http.HandlerFunc, allowedRoles []UserRole) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conf, _ := config.Load()
		if strings.Contains(conf.AlwaysAdmin, "true") {
			next.ServeHTTP(w, r)
			return
		}

		_, user, _ := rb.strategy.AuthenticateRequest(r)

		for _, role := range allowedRoles {
			if contains(user.GetGroups(), role.ToString()) {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}
