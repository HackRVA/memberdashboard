package routes

import (
	"github.com/HackRVA/memberserver/internal/middleware/rbac"
	"github.com/HackRVA/memberserver/pkg/paypal/listener"
)

type PaymentsHTTPHandler interface {
	PaypalSubscriptionWebHookHandler(err error, n *listener.Subscription)
}

func (r Router) setupPaymentRoutes(paymentsServer PaymentsHTTPHandler, accessControl rbac.RBAC) {
	webhook := listener.New(true)
	r.UnAuthedRouter.HandleFunc("/api/paypal/subscription/new", webhook.WebhooksHandler(paymentsServer.PaypalSubscriptionWebHookHandler))
}
