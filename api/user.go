package api

import (
	"context"
	"encoding/json"
	"fmt"
	"memberserver/api/models"
	"memberserver/config"
	"memberserver/datastore"
	"net/http"
	"strings"
	"time"

	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/basic"
	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"github.com/shaj13/libcache"
	_ "github.com/shaj13/libcache/fifo"
	log "github.com/sirupsen/logrus"
)

// get
// register
// login
// logout

type UserStore interface {
	GetMemberByEmail(email string) (models.Member, error)
	RegisterUser(models.Credentials) error
	UserSignin(username, password string) error
}

type UserServer struct {
	store datastore.DataStore
}

// JWTExpireInterval - how long the JWT will last in hours
const JWTExpireInterval = 8

var strategy union.Union
var keeper jwt.SecretsKeeper

func setupAuth(config config.Config, userServer UserServer) {

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
	basicStrategy := basic.NewCached(userServer.getValidator(), cache)
	jwtStrategy := jwt.New(cache, keeper)
	strategy = union.New(jwtStrategy, basicStrategy)
}

func (us *UserServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, user, err := strategy.AuthenticateRequest(r)
		if err != nil && !us.isValidBearer(r) {
			log.Println(err)
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		r = auth.RequestWithUser(user, r)
		next.ServeHTTP(w, r)
	})
}

func NewUserServer(store datastore.DataStore, config config.Config) UserServer {
	userServer := UserServer{
		store: store,
	}

	return userServer
}

// getUser responds with the current logged in user
func (us *UserServer) getUser(w http.ResponseWriter, r *http.Request) {
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

func (us *UserServer) getValidator() basic.AuthenticateFunc {
	validator := func(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
		log.Errorf("signing in: %s", userName)
		err := us.store.UserSignin(userName, password)
		if err != nil {
			log.Errorf("error signing in: %s", err)
			return nil, fmt.Errorf("invalid credentials")
		}
		// If we reach this point, that means the users password was correct, and that they are authorized
		// we could attach some of their privledges to this return val I think

		// get the user's resources/roles from the db
		user, _ := us.store.GetMemberByEmail(userName)
		var resources []string
		for _, resource := range user.Resources {
			resources = append(resources, resource.Name)
		}

		conf, _ := config.Load()
		if strings.Contains(conf.AlwaysAdmin, "true") {
			resources = append(resources, "admin")
		}

		return auth.NewDefaultUser(userName, userName, resources, nil), nil
	}
	return validator
}

// Signup register a user to the db
func (us *UserServer) registerUser(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `Credentials` instance
	creds := &models.Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		http.Error(w, "error registering user", http.StatusBadRequest)
		return
	}

	if len(creds.Password) < 3 {
		http.Error(w, "password must be longer", http.StatusBadRequest)
		return
	}

	err = us.store.RegisterUser(models.Credentials{
		Email:    strings.ToLower(creds.Email),
		Password: creds.Password,
	})

	if err != nil {
		log.Error(err)
		http.Error(w, "error registering user", http.StatusBadRequest)
		return
	}

	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
	j, _ := json.Marshal(models.EndpointSuccess{
		Ack: true,
	})

	ok(w, j)
}

func (us *UserServer) login(w http.ResponseWriter, r *http.Request) {
	exp := jwt.SetExpDuration(time.Hour * JWTExpireInterval)
	u := auth.User(r)
	u.SetUserName(strings.ToLower(u.GetUserName()))
	token, err := jwt.IssueAccessToken(u, keeper, exp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	tokenJSON := &models.TokenResponse{
		Token: token,
	}

	us.setAuthCookie(w, tokenJSON.Token)

	ok(w, tokenJSON)
}

// Logout endpoint for user signin
func (us *UserServer) logout(w http.ResponseWriter, r *http.Request) {
	j, _ := json.Marshal(struct{ Message string }{
		Message: "user logged out!",
	})

	us.removeAuthCookie(w)

	ok(w, j)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (us *UserServer) getAuthCookie(request http.Request) string {
	cookie, err := request.Cookie("memberserver")

	if err != nil {
		return ""
	}

	return cookie.Value
}

func (us *UserServer) setAuthCookie(writer http.ResponseWriter, jwt string) {
	http.SetCookie(writer, &http.Cookie{
		Name:     "memberserver",
		Value:    jwt,
		Expires:  time.Now().Add(JWTExpireInterval * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func (us *UserServer) removeAuthCookie(writer http.ResponseWriter) {
	http.SetCookie(writer, &http.Cookie{
		Name:     "memberserver",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func (us *UserServer) isValidBearer(request *http.Request) bool {
	requestToken := request.Header.Get("Authorization")

	// if there is no authorization then the user hasn't signed in yet
	if len(requestToken) == 0 {
		return false
	}

	bearerToken := strings.Split(requestToken, "Bearer ")[1]

	// this is probably a basic auth if it's empty
	if len(bearerToken) == 0 {
		return true
	}

	return us.getAuthCookie(*request) == bearerToken
}
