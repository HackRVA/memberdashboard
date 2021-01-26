package api

import (
	"encoding/json"
	"fmt"
	"log"
	"memberserver/config"
	"memberserver/database"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTExpireInterval - how long the JWT will last
const JWTExpireInterval = 8

// CookieName - name of the cookie :3
const CookieName = "memberserver-token"

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// tokenReponse for json response of signin
type tokenResponse struct {
	// login response to send token string
	//
	// Example: "<TOKEN_STRING>"
	Token string `json:"token"`
}

// tokenReponseBody for json response of signin
// swagger:response loginResponse
type tokenReponseBody struct {
	// in: body
	Body tokenResponse
}

// swagger:parameters loginRequest
type loginRequest struct {
	// in: body
	Body database.Credentials
}

// swagger:response getUserResponse
type userResponseBody struct {
	// in: body
	Body database.UserResponse
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := verifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}
		name := claims.(jwt.MapClaims)["username"].(string)

		r.Header.Set("name", name)
		r.Header.Set("Authorization", "bearer "+tokenString)

		next.ServeHTTP(w, r)
	})
}

// getUser responds with the current logged in user
func (a API) getUser(w http.ResponseWriter, r *http.Request) {
	userProfile, err := a.db.GetUser(r.Header.Get("name"))
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

	err = a.db.RegisterUser(creds.Username, creds.Password, creds.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
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

func (a API) getToken(name string) (string, error) {

	//Creating Access Token
	atClaims := Claims{}
	atClaims.Username = name
	atClaims.ExpiresAt = time.Now().Add(time.Hour * JWTExpireInterval).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	tokenString, err := token.SignedString([]byte(a.config.AccessSecret))
	return tokenString, err
}

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

func (a API) authenticate(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `Credentials` instance
	creds := &database.Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.db.UserSignin(creds.Username, creds.Password)
	if err != nil {
		fmt.Printf("error signing in: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// If we reach this point, that means the users password was correct, and that they are authorized
	// The default 200 status is sent

	token, err := a.getToken(creds.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	tokenJSON := &tokenResponse{}
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
