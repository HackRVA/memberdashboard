package router

import (
	"memberserver/api/rbac"
	"memberserver/payments/listener"
	"net/http"

	"github.com/gorilla/mux"
)

type PaymentsHTTPHandler interface {
	GetPaymentChart(http.ResponseWriter, *http.Request)
	PaypalSubscriptionWebHookHandler(err error, n *listener.PaypalNotification)
}

func setupPaymentRoutes(unauthedRouter *mux.Router, authedRouter *mux.Router, paymentsServer PaymentsHTTPHandler, accessControl rbac.RBAC) (*mux.Router, *mux.Router) {
	webhook := listener.New(true)
	unauthedRouter.HandleFunc("/api/paypal/subscription/new", webhook.WebhooksHandler(paymentsServer.PaypalSubscriptionWebHookHandler))
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
	authedRouter.HandleFunc("/payments/charts", accessControl.Restrict(paymentsServer.GetPaymentChart, []rbac.UserRole{rbac.Admin}))
	return unauthedRouter, authedRouter
}
