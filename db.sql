CREATE TABLE users (
  uuid UUID PRIMARY KEY,
  username VARCHAR(20) NOT NULL,
  password VARCHAR(60) NOT NULL,
  created BIGINT NOT NULL,
  removed BIGINT NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX on users (username, removed);

CREATE TABLE tokens (
  uuid UUID PRIMARY KEY,
  user_uuid UUID NOT NULL,
  token VARCHAR(60) NOT NULL,
  created BIGINT NOT NULL,
  expires BIGINT NOT NULL,
  removed BIGINT NOT NULL DEFAULT 0
);

CREATE INDEX on tokens (user_uuid, token, removed);
