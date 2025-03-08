package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/datastore"
	"github.com/HackRVA/memberserver/models"

	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/basic"
	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"github.com/shaj13/libcache"
	log "github.com/sirupsen/logrus"
)

type AuthController struct {
	store            datastore.DataStore
	jwtStrategy      auth.Strategy
	AuthStrategy     union.Union
	JWTSecretsKeeper jwt.SecretsKeeper
}

// JWTExpireInterval - how long the JWT will last in hours
const JWTExpireInterval = 8

func New(dataStore datastore.DataStore) *AuthController {
	c, _ := config.Load()
	auth := AuthController{
		store: dataStore,
	}
	keeper := jwt.StaticSecret{
		ID:        "secret-id",
		Secret:    []byte(c.AccessSecret),
		Algorithm: jwt.HS256,
	}
	cache := libcache.FIFO.New(0)
	cache.SetTTL(time.Minute * 5)

	basicStrategy := basic.NewCached(auth.buildValidator(), cache)
	jwtStrategy := jwt.New(cache, keeper)
	auth.jwtStrategy = jwtStrategy
	auth.AuthStrategy = union.New(jwtStrategy, basicStrategy)
	auth.JWTSecretsKeeper = keeper

	return &auth
}

func (a *AuthController) buildValidator() func(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	return func(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
		log.Errorf("signing in: %s", userName)
		err := a.store.UserSignin(userName, password)
		if err != nil {
			log.Errorf("error signing in: %s", err)
			return nil, fmt.Errorf("invalid credentials")
		}
		// If we reach this point, that means the users password was correct, and that they are authorized
		// we could attach some of their privledges to this return val I think

		// get the user's resources/roles from the db
		user, _ := a.store.GetMemberByEmail(userName)
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
}

func (AuthController) getAuthCookie(request http.Request) string {
	cookie, err := request.Cookie("memberserver")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func (AuthController) setAuthCookie(writer http.ResponseWriter, jwt string) {
	http.SetCookie(writer, &http.Cookie{
		Name:     "memberserver",
		Value:    jwt,
		Expires:  time.Now().Add(JWTExpireInterval * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func (AuthController) removeAuthCookie(writer http.ResponseWriter) {
	http.SetCookie(writer, &http.Cookie{
		Name:     "memberserver",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func (a *AuthController) isValidBearer(request *http.Request) bool {
	requestToken := request.Header.Get("Authorization")

	// if there is no authorization then the user hasn't signed in yet
	if len(requestToken) == 0 {
		return false
	}

	bearerAuth := strings.Split(requestToken, "Bearer ")

	// this is probably a basic auth if it didn't split
	if len(bearerAuth) != 2 {
		return true
	}

	bearerToken := bearerAuth[1]

	return a.getAuthCookie(*request) == bearerToken
}

func (a *AuthController) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, user, err := a.AuthStrategy.AuthenticateRequest(r)
		if err != nil || !a.isValidBearer(r) {
			log.Println(err)
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		r = auth.RequestWithUser(user, r)
		next.ServeHTTP(w, r)
	})
}

// Logout endpoint for user signin
func (a *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	j, _ := json.Marshal(struct{ Message string }{
		Message: "user logged out!",
	})

	a.removeAuthCookie(w)

	ok(w, j)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	println("attempting to login")

	var body []byte
	_, err := r.Body.Read(body)
	if err != nil {
		log.Errorf("error reading body %s", err)
	}
	log.Debug(string(body))
	exp := jwt.SetExpDuration(time.Hour * JWTExpireInterval)
	u := auth.User(r)
	u.SetUserName(strings.ToLower(u.GetUserName()))
	token, err := jwt.IssueAccessToken(u, a.JWTSecretsKeeper, exp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	tokenJSON := &models.TokenResponse{
		Token: token,
	}

	a.setAuthCookie(w, tokenJSON.Token)

	ok(w, tokenJSON)
}

// Signup register a user to the db
func (a *AuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {
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

	err = a.store.RegisterUser(models.Credentials{
		Email:    strings.ToLower(creds.Email),
		Password: creds.Password,
	})
	if err != nil {
		log.Error(err)
		http.Error(w, "error registering user", http.StatusBadRequest)
		return
	}

	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
	ok(w, models.EndpointSuccess{Ack: true})
}

func ok(writer http.ResponseWriter, result interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(result)
	if _, err := writer.Write(response); err != nil {
		log.Errorf("error reading rsponse %s", err)
	}
}
