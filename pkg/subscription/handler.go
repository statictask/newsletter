package subscription

import (
	"database/sql"
	"fmt"

	"github.com/statictask/newsletter/internal/database"
)

// insertSubscription inserts a subscription in the database
func insertSubscription(s *Subscription) error {
	query := `
		INSERT INTO subscriptions (
		  project_id,
		  email
	  	)
		VALUES (
		  $1,
		  $2
	  	)
		RETURNING
		  subscription_id,
		  project_id,
		  email,
		  created_at,
		  updated_at
	`

	savedSubscription, err := scanSubscription(query, s.ProjectID, s.Email)
	if err != nil {
		return err
	}

	*s = *savedSubscription

	return nil
}

// getSubscriptions returns all subscriptions in the database
func getSubscriptions(projectID int64) ([]*Subscription, error) {
	query := `
		SELECT
		  subscription_id,
		  project_id,
		  email,
		  created_at,
		  updated_at
		FROM
		  subscriptions
		WHERE
		  project_id = $1
	`

	return scanSubscriptions(query, projectID)
}

// getProjectSubscription returns a single subscription that match both
// subscription and project id
func getSubscription(projectID, subscriptionID int64) (*Subscription, error) {
	query := `
		SELECT
		  subscription_id,
		  project_id,
		  email,
		  created_at,
		  updated_at
		FROM
		  subscriptions
		WHERE
		  project_id = $1
		  subscription_id = $2
	`

	return scanSubscription(query, projectID, subscriptionID)
}

// updateSubscription updates a subscription in the database
func updateSubscription(s *Subscription) error {
	// only allows updates to the email field, the other fields are immutable
	query := `UPDATE subscriptions SET email=$1 WHERE subscription_id=$2`

	if err := database.Exec(query, s.Email, s.ID); err != nil {
		return fmt.Errorf("failed updating subscription: %v", err)
	}

	return nil
}

// deleteSubscription deletes a subscription from database
func deleteSubscription(subscriptionID, projectID int64) error {
	query := `DELETE FROM subscriptions WHERE subscription_id=$1 AND project_id=$2`

	if err := database.Exec(query, subscriptionID, projectID); err != nil {
		return fmt.Errorf("failed deleting subscription: %v", err)
	}

	return nil
}

// scanSubscription returns a single subscription that matches the given query
func scanSubscription(query string, params ...interface{}) (*Subscription, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	row := db.QueryRow(query, params...)
	s := New()

	if err := row.Scan(&s.ID, &s.ProjectID, &s.Email, &s.CreatedAt, &s.UpdatedAt); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("unable to scan subscription row: %v", err)
		}

		return nil, nil
	}

	return s, nil
}

// scanSubscriptions returns multiple subscriptions that match the given query
func scanSubscriptions(query string, params ...interface{}) ([]*Subscription, error) {
	var subscriptions []*Subscription

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(query, params...)
	if err != nil {
		return subscriptions, fmt.Errorf("unable to execute `%s`: %v", query, err)
	}

	defer rows.Close()

	for rows.Next() {
		s := New()

		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Email, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return subscriptions, fmt.Errorf("unable to scan a subscription row: %v", err)
		}

		subscriptions = append(subscriptions, s)
	}

	return subscriptions, nil
}
