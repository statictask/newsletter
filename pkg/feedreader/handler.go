package feedreader

import (
	"fmt"

	"github.com/statictask/newsletter/internal/database"
	"github.com/statictask/newsletter/internal/log"
	"go.uber.org/zap"
)

// insertFeedReader inserts a FeedReader in the database
func insertFeedReader(fr *FeedReader) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `INSERT INTO feed_readers (project_id,content_hash,feed_url) VALUES ($1, $2, $3) RETURNING feed_reader_id,created_at,updated_at`

	if err := db.QueryRow(
		sqlStatement,
		fr.ProjectID,
		fr.ContentHash,
		fr.FeedURL,
	).Scan(&fr.ID, &fr.CreatedAt, &fr.UpdatedAt); err != nil {
		log.L.Fatal("unable to execute the query", zap.Error(err))
	}

	log.L.Info("successfully created feed reader record", zap.Int64("project_id", fr.ProjectID), zap.Int64("feed_reader_id", fr.ID))

	return nil
}

// getFeedReader return a single row that matches a given expression
func getFeedReaderWhere(expression string) (*FeedReader, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := "SELECT feed_reader_id,project_id,content_hash,feed_url,created_at,updated_at FROM feed_readers"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

	row := db.QueryRow(sqlStatement)
	fr := &FeedReader{}

	if err := row.Scan(&fr.ID, &fr.ProjectID, &fr.ContentHash, &fr.FeedURL, &fr.CreatedAt, &fr.UpdatedAt); err != nil {
		return nil, fmt.Errorf("unable to scan feed reader row: %v", err)
	}

	return subscription, nil
}

// getFeedReaders returns all feed readers in the database based
// on a given expression
func getFeedReadersWhere(expression string) ([]*FeedReader, error) {
	var frs []*FeedReader

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := "SELECT feed_reader_id,project_id,content_hash,feed_url,created_at,updated_at FROM fead_readers"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return frs, fmt.Errorf("unable to execute `%s`: %v", sqlStatement, err)
	}

	defer rows.Close()

	for rows.Next() {
		fr := New()

		if err := rows.Scan(&fr.ID, &fr.ProjectID, &fr.ContentHash, &fr.FeedURL, &fr.CreatedAt, &fr.UpdatedAt); err != nil {
			return frs, fmt.Errorf("unable to scan feed reader row: %v", err)
		}

		frs = append(frs, fr)
	}

	return frs, nil

}

// updateFeedReader updates a FeedReader in the database
func updateFeedReader(fr *FeedReader) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	// only allows update to the content_hash and feed_url fields, the other fields are immutable
	sqlStatement := `UPDATE feed_readers SET content_hash=$1,feed_url=$2 WHERE feed_reader_id=$3`

	res, err := db.Exec(sqlStatement, fr.ContentHash, fr.FeedURL, fr.ID)
	if err != nil {
		return fmt.Errorf(
			"unable to execute `%s`: %v",
			sqlStatement, fr.ID, err,
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info("feed reader rows updated", zap.Int64("total", rowsAffected), zap.Int64("feed_reader_id", fr.ID))

	return nil
}

// deleteFeedReader deletes a FeedReader record from database
func deleteFeedReader(feedReaderID, projectID int64) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `DELETE FROM feed_readers WHERE feed_reader_id=$1 AND project_id=$2`

	res, err := db.Exec(sqlStatement, feedReaderID, projectID)
	if err != nil {
		return fmt.Errorf(
			"unable to execute `%s` with feed_reader_id `%d` and project_id `%d`: %v",
			sqlStatement, feedReaderID, projectID, err,
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info(
		"feed reader rows deleted",
		zap.Int64("total", rowsAffected),
		zap.Int64("feed_reader_id", feedReaderID),
		zap.Int64("project_id", projectID),
	)

	return nil
}

// insertEvent inserts a new Event in the database
func insertEvent(e *Event) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `INSERT INTO feed_reader_events (feed_reader_id,content_hash) VALUES ($1, $2) RETURNING feed_reader_event_id,created_at`

	if err := db.QueryRow(
		sqlStatement,
		e.FeedReaderID,
		e.ContentHash,
	).Scan(&e.ID, &e.CreatedAt); err != nil {
		log.L.Fatal("unable to execute the query", zap.Error(err))
	}

	log.L.Info(
		"successfully created feed reader event record",
		zap.Int64("feed_reader_id", e.FeedReaderID),
		zap.Int64("feed_reader_event_id", e.ID),
	)

	return nil
}

// getLastEventWhere
func getLastEventWhere(expression string) (*Event, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := "SELECT feed_reader_event_id,feed_reader_id,content_hash,created_at FROM feed_reader_events"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

	sqlStatement = fmt.Sprintf("%s ORDER BY created_at DESC LIMIT 1", sqlStatement)

	row := db.QueryRow(sqlStatement)
	e := NewEvent()

	if err := row.Scan(&e.ID, &e.FeedReaderID, &e.ContentHash, &fr.CreatedAt); err != nil {
		return nil, fmt.Errorf("unable to scan feed reader row: %v", err)
	}

	return subscription, nil
}

// deleteEvent deletes an Event record from database
func deleteEvent(eventID, feedReaderID int64) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `DELETE FROM feed_reader_events WHERE feed_reader_event_id=$1 AND feed_reader_id=$2`

	res, err := db.Exec(sqlStatement, eventID, feedReaderID)
	if err != nil {
		return fmt.Errorf(
			"unable to execute `%s` with for feed_reader_event_id `%d` and feed_reader_id `%d`: %v",
			sqlStatement, eventID, feedReaderID, err,
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info(
		"feed reader event rows deleted",
		zap.Int64("total", rowsAffected),
		zap.Int64("feed_reader_event_id", eventID),
		zap.Int64("feed_reader_id", feedReaderID),
	)

	return nil
}
