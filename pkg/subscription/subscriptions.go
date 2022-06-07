package subscription

import "fmt"

// Subscriptions is the entity used for controlling
// interactions with many subscriptions in the database
type Subscriptions struct{}

// NewSubscriptions returns a Subscriptions controller
func NewSubscriptions() *Subscriptions {
	return &Subscriptions{}
}

// All returns all the subscriptions registered in the database
func (ss *Subscriptions) All() ([]*Subscription, error) {
	subscriptions, err := getSubscriptionsWhere("")
	if err != nil {
		return subscriptions, fmt.Errorf("unable to get subscriptions: %v", err)
	}

	return subscriptions, nil
}

// Get returns a single subscription according to its ID
func (ss *Subscriptions) Get(subscriptionID int64) (*Subscription, error) {
	exp := fmt.Sprintf("subscription_id=%v", subscriptionID)

	return getSubscriptionWhere(exp)
}

// Where return many subscriptions according to a map of attrs
func (ss *Subscriptions) Where(exp string) ([]*Subscription, error) {
	subscriptions, err := getSubscriptionsWhere(exp)
	if err != nil {
		return nil, fmt.Errorf("unable to get subscriptions: %v", err)
	}

	return subscriptions, nil
}

// Delete deletes a subscription based on its ID
func (ss *Subscriptions) Delete(subscriptionID int64) error {
	if err := deleteSubscription(subscriptionID); err != nil {
		return fmt.Errorf("unable to delete subscription: %v", err)
	}

	return nil
}
