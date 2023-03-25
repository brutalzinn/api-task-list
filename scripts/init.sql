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

-- API KEY INIT DATA
INSERT INTO api_keys (id, apikey, scopes, user_id, "name", name_normalized, expire_at, create_at, update_at) VALUES('574b32d5-f788-4b2c-aefb-979ef17a6602'::uuid, '$2a$04$qIHJRvfn2G7GO1Aq2d.zf.xHkTrB7q8tEISsANDKMuLK3nfdTyjc.', 'task_manager,repo_manager', 'ab7d7136-6c24-4cd0-ba30-97ff0110ecac'::uuid, 'Flutter app', 'flutterapp', '2023-04-23 22:32:07.641', '2023-03-24 22:32:07.656', NULL);

-- USER INIT DATA
INSERT INTO users (id, email, "password", username, firebase_token, create_at, update_at) VALUES('ab7d7136-6c24-4cd0-ba30-97ff0110ecac'::uuid, 'test', '$2a$15$zqJBJTKH7LZbTSQHhnNzeOx9VjcGwv3HamUksu8VQ81E/WbRJCLPW', 'usertest', '', '2023-03-17 09:47:02.881', NULL);

-- OAUTH INIT DATA
INSERT INTO oauth2_clients (id, secret, "domain", "data") VALUES('d45bc6bd-a066-482e-8b64-b56a34d9b2ba', '434977f5-2dc8-4290-bec5-e4a14e629d37', 'http://localhost:8888', '{"ID": "d45bc6bd-a066-482e-8b64-b56a34d9b2ba", "Domain": "http://localhost:8888", "Public": false, "Secret": "434977f5-2dc8-4290-bec5-e4a14e629d37", "UserID": "ab7d7136-6c24-4cd0-ba30-97ff0110ecac"}'::jsonb);
INSERT INTO oauth_client_application (id, appname, "mode", oauth_client_id, create_at, update_at) VALUES(22, 'Meu aplicativo CLI', 0, 'd45bc6bd-a066-482e-8b64-b56a34d9b2ba'::uuid, '2023-03-24 22:29:03.595', NULL);
INSERT INTO users_oauth_client (id, user_id, oauth_client_id) VALUES(22, 'ab7d7136-6c24-4cd0-ba30-97ff0110ecac'::uuid, 'd45bc6bd-a066-482e-8b64-b56a34d9b2ba'::uuid);
