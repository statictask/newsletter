package subscription

import (
	"fmt"

	"github.com/statictask/newsletter/internal/database"
	"github.com/statictask/newsletter/internal/log"
	"go.uber.org/zap"
)

// insertSubscription inserts a subscription in the database
func insertSubscription(subscription *Subscription) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `INSERT INTO subscriptions (project_id,email) VALUES ($1, $2) RETURNING subscription_id`

	if err := db.QueryRow(
		sqlStatement,
		subscription.ProjectID,
		subscription.Email,
	).Scan(&subscription.ID); err != nil {
		log.L.Fatal("unable to execute the query", zap.Error(err))
	}

	log.L.Info(
		"created subscription record",
		zap.Int64("ID", subscription.ID),
		zap.Int64("ProductID", subscription.ProjectID),
		zap.String("Email", subscription.Email),
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

	sqlStatement := "SELECT subscription_id,project_id,email FROM subscriptions"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

	row := db.QueryRow(sqlStatement)
	subscription := New()

	if err := row.Scan(&subscription.ID, &subscription.ProjectID, &subscription.Email); err != nil {
		return nil, fmt.Errorf("unable to scan a subscription row: %v", err)
	}

	return subscription, nil
}

// getSubscriptions returns all subscriptions in the database
func getSubscriptionsWhere(expression string) ([]*Subscription, error) {
	var subscriptions []*Subscription

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := "SELECT subscription_id,project_id,email FROM subscriptions"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return subscriptions, fmt.Errorf("unable to execute `%s`: %v", sqlStatement, err)
	}

	defer rows.Close()

	for rows.Next() {
		subscription := New()

		if err := rows.Scan(&subscription.ID, &subscription.ProjectID, &subscription.Email); err != nil {
			return subscriptions, fmt.Errorf("unable to scan a subscription row: %v", err)
		}

		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil

}

// updateSubscription updates a subscription in the database
func updateSubscription(subscription *Subscription) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `UPDATE subscriptions SET email=$1,project_id=$2 WHERE subscription_id=$3`

	res, err := db.Exec(sqlStatement, subscription.Email, subscription.ProjectID, subscription.ID)
	if err != nil {
		return fmt.Errorf(
			"unable to execute `%s` with subscription_id `%v`: %v",
			sqlStatement, subscription.ID, err,
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info("subscription rows updated", zap.Int64("total", rowsAffected))

	return nil
}

// deleteSubscription deletes a subscription from database
func deleteSubscription(subscriptionID int64) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `DELETE FROM subscriptions WHERE subscription_id=$1`

	res, err := db.Exec(sqlStatement, subscriptionID)
	if err != nil {
		return fmt.Errorf(
			"unable to execute `%s` with subscription_id `%v`: %v",
			sqlStatement, subscriptionID, err,
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info("subscription rows deleted", zap.Int64("total", rowsAffected))

	return nil
}

// getProjectSubscriptions returns all the subscriptions for a given project id
func getProjectSubscriptions(projectID int64) ([]*Subscription, error) {
	var subscriptions []*Subscription

	db, err := database.Connect()
	if err != nil {
		return subscriptions, err
	}

	defer db.Close()

	sqlStatement := `SELECT subscription_id,project_id,email FROM gadgets WHERE project_id=$1`

	rows, err := db.Query(sqlStatement, projectID)
	if err != nil {
		return subscriptions, fmt.Errorf(
			"unable to execute `%s`: %v",
			sqlStatement, err,
		)
	}

	defer rows.Close()

	for rows.Next() {
		s := New()

		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Email); err != nil {
			return subscriptions, fmt.Errorf("unable to scan a gadget row: %v", err)
		}

		subscriptions = append(subscriptions, s)
	}

	return subscriptions, err
}
