package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type serviceAuthResponse struct {
	Token string `json:"token"`
}

type servieAuthRequest struct {
	UserName string `json:"email"`
	Password string `json:"password"`
}

type resource struct {
	ID   string `json:"resourceID"`
	Name string `json:"name"`
}

type member struct {
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	RFID        string     `json:"rfid"`
	MemberLevel uint       `json:"memberLevel"`
	Resources   []resource `json:"resources"`
}

func requestToken() (string, error) {
	c := loadConfig()

	authReq := servieAuthRequest{
		UserName: c.ServiceClientID,
		Password: c.ServiceClientSecret,
	}

	reqJSON, _ := json.Marshal(authReq)

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.ServiceURL+"/api/auth/login", bytes.NewBuffer(reqJSON))

	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var newAccessToken serviceAuthResponse

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(res.Body).Decode(&newAccessToken)
	if err != nil {
		return "", err
	}

	return newAccessToken.Token, err
}

func getMembers(token string) []member {
	c := loadConfig()

	req, _ := http.NewRequest("GET", c.ServiceURL+"/api/member", nil)
	client := &http.Client{}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	members := []member{}
	json.Unmarshal(body, &members)

	return members
}
