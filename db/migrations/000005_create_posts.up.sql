BEGIN;

CREATE TABLE IF NOT EXISTS posts (
	post_id SERIAL PRIMARY KEY,  
	pipeline_id INTEGER REFERENCES pipelines (pipeline_id) ON DELETE CASCADE NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

SELECT db_manage_updated_at('posts');

COMMIT;