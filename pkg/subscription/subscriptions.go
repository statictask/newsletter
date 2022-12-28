package subscription

import "fmt"

// Subscriptions is the entity used for controlling
// interactions with many subscriptions in the database
type Subscriptions struct {
	projectID int64
}

// NewSubscriptions returns a Subscriptions controller
func NewSubscriptions(projectID int64) *Subscriptions {
	return &Subscriptions{projectID}
}

// All returns all the subscriptions registered in the database for
// a given project
func (ss *Subscriptions) All() ([]*Subscription, error) {
	exp := fmt.Sprintf("project_id=%d", ss.projectID)

	subscriptions, err := getSubscriptionsWhere(exp)
	if err != nil {
		return subscriptions, fmt.Errorf("unable to get subscriptions: %v", err)
	}

	return subscriptions, nil
}

// Get returns a single subscription according to its ID
func (ss *Subscriptions) Get(subscriptionID int64) (*Subscription, error) {
	exp := fmt.Sprintf("subscription_id=%d AND project_id=%d", subscriptionID, ss.projectID)

	return getSubscriptionWhere(exp)
}

// Where return many subscriptions according to a map of attrs
func (ss *Subscriptions) Where(exp string) ([]*Subscription, error) {
	exp = fmt.Sprintf("%s AND project_id=%d", exp, ss.projectID)

	subscriptions, err := getSubscriptionsWhere(exp)
	if err != nil {
		return nil, fmt.Errorf("unable to get subscriptions: %v", err)
	}

	return subscriptions, nil
}

// Delete deletes a subscription based on its ID
func (ss *Subscriptions) Delete(subscriptionID int64) error {
	if err := deleteSubscription(subscriptionID, ss.projectID); err != nil {
		return fmt.Errorf("unable to delete subscription: %v", err)
	}

	return nil
}

// Add creates a new entry in the project's subscriptions
// the function creates a new subscription entry in the database
func (ss *Subscriptions) Add(s *Subscription) error {
	// make sure the subscription has the corred ProjectID before adding
	s.ProjectID = ss.projectID

	if err := s.Create(); err != nil {
		return err
	}

	return nil
}
