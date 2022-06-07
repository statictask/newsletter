package server

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/statictask/newsletter/pkg/subscription"
)

var (
	originsOk = handlers.AllowedOrigins([]string{"*"})
	methodsOk = handlers.AllowedMethods(
		[]string{"GET", "HEAD", "POST", "UPDATE", "DELETE", "PUT", "OPTIONS"},
	)
	headersOk = handlers.AllowedHeaders(
		[]string{
			"Accept",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
		},
	)
)

func (s *Server) Listen(bind string) error {
	router := mux.NewRouter()

	// subscriptions routes
	router.HandleFunc("/subscriptions", subscription.GetSubscriptions).Methods("GET")
	router.HandleFunc("/subscriptions", subscription.CreateSubscription).Methods("POST")
	router.HandleFunc("/subscriptions/{subscription_id}", subscription.GetSubscription).Methods("GET")
	router.HandleFunc("/subscriptions/{subscription_id}", subscription.DeleteSubscription).Methods("DELETE")
	router.HandleFunc("/subscriptions/{subscription_id}", subscription.UpdateSubscription).Methods("UPDATE")

	s.L.With(zap.String("bind", bind)).Info("listening")

	r := handlers.CORS(originsOk, headersOk, methodsOk)(router)
	if err := http.ListenAndServe(bind, r); err != nil {
		return err
	}

	return nil
}
