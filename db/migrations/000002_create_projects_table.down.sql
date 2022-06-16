DROP TABLE IF EXISTS projects;

ALTER TABLE subscriptions
  DROP COLUMN project_id;
