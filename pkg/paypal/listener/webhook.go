package listener

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Listener struct {
	debug bool
}

func New(debug bool) *Listener {
	return &Listener{
		debug: debug,
	}
}

// Listen for webhooks
func (l *Listener) WebhooksHandler(cb func(err error, n *Subscription)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			cb(fmt.Errorf("failed to read body: %s", err), nil)
			return
		}

		var subscription Subscription
		err = json.Unmarshal(body, &subscription)
		if err != nil {
			cb(fmt.Errorf("failed to decode request body: %s", err), nil)
			return
		}

		if l.debug {
			fmt.Printf("paypal: body: %s, parsed: %+v\n", body, subscription)
		}

		w.WriteHeader(http.StatusOK)
		cb(nil, &subscription)
	}
}

type Subscription struct {
	ID           string `json:"id"`
	ResourceType string `json:"resource_type"`
	EventType    string `json:"event_type"`
	Summary      string `json:"summary"`
	Resource     struct {
		ID         string `json:"id"`
		Subscriber struct {
			ID        string `json:"id"`
			Summary   string `json:"summary"`
			EventType string `json:"event_type"`
			Name      struct {
				GivenName string `json:"given_name"`
				SurName   string `json:"surname"`
			} `json:"name"`
			Email string `json:"email_address"`
		} `json:"subscriber"`
		ParentPayment string `json:"parent_payment"`
		Amount        struct {
			Total    string `json:"total"`
			Currency string `json:"currency"`
		} `json:"amount"`
	} `json:"resource"`
}
