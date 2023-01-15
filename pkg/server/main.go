package server

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/statictask/newsletter/pkg/project"
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
			"Origin",
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

	// projects routes
	router.HandleFunc("/projects", project.CreateProject).Methods("POST")
	router.HandleFunc("/projects/{project_id}", project.GetProject).Methods("GET")
	router.HandleFunc("/projects/{project_id}", project.DeleteProject).Methods("DELETE")
	router.HandleFunc("/projects/{project_id}", project.UpdateProject).Methods("UPDATE")

	// subscription routes
	router.HandleFunc("/projects/{project_id}/subscriptions", subscription.GetProjectSubscriptions).Methods("GET")
	router.HandleFunc("/projects/{project_id}/subscriptions", subscription.CreateProjectSubscription).Methods("POST")
	router.HandleFunc("/projects/{project_id}/subscriptions/{subscription_id}", subscription.GetProjectSubscription).Methods("GET")
	router.HandleFunc("/projects/{project_id}/subscriptions/{subscription_id}", subscription.DeleteProjectSubscription).Methods("DELETE")
	router.HandleFunc("/projects/{project_id}/subscriptions/{subscription_id}", subscription.UpdateProjectSubscription).Methods("UPDATE")
	router.HandleFunc("/projects/{project_id}/subscriptions/{subscription_id}/encrypt", subscription.GetSubscriptionEncryption).Methods("GET")
	router.HandleFunc("/unsubscribe", subscription.GetUnsubscribePage).Queries("token", "{token}").Methods("GET")
	router.HandleFunc("/unsubscribe", subscription.DeleteUnsubscribe).Queries("token", "{token}").Methods("DELETE")
	router.HandleFunc("/goodbye", subscription.GetGoodbyePage).Methods("GET")

	s.L.With(zap.String("bind", bind)).Info("listening")

	r := handlers.CORS(originsOk, headersOk, methodsOk)(router)
	if err := http.ListenAndServe(bind, r); err != nil {
		return err
	}

	return nil
}
