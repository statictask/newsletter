package subscription

import (
	"fmt"

	"github.com/statictask/newsletter/internal/database"
	"github.com/statictask/newsletter/internal/log"
	"go.uber.org/zap"
)

// insertSubscription inserts a subscription in the database
func insertSubscription(s *Subscription) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `INSERT INTO subscriptions (project_id,email) VALUES ($1, $2) RETURNING subscription_id,created_at,updated_at`

	if err := db.QueryRow(
		sqlStatement,
		s.ProjectID,
		s.Email,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt); err != nil {
		log.L.Fatal("unable to execute the query", zap.Error(err))
	}

	log.L.Info(
		"created subscription record",
		zap.Int64("subscription_id", s.ID),
		zap.Int64("product_id", s.ProjectID),
	)

	return nil
}

// getSubscriptionWhere return a single row that matches a given expression
func getSubscriptionWhere(expression string) (*Subscription, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := "SELECT subscription_id,project_id,email,created_at,updated_at FROM subscriptions"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

	row := db.QueryRow(sqlStatement)
	s := New()

	if err := row.Scan(&s.ID, &s.ProjectID, &s.Email, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return nil, fmt.Errorf("unable to scan a subscription row: %v", err)
	}

	return s, nil
}

// getSubscriptions returns all subscriptions in the database
func getSubscriptionsWhere(expression string) ([]*Subscription, error) {
	var subscriptions []*Subscription

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := "SELECT subscription_id,project_id,email,created_at,updated_at FROM subscriptions"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return subscriptions, fmt.Errorf("unable to execute `%s`: %v", sqlStatement, err)
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

// updateSubscription updates a subscription in the database
func updateSubscription(s *Subscription) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	// only allows updates to the email field, the other fields are immutable
	sqlStatement := `UPDATE subscriptions SET email=$1 WHERE subscription_id=$2`

	res, err := db.Exec(sqlStatement, s.Email, s.ID)
	if err != nil {
		return fmt.Errorf(
			"unable to execute `%s` with subscription_id `%v`: %v",
			sqlStatement, s.ID, err,
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info("subscription rows updated", zap.Int64("total", rowsAffected), zap.Int64("subscription_id", s.ID))

	return nil
}

// deleteSubscription deletes a subscription from database
func deleteSubscription(subscriptionID, projectID int64) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `DELETE FROM subscriptions WHERE subscription_id=$1 AND project_id=$2`

	res, err := db.Exec(sqlStatement, subscriptionID, projectID)
	if err != nil {
		return fmt.Errorf(
			"unable to execute `%s` with subscription_id `%d` and project_id `%d`: %v",
			sqlStatement, subscriptionID, projectID, err,
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info(
		"subscription rows deleted",
		zap.Int64("total", rowsAffected),
		zap.Int64("subscription_id", subscriptionID),
		zap.Int64("project_id", projectID),
	)

	return nil
}
