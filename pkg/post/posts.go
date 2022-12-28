package post

import "fmt"

// Posts is the entity used for lazy controlling
// interactions with many Posts in the database
type Posts struct {
	pipelineID int64
}

// NewPosts returns a Posts controller
func NewPosts(pipelineID int64) *Posts {
	return &Posts{pipelineID}
}

// All returns all the ps registered in the database
// for a given project
func (ps *Posts) All() ([]*Post, error) {
	exp := fmt.Sprintf("pipeline_id=%d", ps.pipelineID)

	pArray, err := getPostsWhere(exp)
	if err != nil {
		return []*Post{}, fmt.Errorf("unable to get posts: %v", err)
	}

	return pArray, nil
}

// Get returns a single Post according to its ID
func (ps *Posts) Get(postID int64) (*Post, error) {
	exp := fmt.Sprintf("post_id=%d AND pipeline_id=%d", postID, ps.pipelineID)

	return getPostWhere(exp)
}

// Where return many Posts according to an expression
func (ps *Posts) Where(exp string) ([]*Post, error) {
	exp = fmt.Sprintf("%s AND pipeline_id=%d", exp, ps.pipelineID)

	tArray, err := getPostsWhere(exp)
	if err != nil {
		return nil, fmt.Errorf("unable to get posts: %v", err)
	}

	return tArray, nil
}

// Delete deletes a Post based on its ID
func (ps *Posts) Delete(postID int64) error {
	if err := deletePost(postID); err != nil {
		return fmt.Errorf("unable to delete post: %v", err)
	}

	return nil
}

// Add creates a new entry in the pipeline's posts
func (ps *Posts) Add(p *Post) error {
	// make sure the Post has the correct ProjectID
	p.PipelineID = ps.pipelineID

	if err := p.Create(); err != nil {
		return err
	}

	return nil
}
