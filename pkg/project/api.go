package project

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/pkg/subscription"
	"go.uber.org/zap"
)

// response format
type response struct {
	ID      interface{} `json:"id,omitempty"`
	Message string      `json:"message,omitempty"`
}

// GetProjects will return all projects
func GetProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	projects := NewProjects()

	allProjects, err := projects.All()
	if err != nil {
		log.L.Fatal("unable to load projects", zap.Error(err))
	}

	json.NewEncoder(w).Encode(allProjects)
}

// CreateProject will create a project
func CreateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	project := New()

	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		log.L.Fatal("unable to decode request body", zap.Error(err))
	}

	if err := project.Create(); err != nil {
		log.L.Fatal("unable to create a project", zap.Error(err))
	}

	json.NewEncoder(w).Encode(project)
}

// GetProject will return a single project by its ID
func GetProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Fatal("unable to load project", zap.Error(err))
	}

	json.NewEncoder(w).Encode(project)
}

// UpdateProject update project's details
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Fatal("unable to load project", zap.Error(err))
	}

	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		log.L.Fatal("unable to decode the request body", zap.Error(err))
	}

	if err := project.Update(); err != nil {
		log.L.Fatal("unable to update project", zap.Error(err))
	}

	log.L.Info("project updated successfully", zap.Int64("projectID", project.ID))

	msg := "project updated successfully"
	res := response{
		ID:      int64(project.ID),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// DeleteProject deletes a project from database
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Fatal("unable to load project", zap.Error(err))
	}

	if err := project.Delete(); err != nil {
		log.L.Fatal("unable to delete project", zap.Error(err))
	}

	log.L.Info("project deleted successfully", zap.String("projectID", projectID))

	id, err := strconv.Atoi(projectID)
	if err != nil {
		log.L.Fatal("unable to convert project_id to int", zap.Error(err))
	}

	msg := fmt.Sprintf("project deleted successfully")
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// CreateProjectSubscription creates a subscription entry in the project
func CreateProjectSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Fatal("unable to load project", zap.Error(err))
	}

	s := subscription.New()
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		log.L.Fatal("unable to decode the request body", zap.Error(err))
	}

	if err := project.Subscriptions().Add(s); err != nil {
		log.L.Fatal("unable to add the new subscription", zap.Error(err))
	}

	msg := fmt.Sprintf("subscription %d added successfully", s.ID)

	res := response{
		ID:      int64(s.ID),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// GetProjectSubscriptions return all subscriptions related to a given project
func GetProjectSubscriptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Fatal("unable to load project", zap.Error(err))
	}

	subscriptions, err := project.Subscriptions().All()
	if err != nil {
		log.L.Fatal("unable to load project subscriptions", zap.Error(err))
	}

	json.NewEncoder(w).Encode(subscriptions)
}

// GetProjectSubscription return a specific subscription related to a given project
func GetProjectSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Fatal("unable to load project", zap.Error(err))
	}

	subscriptionID, err := castID(params["subscription_id"])
	if err != nil {
		log.L.Fatal("unable to parse subscription id", zap.Error(err))
	}

	subscription, err := project.Subscriptions().Get(subscriptionID)
	if err != nil {
		log.L.Fatal("unable to load project subscriptions", zap.Error(err))
	}

	json.NewEncoder(w).Encode(subscription)
}

// UpdateProjectSubscription updates a subscription entry from the project
func UpdateProjectSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Fatal("unable to load project", zap.Error(err))
	}

	subscriptionID, err := castID(params["subscription_id"])
	if err != nil {
		log.L.Fatal("unable to parse subscription id", zap.Error(err))
	}

	subscription, err := project.Subscriptions().Get(subscriptionID)
	if err != nil {
		log.L.Fatal("unable to get subscription", zap.Error(err))
	}

	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		log.L.Fatal("unable to decode the request body", zap.Error(err))
	}

	if err := subscription.Update(); err != nil {
		log.L.Fatal("unable to update subscription", zap.Error(err))
	}

	log.L.Info("subscription updated successfully", zap.Int64("ID", subscription.ID))

	msg := "subscription updated successfully"
	res := response{
		ID:      int64(project.ID),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// DeleteProjectSubscription deletes a subscription entry from the project
func DeleteProjectSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Fatal("unable to load project", zap.Error(err))
	}

	subscriptionID, err := castID(params["subscription_id"])
	if err != nil {
		log.L.Fatal("unable to parse subscription id", zap.Error(err))
	}

	if err := project.Subscriptions().Delete(subscriptionID); err != nil {
		log.L.Fatal("unable to delete subscription", zap.Error(err))
	}

	msg := fmt.Sprintf("subscription %d removed successfully", subscriptionID)

	res := response{
		ID:      int64(subscriptionID),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}
