package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dfirebaugh/memberserver/database"
	jwt "github.com/dgrijalva/jwt-go"
)

const JWTExpireMinutes = 1
const CookieName = "memberserver-token"

type accessDetails struct {
	Authorized string `json:"authorized"`
	Expires    string `json:"exp"`
	UserID     string `json:"userID"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

// authJWT middleware that checks JWT token
func (a *API) authJWT(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string
		for _, cookie := range r.Cookies() {
			if cookie.Name == CookieName {
				token = cookie.Value
				break
			}
		}
		r.Header.Set("Authorization", "bearer "+token)
		err := a.TokenValid(r)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		f(w, r) // original function call
	}
}

// getUser responds with the current logged in user
func (a *API) getUser(w http.ResponseWriter, r *http.Request) {
	user, err := a.ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userProfile, err := a.db.GetUser(user.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(userProfile)
	w.Write(j)
}

// Signup register a user to the db
func (a *API) Signup(w http.ResponseWriter, r *http.Request) {
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

// CreateToken create a token for the user to use
func (a *API) CreateToken(userid string) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * JWTExpireMinutes).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(a.config.AccessSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ExtractToken extracts bearer token from httpRequest
func (a *API) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// VerifyToken call ExtractToken
func (a *API) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := a.ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.config.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// TokenValid check to see if the token is valid, still useful, or if it has expired
func (a *API) TokenValid(r *http.Request) error {
	token, err := a.VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

// ExtractTokenMetadata grabs metadata from token
func (a *API) ExtractTokenMetadata(r *http.Request) (*accessDetails, error) {
	token, err := a.VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, err
		}
		return &accessDetails{
			UserID: userID,
		}, nil
	}
	return nil, err
}

// Signin endpoint for user signin
func (a *API) Signin(w http.ResponseWriter, r *http.Request) {
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

	token, err := a.CreateToken(creds.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	tokenJSON := &tokenResponse{}
	tokenJSON.Token = token

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Expires:  time.Now().Add(JWTExpireMinutes * time.Minute),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(tokenJSON)
	w.Write(j)
}

// Logout endpoint for user signin
func (a *API) Logout(w http.ResponseWriter, r *http.Request) {
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
}
