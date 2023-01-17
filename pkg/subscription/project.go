package subscription

import "fmt"

// ProjectSubscriptions is the entity used for controlling
// interactions with many subscriptions in the database
type ProjectSubscriptions struct {
	projectID int64
}

// NewProjectSubscriptions returns a ProjectSubscriptions controller
func NewProjectSubscriptions(projectID int64) *ProjectSubscriptions {
	return &ProjectSubscriptions{projectID}
}

// All returns all the subscriptions registered in the database for
// a given project
func (ps *ProjectSubscriptions) All() ([]*Subscription, error) {
	subscriptions, err := getSubscriptions(ps.projectID)
	if err != nil {
		return subscriptions, fmt.Errorf("unable to get subscriptions: %v", err)
	}

	return subscriptions, nil
}

// Get a single subscription based on the project and the subscriptionID
func (ps *ProjectSubscriptions) Get(subscriptionID int64) (*Subscription, error) {
	subscription, err := getSubscription(ps.projectID, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("unable to get this subscription: %v", err)
	}

	return subscription, nil
}

// Delete deletes a subscription based on its ID
func (ps *ProjectSubscriptions) Delete(subscriptionID int64) error {
	if err := deleteSubscription(subscriptionID, ps.projectID); err != nil {
		return fmt.Errorf("unable to delete subscription: %v", err)
	}

	return nil
}

// Add creates a new entry in the project's subscriptions
func (ps *ProjectSubscriptions) Add(s *Subscription) error {
	// make sure the subscription has the corred ProjectID before adding
	s.ProjectID = ps.projectID

	if err := s.Create(); err != nil {
		return err
	}

	return nil
}
