package api

import (
	"context"
	"encoding/json"
	"fmt"
	"memberserver/api/models"
	"memberserver/config"
	"memberserver/database"
	"net/http"
	"time"

	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/basic"
	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"github.com/shaj13/libcache"
	_ "github.com/shaj13/libcache/fifo"
	log "github.com/sirupsen/logrus"
)

// JWTExpireInterval - how long the JWT will last in hours
const JWTExpireInterval = 8

var strategy union.Union
var keeper jwt.SecretsKeeper

func setupGoGuardian(config config.Config, db *database.Database) {

	keeper = jwt.StaticSecret{
		ID:        "secret-id",
		Secret:    []byte(config.AccessSecret),
		Algorithm: jwt.HS256,
	}
	cache := libcache.FIFO.New(0)
	cache.SetTTL(time.Minute * 5)
	cache.RegisterOnExpired(func(key, _ interface{}) {
		cache.Peek(key)
	})
	basicStrategy := basic.NewCached(getValidator(db), cache)
	jwtStrategy := jwt.New(cache, keeper)
	strategy = union.New(jwtStrategy, basicStrategy)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, user, err := strategy.AuthenticateRequest(r)
		if err != nil {
			log.Println(err)
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		r = auth.RequestWithUser(user, r)
		next.ServeHTTP(w, r)
	})
}

// getUser responds with the current logged in user
func (a API) getUser(w http.ResponseWriter, r *http.Request) {
	u := auth.User(r)
	userProfile, err := a.db.GetUser(u.GetUserName())
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userProfile == (database.UserResponse{}) {
		log.Println("user not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(userProfile)
	w.Write(j)
}

// Signup register a user to the db
func (a API) signup(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `Credentials` instance
	creds := &database.Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.db.RegisterUser(creds.Email, creds.Password)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(models.EndpointSuccess{
		Ack: true,
	})
	w.Write(j)
}

// Logout endpoint for user signin
func (a API) logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(struct{ Message string }{
		Message: "user logged out!",
	})
	w.Write(j)
	http.Redirect(w, r, "/", http.StatusFound)
}

func getValidator(db *database.Database) basic.AuthenticateFunc {
	validator := func(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
		log.Errorf("signing in: %s", userName)
		err := db.UserSignin(userName, password)
		if err != nil {
			log.Errorf("error signing in: %s", err)
			return nil, fmt.Errorf("Invalid credentials")
		}
		// If we reach this point, that means the users password was correct, and that they are authorized
		// we could attach some of their privledges to this return val I think
		return auth.NewDefaultUser(userName, userName, nil, nil), nil
	}
	return validator
}

func (a API) authenticate(w http.ResponseWriter, r *http.Request) {
	exp := jwt.SetExpDuration(time.Hour * JWTExpireInterval)
	u := auth.User(r)
	token, err := jwt.IssueAccessToken(u, keeper, exp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	tokenJSON := &models.TokenResponse{}
	tokenJSON.Token = token

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(tokenJSON)
	w.Write(j)
}
