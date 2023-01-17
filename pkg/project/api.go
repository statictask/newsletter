package project

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/internal/utils"
	"go.uber.org/zap"
)

// CreateProject will create a project
func CreateProject(w http.ResponseWriter, r *http.Request) {
	project := New()

	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		log.L.Error("Failed decoding request body,", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err := project.Create(); err != nil {
		log.L.Error("Failed creating a new Project.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log := log.L.With(zap.Int64("project_id", project.ID))

	et, err := project.EmailTemplates().CreateDefault()
	if err != nil {
		_log.Error("Failed creating the default EmailTemplate for the new Project")
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	if err := et.Activate(); err != nil {
		_log.Error("Failed activating default EmailTemplate for the new Project")
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log.Info("Project created successfully.")
	utils.WriteJSONResponseData(w, http.StatusOK, project)
}

// GetProject will return a single project by its ID
func GetProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	projects := NewProjects()
	project, err := projects.Get(int64(id))
	if err != nil {
		log.L.Error("Failed loading Project.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log := log.L.With(zap.Int("project_id", id))

	if project == nil {
		err := fmt.Errorf("Project %d not found.", id)
		_log.Error("Failed loading Project.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusNotFound, err)
		return
	}

	_log.Info("Project loaded successfully")
	utils.WriteJSONResponseData(w, http.StatusOK, project)
}

// UpdateProject update project's details
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	projects := NewProjects()
	project, err := projects.Get(int64(id))
	if err != nil {
		log.L.Error("Failed loading Project.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log := log.L.With(zap.Int("project_id", id))

	if project == nil {
		err := fmt.Errorf("Project %d not found.", id)
		_log.Error("Failed loading Project.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusNotFound, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		_log.Error("Failed decoding the request body.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err := project.Update(); err != nil {
		_log.Error("Failed updating project.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log.Info("Project updated successfully")
	utils.WriteJSONResponseData(w, http.StatusOK, project)
}

// DeleteProject deletes a project from database
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	projects := NewProjects()
	project, err := projects.Get(int64(id))
	if err != nil {
		log.L.Error("Failed loading Project.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log := log.L.With(zap.Int("project_id", id))

	if project == nil {
		err := fmt.Errorf("Project %d not found.", id)
		_log.Error("Failed loading Project.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusNotFound, err)
		return
	}

	if err = project.Delete(); err != nil {
		log.L.Error("Failed deleting Project.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	msg := "Project deleted successfully."
	_log.Info(msg)
	utils.WriteJSONResponseMessage(w, http.StatusNoContent, msg)
}
