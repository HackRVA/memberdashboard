package integrations

import "time"

type PaymentProvider interface {
	GetSubscription(subscriptionID string) (status string, lastPaymentAmount string, lastPaymentTime time.Time, err error)
	GetSubscriber(subscriptionID string) (name string, email string, err error)
}
