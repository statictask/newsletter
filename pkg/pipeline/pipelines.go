package pipeline

import "fmt"

// Pipelines is the entity used for lazy controlling
// interactions with many Pipelines in the database
type Pipelines struct {
	projectID int64
}

// NewPipelines returns a Pipelines controller
func NewPipelines(projectID int64) *Pipelines {
	return &Pipelines{projectID}
}

// All returns all the ps registered in the database
// for a given project
func (ps *Pipelines) All() ([]*Pipeline, error) {
	exp := fmt.Sprintf("project_id=%d", ps.projectID)

	pArray, err := getPipelinesWhere(exp)
	if err != nil {
		return []*Pipeline{}, fmt.Errorf("unable to get pipelines: %v", err)
	}

	return pArray, nil
}

// Get returns a single Pipeline according to its ID
func (ps *Pipelines) Get(pipelineID int64) (*Pipeline, error) {
	exp := fmt.Sprintf("pipeline_id=%d AND project_id=%d", pipelineID, ps.projectID)

	return getPipelineWhere(exp)
}

// Where return many Pipelines according to an expression
func (ps *Pipelines) Where(exp string) ([]*Pipeline, error) {
	exp = fmt.Sprintf("%s AND project_id=%d", exp, ps.projectID)

	pArray, err := getPipelinesWhere(exp)
	if err != nil {
		return nil, fmt.Errorf("unable to get pipelines: %v", err)
	}

	return pArray, nil
}

// Delete deletes a Pipeline based on its ID
func (ps *Pipelines) Delete(pipelineID int64) error {
	if err := deletePipeline(pipelineID); err != nil {
		return fmt.Errorf("unable to delete pipeline: %v", err)
	}

	return nil
}

// Add creates a new entry in the project's ps
// the function creates a new pipeline entry in the database
func (ps *Pipelines) Create() (*Pipeline, error) {
	// make sure the Pipeline has the correct ProjectID
	p := New()
	p.ProjectID = ps.projectID

	if err := p.Create(); err != nil {
		return nil, err
	}

	return p, nil
}

// Running returns the running pipeline if it exists
func (ps *Pipelines) Last() (*Pipeline, error) {
	return getLastPipeline(ps.projectID)
}
