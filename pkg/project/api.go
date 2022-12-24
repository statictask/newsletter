package project

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/pkg/subscription"
	"go.uber.org/zap"
)

// response format
type response struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}

// response error format
type errorMessage struct {
	Error string `json:"error"`
}

// response info format
type infoMessage struct {
	Message string `json:"message"`
}

// GetProjects will return all projects
func GetProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var res response
	projects := NewProjects()

	allProjects, err := projects.All()
	if err != nil {
		log.L.Error("unable to load projects", zap.Error(err))
		res = response{
			Status:     "Service Unavailable",
			StatusCode: 503,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("projects retrieved successfully")

	res = response{
		Status:     "OK",
		StatusCode: 200,
		Data:       allProjects,
	}

	json.NewEncoder(w).Encode(res)
}

// CreateProject will create a project
func CreateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var err error
	var res response
	project := New()

	if err = json.NewDecoder(r.Body).Decode(&project); err != nil {
		log.L.Error("unable to decode request body", zap.Error(err))
		res = response{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = project.Create(); err != nil {
		log.L.Error("unable to create a project", zap.Error(err))
		res = response{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("project created successfully", zap.Int64("projectID", project.ID))
	res = response{
		Status:     "OK",
		StatusCode: 200,
		Data:       project,
	}

	json.NewEncoder(w).Encode(res)
}

// GetProject will return a single project by its ID
func GetProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var res response
	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Error("unable to load project", zap.Error(err))
		res = response{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("project retrieved successfully", zap.String("projectID", projectID))
	res = response{
		Status:     "OK",
		StatusCode: 200,
		Data:       project,
	}

	json.NewEncoder(w).Encode(res)
}

// UpdateProject update project's details
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var err error
	var res response
	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Error("unable to load project", zap.Error(err))
		res = response{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&project); err != nil {
		log.L.Error("unable to decode the request body", zap.Error(err))
		res = response{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = project.Update(); err != nil {
		log.L.Error("unable to update project", zap.Error(err))
		res = response{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("project updated successfully", zap.Int64("projectID", project.ID))
	res = response{
		Status:     "OK",
		StatusCode: 200,
		Data:       project,
	}

	json.NewEncoder(w).Encode(res)
}

// DeleteProject deletes a project from database
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var err error
	var res response
	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Error("unable to load project", zap.Error(err))
		res = response{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = project.Delete(); err != nil {
		log.L.Error("unable to delete project", zap.Error(err))
		res = response{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("project deleted successfully", zap.String("projectID", projectID))

	msg := fmt.Sprintf("project %s deleted successfully", projectID)
	res = response{
		Status:     "OK",
		StatusCode: 200,
		Data:       infoMessage{msg},
	}

	json.NewEncoder(w).Encode(res)
}

// CreateProjectSubscription creates a subscription entry in the project
func CreateProjectSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var err error
	var res response
	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Error("unable to load project", zap.Error(err))
		res = response{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	s := subscription.New()
	if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
		log.L.Error("unable to decode the request body", zap.Error(err))
		res = response{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = project.Subscriptions().Add(s); err != nil {
		log.L.Error("unable to add the new subscription", zap.Error(err))
		res = response{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("subscription created successfully", zap.String("projectID", projectID), zap.Int64("subscriptionID", s.ID))
	res = response{
		Status:     "OK",
		StatusCode: 200,
		Data:       s,
	}

	json.NewEncoder(w).Encode(res)
}

// GetProjectSubscriptions return all subscriptions related to a given project
func GetProjectSubscriptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var err error
	var res response
	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Error("unable to load project", zap.Error(err))
		res = response{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	subscriptions, err := project.Subscriptions().All()
	if err != nil {
		log.L.Error("unable to load project subscriptions", zap.Error(err))
		res = response{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("subscriptions retrieved successfully", zap.String("projectID", projectID))
	res = response{
		Status:     "OK",
		StatusCode: 200,
		Data:       subscriptions,
	}

	json.NewEncoder(w).Encode(res)
}

// GetProjectSubscription return a specific subscription related to a given project
func GetProjectSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var err error
	var res response
	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Error("unable to load project", zap.Error(err))
		res = response{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	subscriptionID, err := castID(params["subscription_id"])
	if err != nil {
		log.L.Error("unable to parse subscription id", zap.Error(err))
		res = response{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	subscription, err := project.Subscriptions().Get(subscriptionID)
	if err != nil {
		log.L.Error("unable to load project subscription", zap.Error(err))
		res = response{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("subscription retrieved successfully", zap.String("projectID", projectID), zap.Int64("subscriptionID", subscriptionID))
	res = response{
		Status:     "OK",
		StatusCode: 200,
		Data:       subscription,
	}

	json.NewEncoder(w).Encode(res)
}

// UpdateProjectSubscription updates a subscription entry from the project
func UpdateProjectSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var err error
	var res response
	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Error("unable to load project", zap.Error(err))
		res = response{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	subscriptionID, err := castID(params["subscription_id"])
	if err != nil {
		log.L.Error("unable to parse subscription id", zap.Error(err))
		res = response{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	subscription, err := project.Subscriptions().Get(subscriptionID)
	if err != nil {
		log.L.Error("unable to get subscription", zap.Error(err))
		res = response{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		log.L.Error("unable to decode the request body", zap.Error(err))
		res = response{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = subscription.Update(); err != nil {
		log.L.Error("unable to update subscription", zap.Error(err))
		res = response{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("subscription updated successfully", zap.String("projectID", projectID), zap.Int64("ID", subscription.ID))
	res = response{
		Status:     "OK",
		StatusCode: 200,
		Data:       subscription,
	}

	json.NewEncoder(w).Encode(res)
}

// DeleteProjectSubscription deletes a subscription entry from the project
func DeleteProjectSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var err error
	var res response
	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Error("unable to load project", zap.Error(err))
		res = response{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	subscriptionID, err := castID(params["subscription_id"])
	if err != nil {
		log.L.Error("unable to parse subscription id", zap.Error(err))
		res = response{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = project.Subscriptions().Delete(subscriptionID); err != nil {
		log.L.Error("unable to delete subscription", zap.Error(err))
		res = response{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       errorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("subscription deleted successfully", zap.Int64("subscriptionID", subscriptionID))

	msg := fmt.Sprintf("subscription %d deleted successfully", subscriptionID)
	res = response{
		Status:     "OK",
		StatusCode: 200,
		Data:       infoMessage{msg},
	}

	json.NewEncoder(w).Encode(res)
}
