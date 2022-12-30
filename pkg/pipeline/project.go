package pipeline

// ProjectPipelines is the entity used for lazy controlling
// interactions with many ProjectPipelines in the database
type ProjectPipelines struct {
	projectID int64
}

// NewProjectPipelines returns a ProjectPipelines controller
func NewProjectPipelines(projectID int64) *ProjectPipelines {
	return &ProjectPipelines{projectID}
}

// All returns all the ps registered in the database
// for a given project
func (ps *ProjectPipelines) All() ([]*Pipeline, error) {
	return getPipelinesByProjectID(ps.projectID)
}

// Add creates a new entry in the project's ps
// the function creates a new pipeline entry in the database
func (ps *ProjectPipelines) Create() (*Pipeline, error) {
	// make sure the Pipeline has the correct ProjectID
	p := New()
	p.ProjectID = ps.projectID

	if err := p.Create(); err != nil {
		return nil, err
	}

	return p, nil
}

// Last returns the last pipeline of the respective Project
func (ps *ProjectPipelines) Last() (*Pipeline, error) {
	return getLastPipeline(ps.projectID)
}
