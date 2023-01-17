package template

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/internal/utils"
	"go.uber.org/zap"
)

// CreateEmailTemplate creates a new EmailTemplate record in the database
func CreateEmailTemplate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	controller := NewProjectEmailTemplates(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	et := New()
	if err = json.NewDecoder(r.Body).Decode(&et); err != nil {
		_log.Error("Failed decoding request body.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err = controller.Add(et); err != nil {
		_log.Error("Failed adding new EmailTemplate.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log.Info(
		"EmailTemplate created successfully",
		zap.Int64("email_template_id", et.ID),
	)

	utils.WriteJSONResponseData(w, http.StatusOK, et)
}

// GetEmailTemplate return a single subscription
func GetEmailTemplate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	controller := NewProjectEmailTemplates(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	etID, err := strconv.Atoi(params["email_template_id"])
	if err != nil {
		log.L.Error("Failed parsing email_template_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	id := int64(etID)
	_log = _log.With(zap.Int64("email_template_id", id))

	et, err := controller.Get(id)
	if err != nil {
		_log.Error("Failed getting EmailTemplate.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	_log.Info("EmailTemplate retrieved successfully.")
	utils.WriteJSONResponseData(w, http.StatusOK, et)
}

// GetEmailTemplates return all subscriptions related to a given project
func GetEmailTemplates(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	controller := NewProjectEmailTemplates(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	ets, err := controller.All()
	if err != nil {
		_log.Error("Failed loading EmailTemplates.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log.Info("EmailTemplates retrieved successfully.")
	utils.WriteJSONResponseData(w, http.StatusOK, ets)
}

// UpdateEmailTemplate updates an EmailTemplate entry from the project
func UpdateEmailTemplate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	controller := NewProjectEmailTemplates(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	etID, err := strconv.Atoi(params["email_template_id"])
	if err != nil {
		log.L.Error("Failed parsing email_template_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	id := int64(etID)
	_log = _log.With(zap.Int64("email_template_id", id))

	et, err := controller.Get(id)
	if err != nil {
		_log.Error("Failed getting EmailTemplate.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&et); err != nil {
		_log.Error("Failed decoding request body.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err := et.Update(); err != nil {
		_log.Error("Failed updating EmailTemplate.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log.Info("EmailTemplate updated successfully.")
	utils.WriteJSONResponseData(w, http.StatusOK, et)
}

// DeleteEmailTemplate deletes an EmailTemplate entry from the project
func DeleteEmailTemplate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	controller := NewProjectEmailTemplates(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	etID, err := strconv.Atoi(params["email_template_id"])
	if err != nil {
		log.L.Error("Failed parsing email_template_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	id := int64(etID)
	_log = _log.With(zap.Int64("email_template_id", id))

	et, err := controller.Get(id)
	if err != nil {
		_log.Error("Failed getting EmailTemplate.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err := et.Delete(); err != nil {
		_log.Error("Failed deleting EmailTemplate.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log.Info("EmailTemplate deleted successfully.")
	msg := "Email template deleted successfully."
	utils.WriteJSONResponseMessage(w, http.StatusNoContent, msg)
}

// ActivateEmplateTemplate deactivates all the other email templates for
// this project and enables only the given email template
func ActivateEmailTemplate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	controller := NewProjectEmailTemplates(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	etID, err := strconv.Atoi(params["email_template_id"])
	if err != nil {
		log.L.Error("Failed parsing email_template_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	id := int64(etID)
	_log = _log.With(zap.Int64("email_template_id", id))

	newActiveEt, err := controller.Get(id)
	if err != nil {
		_log.Error("Failed getting EmailTemplate.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	currentActiveEt, err := controller.GetActive()
	if err != nil {
		_log.Error("Failed getting current active EmailTemplate.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	if currentActiveEt != nil {
		if err := currentActiveEt.Deactivate(); err != nil {
			_log.Error("Failed deactivating the current active EmailTemplate.", zap.Error(err))
			utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
			return
		}
	}

	if err := newActiveEt.Activate(); err != nil {
		_log.Error("Failed activating EmailTemplate.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log.Info("EmailTemplate activated successfully.")
	msg := "Email template activated successfully."
	utils.WriteJSONResponseMessage(w, http.StatusOK, msg)
}
