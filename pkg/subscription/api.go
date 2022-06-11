package subscription

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/statictask/newsletter/internal/log"
	"go.uber.org/zap"
)

// response format
type response struct {
	ID      interface{} `json:"id,omitempty"`
	Message string      `json:"message,omitempty"`
}

// GetSubscriptions will return all subscriptions
func GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	subscriptions := NewSubscriptions()

	allSubscriptions, err := subscriptions.All()
	if err != nil {
		log.L.Fatal("unable to load subscriptions", zap.Error(err))
	}

	json.NewEncoder(w).Encode(allSubscriptions)
}

// CreateSubscription will create a subscription
func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	subscription := New()

	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		log.L.Fatal("unable to decode request body", zap.Error(err))
	}

	if err := subscription.Create(); err != nil {
		log.L.Fatal("unable to create a subscription", zap.Error(err))
	}

	json.NewEncoder(w).Encode(subscription)
}

// GetSubscription will return a single subscription by its ID
func GetSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	params := mux.Vars(r)
	subscriptionID := params["subscription_id"]

	subscription, err := loadSubscription(subscriptionID)
	if err != nil {
		log.L.Fatal("unable to load subscription", zap.Error(err))
	}

	json.NewEncoder(w).Encode(subscription)
}

// UpdateSubscription update subscription's details
func UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	subscriptionID := params["subscription_id"]

	subscription, err := loadSubscription(subscriptionID)
	if err != nil {
		log.L.Fatal("unable to load subscription", zap.Error(err))
	}

	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		log.L.Fatal("unable to decode the request body", zap.Error(err))
	}

	if err := subscription.Update(); err != nil {
		log.L.Fatal("unable to update subscription", zap.Error(err))
	}

	log.L.Info("subscription updated successfully", zap.Int64("subscriptionID", subscription.ID))

	msg := "subscription updated successfully"
	res := response{
		ID:      int64(subscription.ID),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// DeleteSubscription deletes a subscription from database
func DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")

	params := mux.Vars(r)
	subscriptionID := params["subscription_id"]

	subscription, err := loadSubscription(subscriptionID)
	if err != nil {
		log.L.Fatal("unable to load subscription", zap.Error(err))
	}

	if err := subscription.Delete(); err != nil {
		log.L.Fatal("unable to delete subscription", zap.Error(err))
	}

	log.L.Info("subscription deleted successfully", zap.String("subscriptionID", subscriptionID))

	id, err := strconv.Atoi(subscriptionID)
	if err != nil {
		log.L.Fatal("unable to convert subscription_id to int", zap.Error(err))
	}

	msg := fmt.Sprintf("subscription deleted successfully")
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}
