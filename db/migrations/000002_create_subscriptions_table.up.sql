BEGIN;

CREATE TABLE IF NOT EXISTS subscriptions(
	subscription_id serial PRIMARY KEY,
	email VARCHAR (300) UNIQUE NOT NULL,
	project_id INTEGER REFERENCES projects (project_id) ON DELETE CASCADE NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

SELECT db_manage_updated_at('subscriptions');

COMMIT;
