package server

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/statictask/newsletter/pkg/project"
	"github.com/statictask/newsletter/pkg/subscription"
	"github.com/statictask/newsletter/pkg/template"
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

	// email template routes
	router.HandleFunc("/projects/{project_id}/templates", template.GetEmailTemplates).Methods("GET")
	router.HandleFunc("/projects/{project_id}/templates", template.CreateEmailTemplate).Methods("POST")
	router.HandleFunc("/projects/{project_id}/templates/{email_template_id}", template.GetEmailTemplate).Methods("GET")
	router.HandleFunc("/projects/{project_id}/templates/{email_template_id}", template.DeleteEmailTemplate).Methods("DELETE")
	router.HandleFunc("/projects/{project_id}/templates/{email_template_id}", template.UpdateEmailTemplate).Methods("UPDATE")
	router.HandleFunc("/projects/{project_id}/templates/{email_template_id}/_activate", template.ActivateEmailTemplate).Methods("GET")


	// subscription routes
	router.HandleFunc("/projects/{project_id}/subscriptions", subscription.GetSubscriptions).Methods("GET")
	router.HandleFunc("/projects/{project_id}/subscriptions", subscription.CreateSubscription).Methods("POST")
	router.HandleFunc("/projects/{project_id}/subscriptions/{subscription_id}", subscription.GetSubscription).Methods("GET")
	router.HandleFunc("/projects/{project_id}/subscriptions/{subscription_id}", subscription.DeleteSubscription).Methods("DELETE")
	router.HandleFunc("/projects/{project_id}/subscriptions/{subscription_id}", subscription.UpdateSubscription).Methods("UPDATE")
	router.HandleFunc("/projects/{project_id}/subscriptions/{subscription_id}/_token", subscription.GetSubscriptionToken).Methods("GET")
	router.HandleFunc("/unsubscribe", subscription.GetUnsubscribePage).Queries("token", "{token}").Methods("GET")
	router.HandleFunc("/unsubscribe", subscription.DeleteSubscriptionByToken).Queries("token", "{token}").Methods("DELETE")
	router.HandleFunc("/goodbye", subscription.GetGoodbyePage).Methods("GET")

	s.L.With(zap.String("bind", bind)).Info("listening")

	r := handlers.CORS(originsOk, headersOk, methodsOk)(router)
	if err := http.ListenAndServe(bind, r); err != nil {
		return err
	}

	return nil
}
