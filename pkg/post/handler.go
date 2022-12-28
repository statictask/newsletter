package post

import (
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

	sqlStatement := `INSERT INTO posts (pipeline_id,content) VALUES ($1) RETURNING post_id,created_at,updated_at`

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

// getPost return a single row that matches a given expression
func getPostWhere(expression string) (*Post, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := "SELECT post_id,pipeline_id,content,created_at,updated_at FROM posts"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

	row := db.QueryRow(sqlStatement)
	p := &Post{}

	if err := row.Scan(&p.ID, &p.PipelineID, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, fmt.Errorf("unable to scan post row: %v", err)
	}

	return p, nil
}

// getPosts returns all posts in the database based
// on a given expression
func getPostsWhere(expression string) ([]*Post, error) {
	var ps []*Post

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	sqlStatement := "SELECT post_id,pipeline_id,content,created_at,updated_at FROM posts"

	if expression != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, expression)
	}

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

// updatePost updates a Post in the database
func updatePost(p *Post) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	// only allows update to the post_status field, the other fields are immutable
	sqlStatement := `UPDATE posts SET content WHERE post_id=$3`

	res, err := db.Exec(sqlStatement, p.Content, p.ID)
	if err != nil {
		return fmt.Errorf(
			"unable to execute `%s` with post_id `%d`: %v",
			sqlStatement, p.ID, err,
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info(
		"posts rows updated",
		zap.Int64("total", rowsAffected),
		zap.Int64("post_id", p.ID),
	)

	return nil
}

// deletePost deletes a Post record from database
func deletePost(postID int64) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	sqlStatement := `DELETE FROM posts WHERE post_id=$1`

	res, err := db.Exec(sqlStatement, postID)
	if err != nil {
		return fmt.Errorf(
			"unable to execute `%s` with post_id `%d`: %v",
			sqlStatement, postID, err,
		)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed checking the affected rows: %v", err)
	}

	log.L.Info(
		"post rows deleted",
		zap.Int64("total", rowsAffected),
		zap.Int64("post_id", postID),
	)

	return nil
}
