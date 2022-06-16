CREATE TABLE IF NOT EXISTS projects(
  project_id SERIAL PRIMARY KEY,
  domain VARCHAR (300) UNIQUE NOT NULL
);

ALTER TABLE subscriptions
  ADD COLUMN project_id INTEGER
  REFERENCES projects (project_id)
  ON DELETE CASCADE;
