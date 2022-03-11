package integrations

import "memberserver/internal/models"

type PaymentProvider interface {
	GetSubscription(subscriptionID string) (string, models.Payment, error)
	GetMemberFromSubscription(subscriptionID string) (models.Member, error)
}
