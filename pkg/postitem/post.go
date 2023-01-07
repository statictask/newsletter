package postitem

// PostPostItems is the entity used for lazy controlling
// interactions with many PostItems in the database
type PostPostItems struct {
	postID int64
}

// NewPostPostItems returns a Posts controller
func NewPostPostItems(postID int64) *PostPostItems {
	return &PostPostItems{postID}
}

// All returns all project's posts
func (pp *PostPostItems) All() ([]*PostItem, error) {
	return getPostItemsByPostID(pp.postID)
}
