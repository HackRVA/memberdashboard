package router

import (
	"memberserver/api/rbac"
	"memberserver/payments/listener"
)

type PaymentsHTTPHandler interface {
	PaypalSubscriptionWebHookHandler(err error, n *listener.PaypalNotification)
}

func (r Router) setupPaymentRoutes(paymentsServer PaymentsHTTPHandler, accessControl rbac.RBAC) {
	webhook := listener.New(true)
	r.UnAuthedRouter.HandleFunc("/api/paypal/subscription/new", webhook.WebhooksHandler(paymentsServer.PaypalSubscriptionWebHookHandler))
}
