BEGIN;

CREATE TABLE IF NOT EXISTS feed_readers(
	feed_reader_id serial PRIMARY KEY,
	content_hash VARCHAR (300) NOT NULL,
	project_id INTEGER REFERENCES projects (project_id) ON DELETE CASCADE NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

SELECT db_manage_updated_at('feed_readers');

COMMIT;
