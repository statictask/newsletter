package postitem

import (
	"database/sql"
	"fmt"

	"github.com/statictask/newsletter/internal/database"
)

// insertPostItem inserts a PostItem in the database
func insertPostItem(p *PostItem) error {
	query := `
		INSERT INTO post_items (
		  post_id,
		  title,
		  link,
		  content
	        )
		VALUES (
		  $1,
		  $2,
		  $3,
		  $4
	        )
		RETURNING
		  post_item_id,
		  post_id,
		  title,
		  link,
		  content,
		  created_at,
		  updated_at
	`

	savedPostItem, err := scanPostItem(query, p.PostID, p.Title, p.Link, p.Content)
	if err != nil {
		return err
	}

	*p = *savedPostItem

	return nil
}

// getPostItemsByPostID return a single row that matches a given expression
func getPostItemsByPostID(postID int64) ([]*PostItem, error) {
	query := `
		SELECT
		  post_item_id,
		  post_id,
		  title,
		  link,
		  content,
		  created_at,
		  updated_at
		FROM
		  post_items
		WHERE
		  post_id = $1
	`

	return scanPostItems(query, postID)
}

// scanPostItem returns a single post based on the given query
func scanPostItem(query string, params ...interface{}) (*PostItem, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	row := db.QueryRow(query, params...)
	p := &PostItem{}

	if err := row.Scan(&p.ID, &p.PostID, &p.Title, &p.Link, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("unable to scan post_item row: %v", err)
		}

		return nil, nil
	}

	return p, nil
}

// scanPostItems returns multiple posts based on the given query
func scanPostItems(query string, params ...interface{}) ([]*PostItem, error) {
	var ps []*PostItem

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(query, params...)
	if err != nil {
		return ps, fmt.Errorf("unable to execute `%s`: %v", query, err)
	}

	defer rows.Close()

	for rows.Next() {
		p := New()

		if err := rows.Scan(&p.ID, &p.PostID, &p.Title, &p.Link, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return ps, fmt.Errorf("unable to scan post_item row: %v", err)
		}

		ps = append(ps, p)
	}

	return ps, nil
}
