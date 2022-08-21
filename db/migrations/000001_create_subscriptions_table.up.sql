BEGIN;

CREATE TABLE IF NOT EXISTS subscriptions(
   subscription_id serial PRIMARY KEY,
   email VARCHAR (300) UNIQUE NOT NULL
);

COMMIT;
