package subscription

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

// CreateProjectSubscription creates a subscription entry in the project
func CreateProjectSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var err error
	var res utils.HTTPResponse
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("unable to parse project id", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	controller := NewProjectSubscriptions(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	s := New()
	if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
		_log.Error("unable to decode the request body", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = controller.Add(s); err != nil {
		_log.Error("unable to add the new subscription", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	_log.Info(
		"subscription created successfully",
		zap.Int64("subscription_id", s.ID),
	)

	res = utils.HTTPResponse{
		Status:     "OK",
		StatusCode: 200,
		Data:       s,
	}

	json.NewEncoder(w).Encode(res)
}

// GetProjectSubscription return a single subscription
func GetProjectSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var err error
	var res utils.HTTPResponse
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("unable to parse project id", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	controller := NewProjectSubscriptions(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	subscriptionID, err := strconv.Atoi(params["subscription_id"])
	if err != nil {
		_log.Error("unable to parse subscription id", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	id := int64(subscriptionID)
	_log = _log.With(zap.Int64("subscription_id", id))

	subscription, err := controller.Get(id)
	if err != nil {
		_log.Error("unable to get subscription", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	_log.Info("subscription retrieved successfully")
	res = utils.HTTPResponse{
		Status:     "OK",
		StatusCode: 200,
		Data:       subscription,
	}

	json.NewEncoder(w).Encode(res)
}

// GetProjectSubscriptions return all subscriptions related to a given project
func GetProjectSubscriptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var err error
	var res utils.HTTPResponse
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("unable to parse project id", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	controller := NewProjectSubscriptions(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	subscriptions, err := controller.All()
	if err != nil {
		_log.Error("unable to load project subscriptions", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	_log.Info("subscriptions retrieved successfully")
	res = utils.HTTPResponse{
		Status:     "OK",
		StatusCode: 200,
		Data:       subscriptions,
	}

	json.NewEncoder(w).Encode(res)
}

// UpdateProjectSubscription updates a subscription entry from the project
func UpdateProjectSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	var err error
	var res utils.HTTPResponse
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("unable to parse project id", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	controller := NewProjectSubscriptions(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	subscriptionID, err := strconv.Atoi(params["subscription_id"])
	if err != nil {
		_log.Error("unable to parse subscription id", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	id := int64(subscriptionID)
	_log = _log.With(zap.Int64("subscription_id", id))

	subscription, err := controller.Get(id)
	if err != nil {
		_log.Error("unable to get subscription", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		_log.Error("unable to decode the request body", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	if err = subscription.Update(); err != nil {
		_log.Error("unable to update subscription", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	_log.Info("subscription updated successfully")
	res = utils.HTTPResponse{
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
	var res utils.HTTPResponse
	params := mux.Vars(r)

	projectID, err := strconv.Atoi(params["project_id"])
	if err != nil {
		log.L.Error("unable to parse project id", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Not Found",
			StatusCode: 404,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	controller := NewProjectSubscriptions(int64(projectID))
	_log := log.L.With(zap.Int64("project_id", int64(projectID)))

	subscriptionID, err := strconv.Atoi(params["subscription_id"])
	if err != nil {
		_log.Error("unable to parse subscription id", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Bad Request",
			StatusCode: 400,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	id := int64(subscriptionID)
	_log = _log.With(zap.Int64("subscription_id", id))

	if err = controller.Delete(id); err != nil {
		_log.Error("unable to delete subscription", zap.Error(err))
		res = utils.HTTPResponse{
			Status:     "Internal Server Error",
			StatusCode: 500,
			Data:       utils.HTTPErrorMessage{err.Error()},
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	_log.Info("subscription deleted successfully")

	msg := fmt.Sprintf("subscription %d deleted successfully", subscriptionID)
	res = utils.HTTPResponse{
		Status:     "OK",
		StatusCode: 200,
		Data:       utils.HTTPInfoMessage{msg},
	}

	json.NewEncoder(w).Encode(res)
}
