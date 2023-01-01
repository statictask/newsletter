package post

import (
	"database/sql"
	"fmt"

	"github.com/statictask/newsletter/internal/database"
)

// insertPost inserts a Post in the database
func insertPost(p *Post) error {
	query := `
		INSERT INTO posts (
		  pipeline_id,
		  content
	        )
		VALUES (
		  $1,
		  $2
	        )
		RETURNING
		  post_id,
		  pipeline_id,
		  content,
		  created_at,
		  updated_at
	`

	savedPost, err := scanPost(query, p.PipelineID, p.Content)
	if err != nil {
		return err
	}

	*p = *savedPost

	return nil
}

// getPostsByProjectID returns all posts in the database based
// on a given expression
func getPostsByProjectID(projectID int64) ([]*Post, error) {
	query := `
		SELECT
		  p.post_id,
		  p.pipeline_id,
		  p.content,
		  p.created_at,
		  p.updated_at
		FROM
		  posts AS p
		JOIN pipelines AS pl
		  ON p.pipeline_id = pl.pipeline_id
		WHERE
		  pl.project_id = $1
	`

	return scanPosts(query, projectID)
}

// getPostByPipelineID return a single row that matches a given expression
func getPostByPipelineID(pipelineID int64) (*Post, error) {
	query := `
		SELECT
		  post_id,
		  pipeline_id,
		  content,
		  created_at,
		  updated_at
		FROM
		  posts
		WHERE
		  pipeline_id = $1
	`

	return scanPost(query, pipelineID)
}

// getLastPostByProjectID returns all posts in the database based
// on a given expression
func getLastPostByProjectID(projectID int64) (*Post, error) {
	query := `
		SELECT
		  p.post_id,
		  p.pipeline_id,
		  p.content,
		  p.created_at,
		  p.updated_at
		FROM
		  posts AS p
		JOIN pipelines AS pl
		  ON p.pipeline_id = pl.pipeline_id
		WHERE
		  pl.project_id = $1
		ORDER BY
		  p.created_at
		DESC
		LIMIT 1
	`

	return scanPost(query, projectID)
}

// scanPost returns a single post based on the given query
func scanPost(query string, params ...interface{}) (*Post, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	row := db.QueryRow(query, params...)
	p := &Post{}

	if err := row.Scan(&p.ID, &p.PipelineID, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("unable to scan post row: %v", err)
		}

		return nil, nil
	}

	return p, nil
}

// scanPosts returns multiple posts based on the given query
func scanPosts(query string, params ...interface{}) ([]*Post, error) {
	var ps []*Post

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

		if err := rows.Scan(&p.ID, &p.PipelineID, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return ps, fmt.Errorf("unable to scan post row: %v", err)
		}

		ps = append(ps, p)
	}

	return ps, nil
}
