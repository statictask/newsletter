package subscription

import (
	"fmt"
	"strconv"

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

	sqlStatement := `INSERT INTO subscriptions (email) VALUES ($1) RETURNING subscription_id`

	if err := db.QueryRow(
		sqlStatement,
		subscription.Email,
	).Scan(&subscription.ID); err != nil {
		log.L.Fatal("unable to execute the query", zap.Error(err))
	}

	log.L.Info(
		"created subscription record",
		zap.Int64("subscriptionID", subscription.ID),
		zap.String("email", subscription.Email),
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

	sqlStatement := "SELECT subscription_id,email FROM subscriptions"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

	row := db.QueryRow(sqlStatement)
	subscription := New()

	if err := row.Scan(&subscription.ID, &subscription.Email); err != nil {
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

	sqlStatement := "SELECT subscription_id,email FROM subscriptions"

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

		if err := rows.Scan(&subscription.ID, &subscription.Email); err != nil {
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

	sqlStatement := `UPDATE subscriptions SET email=$1 WHERE subscription_id=$2`

	res, err := db.Exec(sqlStatement, subscription.Email, subscription.ID)
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

// loadSubscription is a helper function that receives an string with the
// subscription_id and returns a subscription instance loaded from db
func loadSubscription(subscriptionID string) (*Subscription, error) {
	id, err := castID(subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed casting id: %v", err)
	}

	subscriptions := NewSubscriptions()

	subscription, err := subscriptions.Get(int64(id))
	if err != nil {
		return nil, fmt.Errorf("failed retrieving subscription: %v", err)
	}

	return subscription, nil
}

// castID converts a string ID to an int64 ID
func castID(strID string) (int64, error) {
	id, err := strconv.Atoi(strID)
	if err != nil {
		return -1, fmt.Errorf("unable to parse subscription_id into int: %v", err)
	}

	return int64(id), nil
}
