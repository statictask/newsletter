package project

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/internal/utils"
	"go.uber.org/zap"
)

// CreateProject will create a project
func CreateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var err error
	var res utils.HTTPResponse
	project := New()

	if err = json.NewDecoder(r.Body).Decode(&project); err != nil {
		log.L.Error("unable to decode request body", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = project.Create(); err != nil {
		log.L.Error("unable to create a project", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("project created successfully", zap.Int64("projectID", project.ID))
	res = utils.HTTPResponse{
		Status:     "OK",
		StatusCode: 200,
		Data:       project,
	}

	json.NewEncoder(w).Encode(res)
}

// GetProject will return a single project by its ID
func GetProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var res utils.HTTPResponse
	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Error("unable to load project", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("project retrieved successfully", zap.String("projectID", projectID))
	res = utils.HTTPResponse{
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
	var res utils.HTTPResponse
	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Error("unable to load project", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&project); err != nil {
		log.L.Error("unable to decode the request body", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = project.Update(); err != nil {
		log.L.Error("unable to update project", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("project updated successfully", zap.Int64("projectID", project.ID))
	res = utils.HTTPResponse{
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
	var res utils.HTTPResponse
	params := mux.Vars(r)
	projectID := params["project_id"]

	project, err := loadProject(projectID)
	if err != nil {
		log.L.Error("unable to load project", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = project.Delete(); err != nil {
		log.L.Error("unable to delete project", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	log.L.Info("project deleted successfully", zap.String("projectID", projectID))

	msg := fmt.Sprintf("project %s deleted successfully", projectID)
	res = utils.HTTPResponse{
		Status:     "OK",
		StatusCode: 200,
		Data:       utils.HTTPInfoMessage{msg},
	}

	json.NewEncoder(w).Encode(res)
}
