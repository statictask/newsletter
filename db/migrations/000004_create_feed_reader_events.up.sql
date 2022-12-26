BEGIN;

CREATE TABLE IF NOT EXISTS feed_reader_events (
	feed_reader_event_id serial PRIMARY KEY,
	content_hash VARCHAR (300) NOT NULL,
	feed_reader_id INTEGER REFERENCES feed_readers (feed_reader_id) ON DELETE CASCADE NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
);

COMMIT;
