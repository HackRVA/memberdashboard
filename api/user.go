package api

import (
	"encoding/json"
	"memberserver/api/models"
	"memberserver/config"
	"memberserver/database"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
    "github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/basic"
	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	log "github.com/sirupsen/logrus"
)

// JWTExpireInterval - how long the JWT will last in hours
const JWTExpireInterval = 8

// CookieName - name of the cookie :3
const CookieName = "memberserver-token"

var strategy union.Union
var keeper jwt.SecretsKeeper

func setupGoGuardian(config) {

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
    basicStrategy := basic.NewCached(validateUser, cache)
    jwtStrategy := jwt.New(cache, keeper)
    strategy = union.New(jwtStrategy, basicStrategy)
}


func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("Executing Auth Middleware")
        _, user, err := strategy.AuthenticateRequest(r)
        if err != nil {
            fmt.Println(err)
            code := http.StatusUnauthorized
            http.Error(w, http.StatusText(code), code)
            return
        }
        log.Printf("User %s Authenticated\n", user.GetUserName())
        r = auth.RequestWithUser(user, r)
        next.ServeHTTP(w, r)
    })
}

// getUser responds with the current logged in user
func (a API) getUser(w http.ResponseWriter, r *http.Request) {
	userProfile, err := a.db.GetUser(r.Header.Get("email"))
	if err != nil {
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
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(struct{ Message string }{
		Message: "user logged out!",
	})
	w.Write(j)
	http.Redirect(w, r, "/", http.StatusFound)
}

//No longer needed
func (a API) getToken(email string) (string, error) {
	//Creating Access Token
	atClaims := models.Claims{}
	atClaims.Email = email
	atClaims.ExpiresAt = time.Now().Add(time.Hour * JWTExpireInterval).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	tokenString, err := token.SignedString([]byte(a.config.AccessSecret))
	return tokenString, err
}


//no longer needed
func verifyToken(tokenString string) (jwt.Claims, error) {
	c, _ := config.Load()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}

func (a API) validateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	err = a.db.UserSignin(userName, password)
	if err != nil {
		log.Errorf("error signing in: %s", err)
        return nil, fmt.Errorf("Invalid credentials")
	}
	// If we reach this point, that means the users password was correct, and that they are authorized
    // we could attach some of their privledges to this return val I think
    return auth.NewDefaultUser(userName, userName, nil, nil), nil
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

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Expires:  time.Now().Add(JWTExpireInterval * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(tokenJSON)
	w.Write(j)
}
