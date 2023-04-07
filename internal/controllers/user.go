package controllers

import (
	"net/http"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/internal/datastore"
	"github.com/HackRVA/memberserver/internal/models"

	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	_ "github.com/shaj13/libcache/fifo"
)

type UserServer struct {
	store        datastore.DataStore
	AuthStrategy union.Union
	JWTKeeper    jwt.SecretsKeeper
}

func NewUserServer(store datastore.DataStore, config config.Config) UserServer {
	return UserServer{
		store: store,
	}
}

// getUser responds with the current logged in user
func (us *UserServer) GetUser(w http.ResponseWriter, r *http.Request) {
	u := auth.User(r)
	userProfile, err := us.store.GetMemberByEmail(u.GetUserName())
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	if userProfile.Email == (models.UserResponse{}).Email {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	ok(w, userProfile)
}
