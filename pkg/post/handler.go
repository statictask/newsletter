package post

import (
	"database/sql"
	"fmt"

	"github.com/statictask/newsletter/internal/database"
	"github.com/statictask/newsletter/internal/log"
	"go.uber.org/zap"
)

// insertPost inserts a Post in the database
func insertPost(p *Post) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `INSERT INTO posts (pipeline_id,content) VALUES ($1,$2) RETURNING post_id,created_at,updated_at`

	if err := db.QueryRow(
		sqlStatement,
		p.PipelineID,
		p.Content,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		log.L.Fatal("unable to execute the query", zap.Error(err))
	}

	log.L.Info(
		"successfully created post record",
		zap.Int64("pipeline_id", p.PipelineID),
		zap.Int64("post_id", p.ID),
	)

	return nil
}

// getPostsByProjectID returns all posts in the database based
// on a given expression
func getPostsByProjectID(projectID int64) ([]*Post, error) {
	var ps []*Post

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := fmt.Sprintf(
		`SELECT p.post_id, p.pipeline_id, p.content, p.created_at, p.updated_at FROM posts AS p
		JOIN pipelines AS pl ON p.pipeline_id = pl.pipeline_id WHERE pl.project_id = %d`,
		projectID,
	)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return ps, fmt.Errorf("unable to execute `%s`: %v", sqlStatement, err)
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

// getPostByPipelineID return a single row that matches a given expression
func getPostByPipelineID(pipelineID int64) (*Post, error) {
	query := fmt.Sprintf(
		"SELECT post_id,pipeline_id,content,created_at,updated_at FROM posts WHERE pipeline_id=%d",
		pipelineID,
	)

	return scanPost(query)
}

// getLastPostByProjectID returns all posts in the database based
// on a given expression
func getLastPostByProjectID(projectID int64) (*Post, error) {
	query := fmt.Sprintf(
		`SELECT p.post_id, p.pipeline_id, p.content, p.created_at, p.updated_at FROM posts AS p
		JOIN pipelines AS pl ON p.pipeline_id = pl.pipeline_id WHERE pl.project_id = %d
		ORDER BY p.created_at DESC LIMIT 1`,
		projectID,
	)

	return scanPost(query)
}

// scanPost returns a single post based on the given query
func scanPost(query string) (*Post, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	row := db.QueryRow(query)
	p := &Post{}

	if err := row.Scan(&p.ID, &p.PipelineID, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("unable to scan post row: %v", err)
		}

		return nil, nil
	}

	return p, nil
}
