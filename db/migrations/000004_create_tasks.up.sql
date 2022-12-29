BEGIN;
	
DO $$ BEGIN
	CREATE TYPE task_type_t AS ENUM ('Scrape', 'Publish');
EXCEPTION
    	WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
	CREATE TYPE task_status_t AS ENUM ('Waiting', 'Ready', 'Running', 'Finished', 'Failed', 'Unknown', 'Aborted');
EXCEPTION
    	WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS tasks (
	task_id SERIAL PRIMARY KEY,  
	pipeline_id INTEGER REFERENCES pipelines (pipeline_id) ON DELETE CASCADE NOT NULL,
	task_type task_type_t NOT NULL,
	task_status task_status_t NOT NULL DEFAULT 'Waiting',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

SELECT db_manage_updated_at('tasks');

COMMIT;
