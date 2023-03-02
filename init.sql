CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tasks (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	title TEXT,
    repo_id INT NOT NULL,
    description TEXT,
    text TEXT,
    create_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    update_at timestamptz default NULL
);

CREATE TABLE IF NOT EXISTS repos (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	title TEXT default NULL,
    description TEXT default NULL,
    user_id uuid NOT NULL,
    create_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    update_at timestamptz default NULL
);

CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4(),
	email varchar(100),
    password varchar(100),
    username TEXT,
    firebaseToken TEXT,
    create_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    update_at timestamptz default NULL
);

CREATE TABLE IF NOT EXISTS api_keys (
    id uuid DEFAULT uuid_generate_v4(),
	apikey text NOT NULL,
	scopes text NOT NULL,
	user_id uuid NOT NULL,
    name text NOT NULL,
    name_normalized text NOT NULL,
    expire_at timestamptz default NULL,
	create_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	update_at timestamptz default NULL
);

ALTER TABLE api_keys
ADD FOREIGN KEY (user_id) REFERENCES users (id)
on delete cascade on update cascade
DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE repos
ADD FOREIGN KEY (user_id) REFERENCES users (id)
on delete cascade on update cascade
DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE tasks
ADD FOREIGN KEY (repo_id) REFERENCES repos (id)
on delete cascade on update cascade
DEFERRABLE INITIALLY DEFERRED;