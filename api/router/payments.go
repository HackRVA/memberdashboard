package router

import (
	"memberserver/api/rbac"
	"memberserver/payments/listener"
	"net/http"
)

type PaymentsHTTPHandler interface {
	GetPaymentChart(http.ResponseWriter, *http.Request)
	PaypalSubscriptionWebHookHandler(err error, n *listener.PaypalNotification)
}

func (r Router) setupPaymentRoutes(paymentsServer PaymentsHTTPHandler, accessControl rbac.RBAC) {
	webhook := listener.New(true)
	r.UnAuthedRouter.HandleFunc("/api/paypal/subscription/new", webhook.WebhooksHandler(paymentsServer.PaypalSubscriptionWebHookHandler))
	// swagger:route GET /api/payments/charts payments searchPaymentChartRequest
	//
	// Get Chart information of payments
	//
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: getPaymentChartResponse
	r.authedRouter.HandleFunc("/payments/charts", accessControl.Restrict(paymentsServer.GetPaymentChart, []rbac.UserRole{rbac.Admin}))
}
