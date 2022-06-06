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
func getUserWhere(expression string) (*Subscription, error) {
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

	log.L.Info("user rows updated", zap.Int64("total", rowsAffected))

	return nil
}

// deleteSubscription deletes a subscription from database
func deleteSubscription(subscription *Subscription) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `DELETE FROM subscriptions WHERE subscription_id=$1`

	res, err := db.Exec(sqlStatement, subscription.ID)
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

	log.L.Info("user rows deleted", zap.Int64("total", rowsAffected))

	return nil
}
