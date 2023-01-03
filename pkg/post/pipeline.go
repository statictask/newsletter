package post

// PipelinePosts is the entity used for lazy controlling
// interactions with many PipelinePosts in the database
type PipelinePosts struct {
	pipelineID int64
}

// NewPipelinePosts returns a PipelinePosts controller
func NewPipelinePosts(pipelineID int64) *PipelinePosts {
	return &PipelinePosts{pipelineID}
}

// Last returns the only post attached to the pipeline
func (ps *PipelinePosts) Last() (*Post, error) {
	return getPostByPipelineID(ps.pipelineID)
}
