package server

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/statictask/newsletter/pkg/project"
)

var (
	originsOk = handlers.AllowedOrigins([]string{"http://luan.com"})
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

	// projects routes
	router.HandleFunc("/projects", project.GetProjects).Methods("GET")
	router.HandleFunc("/projects", project.CreateProject).Methods("POST")
	router.HandleFunc("/projects/{project_id}", project.GetProject).Methods("GET")
	router.HandleFunc("/projects/{project_id}", project.DeleteProject).Methods("DELETE")
	router.HandleFunc("/projects/{project_id}", project.UpdateProject).Methods("UPDATE")

	// subscription routes
	router.HandleFunc("/projects/{project_id}/subscriptions", project.GetProjectSubscriptions).Methods("GET")
	router.HandleFunc("/projects/{project_id}/subscriptions", project.CreateProjectSubscription).Methods("POST")
	router.HandleFunc("/projects/{project_id}/subscriptions/{subscription_id}", project.GetProjectSubscription).Methods("GET")
	router.HandleFunc("/projects/{project_id}/subscriptions/{subscription_id}", project.DeleteProjectSubscription).Methods("DELETE")
	router.HandleFunc("/projects/{project_id}/subscriptions/{subscription_id}", project.UpdateProjectSubscription).Methods("UPDATE")

	s.L.With(zap.String("bind", bind)).Info("listening")

	r := handlers.CORS(originsOk, headersOk, methodsOk)(router)
	if err := http.ListenAndServe(bind, r); err != nil {
		return err
	}

	return nil
}
