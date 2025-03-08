package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Profile struct {
	DisplayName string `json:"display_name"`
	RealName    string `json:"real_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
}

type SlackUser struct {
	Name     string  `json:"name"`
	RealName string  `json:"real_name"`
	Profile  Profile `json:"profile"`
	IsBot    bool    `json:"is_bot"`
	Deleted  bool    `json:"deleted"`
}

type SlackUserResponse struct {
	Members []SlackUser `json:"members"`
}

func GetUsers(token string) ([]SlackUser, error) {
	var slackUsers []SlackUser
	if len(token) == 0 {
		return slackUsers, errors.New("slack Token not found")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", `https://slack.com/api/users.list`, nil)
	if err != nil {
		return slackUsers, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return slackUsers, err
	}
	defer res.Body.Close()

	var slackUsersResponse SlackUserResponse

	err = json.NewDecoder(res.Body).Decode(&slackUsersResponse)
	if err != nil {
		return slackUsers, err
	}

	slackUsers = slackUsersResponse.Members

	return slackUsers, err
}

func Send(slackWebHook string, msg string) {
	newMsg := fmt.Sprint("{\"text\":'```", msg, "```'}")
	jsonStr := []byte(newMsg)
	if len(slackWebHook) == 0 {
		log.Debugf("slack web hook isn't set")
		return
	}

	c := &http.Client{}
	req, err := http.NewRequest("POST", slackWebHook, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Errorf("some error sending to slack hook: %s", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		log.Errorf("some error sending to slack hook: %s", err)
		return
	}

	defer res.Body.Close()
}
