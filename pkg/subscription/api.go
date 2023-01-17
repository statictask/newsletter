package subscription

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"html/template"

	"github.com/gorilla/mux"
	"github.com/statictask/newsletter/internal/log"
	"github.com/statictask/newsletter/internal/utils"
	"go.uber.org/zap"
)

// CreateSubscription creates a subscription entry in the project
func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	controller := NewProjectSubscriptions(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	s := New()
	if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
		_log.Error("Failed decoding request body.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err = controller.Add(s); err != nil {
		_log.Error("Failed adding new Subscription.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log.Info(
		"Subscription created successfully.",
		zap.Int64("subscription_id", s.ID),
	)

	utils.WriteJSONResponseData(w, http.StatusOK, s)
}

// GetProjectSubscription return a single subscription
func GetSubscription(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	controller := NewProjectSubscriptions(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	subscriptionID, err := strconv.Atoi(params["subscription_id"])
	if err != nil {
		log.L.Error("Failed parsing subscription_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	id := int64(subscriptionID)
	_log = _log.With(zap.Int64("subscription_id", id))

	s, err := controller.Get(id)
	if err != nil {
		_log.Error("Failed getting Subscription.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	_log.Info("Subscription retrieved successfully.")
	utils.WriteJSONResponseData(w, http.StatusOK, s)
}

// GetProjectSubscriptions return all subscriptions related to a given project
func GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	controller := NewProjectSubscriptions(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	subscriptions, err := controller.All()
	if err != nil {
		_log.Error("Failed loading Subscriptions.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log.Info("Subscriptions retrieved successfully.")
	utils.WriteJSONResponseData(w, http.StatusOK, subscriptions)
}

// UpdateSubscription updates a subscription entry from the project
func UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	controller := NewProjectSubscriptions(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	subscriptionID, err := strconv.Atoi(params["subscription_id"])
	if err != nil {
		log.L.Error("Failed parsing subscription_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	id := int64(subscriptionID)
	_log = _log.With(zap.Int64("subscription_id", id))

	s, err := controller.Get(id)
	if err != nil {
		_log.Error("Failed getting Subscription.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
		_log.Error("Failed decoding request body.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err = s.Update(); err != nil {
		_log.Error("Failed updating Subscription.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log.Info("Subscription updated successfully.")
	utils.WriteJSONResponseData(w, http.StatusOK, s)
}

// DeleteSubscription deletes a subscription entry from the project
func DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	controller := NewProjectSubscriptions(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	subscriptionID, err := strconv.Atoi(params["subscription_id"])
	if err != nil {
		log.L.Error("Failed parsing subscription_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	id := int64(subscriptionID)
	_log = _log.With(zap.Int64("subscription_id", id))

	s, err := controller.Get(id)
	if err != nil {
		_log.Error("Failed getting Subscription.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err = s.Delete(); err != nil {
		_log.Error("Failed deleting Subscription.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log.Info("Subscription deleted successfully.")
	msg := "Subscription deleted successfully."
	utils.WriteJSONResponseMessage(w, http.StatusNoContent, msg)
}

// GetUnsubscribePage builds an HTML response with an unsbscribe page
func GetUnsubscribePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "text/html")

	token := r.URL.Query().Get("token")

	s, err := Decrypt(token)
	if err != nil {
		tmpl :=	template.Must(template.ParseFiles("static/404/index.html"))
        	tmpl.Execute(w, nil)

		return
	}

	// Parse the template file
        tmpl := template.Must(template.ParseFiles("static/unsubscribe/index.html"))

	data := map[string]interface{} {
		"token": token,
		"email": s.Email,
	}

        // Execute the template and write it to the response writer
        tmpl.Execute(w, data)
}

// DeleteUnsubscribe builds an HTML response with an unsbscribe page
func DeleteSubscriptionByToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	s, err := Decrypt(token)
	if err != nil {
		log.L.Error("Failed to decrypt subscription token.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusNotFound, err)
		return
	}

	_log := log.L.With(zap.Int64("project_id", s.ProjectID), zap.Int64("subscription_id", s.ID))
	
	if err = s.Delete(); err != nil {
		log.L.Error("Failed deleting subscription by token.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log.Info("Subscription deleted successfully")
	msg := fmt.Sprintf("subscription %d deleted successfully", s.ID)
	utils.WriteJSONResponseMessage(w, http.StatusNoContent, msg)
}

// GetGoodbyePage builds an HTML response with an unsbscribe page
func GetGoodbyePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "text/html")

	// Parse the template file
        tmpl := template.Must(template.ParseFiles("static/goodbye/index.html"))

        // Execute the template and write it to the response writer
        tmpl.Execute(w, nil)
}

// GetSubscriptionToken return a single subscription token
func GetSubscriptionToken(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("Failed parsing project_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	controller := NewProjectSubscriptions(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	subscriptionID, err := strconv.Atoi(params["subscription_id"])
	if err != nil {
		log.L.Error("Failed parsing subscription_id.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	id := int64(subscriptionID)
	_log = _log.With(zap.Int64("subscription_id", id))

	s, err := controller.Get(id)
	if err != nil {
		_log.Error("Failed getting Subscription.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	if s == nil {
		err = fmt.Errorf("Subscription not found.")
		_log.Error("Failed getting Subscription.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusBadRequest, err)
		return
	}

	token, err := s.Encrypt()
	if err != nil {
		_log.Error("Failed generating Subscription token.", zap.Error(err))
		utils.WriteJSONResponseError(w, http.StatusInternalServerError, err)
		return
	}

	_log.Info("Subscription token generated successfully.")
	data := map[string]string {
		"token": token,
	}

	fmt.Println(token)

	utils.WriteJSONResponseData(w, http.StatusOK, &data)
}
