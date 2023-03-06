CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tasks (
	id bigserial NOT NULL,
	"title" TEXT,
    "repo_id" INT NOT NULL,
    "description" TEXT,
    "text" TEXT,
    "create_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    "update_at" timestamptz default NULL,
    CONSTRAINT tasks_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS repos (
	id bigserial NOT NULL,
	"title" TEXT default NULL,
    "description" TEXT default NULL,
    "user_id" uuid NOT NULL,
    "create_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    "update_at" timestamptz default NULL,
    CONSTRAINT repos_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4(),
	"email" varchar(100),
    "password" varchar(100),
    "username" TEXT,
    "firebase_token" TEXT,
    "create_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    "update_at" timestamptz default NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS api_keys (
    id uuid DEFAULT uuid_generate_v4(),
	"apikey" text NOT NULL,
	"scopes" text NOT NULL,
	"user_id" uuid NOT NULL,
    "name" text NOT NULL,
    "name_normalized" text NOT NULL,
    "expire_at" timestamptz default NULL,
	"create_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	"update_at" timestamptz default NULL,
    CONSTRAINT api_keys_pkey PRIMARY KEY (id)
);

-- MY UNDERGROUND ACTION TO PERFMORM WHAT I NEED TO EASILY IMPLEMENTATION
-- START TABLEs WITH WILL BE USED AS A PROXY BEHIND MY BUSINESS RULE WITH GOLANG OAUTH2 PACKAGE
CREATE TABLE IF not  exists  users_oauth_client(
	id bigserial NOT NULL,
    "user_id" uuid NOT NULL,
	"oauth_client_id" uuid NOT NULL,
    CONSTRAINT users_oauth_client_pkey PRIMARY KEY (id)

);

CREATE TABLE IF not exists  oauth_client_application(
	id bigserial NOT NULL,
    "appname" text NOT NULL,
    "mode" INT NOT NULL,
	"oauth_client_id" uuid NOT NULL,
    "create_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	"update_at" timestamptz default NULL,
    CONSTRAINT oauth_client_application_pkey PRIMARY KEY (id)
);
-- END TABLES

-- DONT CONTROLLED TABLE THOSE TABLES IS PROVIDED BY GO OAUTH PACKAGE. I THINK I DO THE WRONG CHOOSE. BUT ONLY HAVE MORE ONE WEEK TO CLOSE ALL MY MIND THREADS.
-- START TABLEs THAT WILL BE CREATED BY OAUTH2 PACKAGE.

CREATE TABLE IF not exists oauth2_clients (
	id text NOT NULL,
	"secret" text NOT NULL,
	"domain" text NOT NULL,
	"data" jsonb NOT NULL,
	CONSTRAINT oauth2_clients_pkey PRIMARY KEY (id)
);

CREATE TABLE IF not exists oauth2_tokens (
	id bigserial NOT NULL,
	"created_at" timestamptz NOT NULL,
	"expires_at" timestamptz NOT NULL,
	"code" text NOT NULL,
	"access" text NOT NULL,
	"refresh" text NOT NULL,
	"data" jsonb NOT NULL,
	CONSTRAINT oauth2_tokens_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_oauth2_tokens_access ON oauth2_tokens USING btree (access);
CREATE INDEX idx_oauth2_tokens_code ON oauth2_tokens USING btree (code);
CREATE INDEX idx_oauth2_tokens_expires_at ON oauth2_tokens USING btree (expires_at);
CREATE INDEX idx_oauth2_tokens_refresh ON oauth2_tokens USING btree (refresh);

-- END TABLE NOT UNDER MY CONTROL


ALTER TABLE users_oauth_client
ADD FOREIGN KEY ("user_id") REFERENCES users (id)
ON DELETE CASCADE
DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE api_keys
ADD FOREIGN KEY ("user_id") REFERENCES users (id)
ON DELETE CASCADE
DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE repos
ADD FOREIGN KEY ("user_id") REFERENCES users (id)
ON DELETE CASCADE
DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE tasks
ADD FOREIGN KEY ("repo_id") REFERENCES repos (id)
ON DELETE CASCADE
DEFERRABLE INITIALLY DEFERRED