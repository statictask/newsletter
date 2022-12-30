package post

// Posts is the entity used for lazy controlling
// interactions with many Posts in the database
type ProjectPosts struct {
	projectID int64
}

// NewProjectPosts returns a Posts controller
func NewProjectPosts(projectID int64) *ProjectPosts {
	return &ProjectPosts{projectID}
}

// All returns all project's posts
func (pp *ProjectPosts) All() ([]*Post, error) {
	return getPostsByProjectID(pp.projectID)
}

// Last returns the last project's post
func (pp *ProjectPosts) Last() (*Post, error) {
	return getLastPostByProjectID(pp.projectID)
}
