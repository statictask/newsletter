package subscription

import "fmt"

type Subscription struct {
	ID    int64  `json:"subscription_id"`
	Email string `json:"email"`
}

// New returns an empty Subscription
func New() *Subscription {
	return &Subscription{}
}

// Create the subscription in the database
func (s *Subscription) Create() error {
	if err := insertSubscription(s); err != nil {
		return fmt.Errorf("unable to create subscription: %v", err)
	}

	return nil
}

// Update the subscription in the database
func (s *Subscription) Update() error {
	if err := updateSubscription(s); err != nil {
		return fmt.Errorf("unable to update subscription: %v", err)
	}

	return nil
}

// Delete the subscription from the database
func (s *Subscription) Delete() error {
	if err := deleteSubscription(s.ID); err != nil {
		return fmt.Errorf("unable to delete subscription: %v", err)
	}

	return nil
}
