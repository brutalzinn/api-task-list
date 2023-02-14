CREATE TABLE IF NOT EXISTS tasks (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	title TEXT,
    description TEXT,
    user_id INT NOT NULL,
    create_at timestamp default NULL,
    update_at timestamp default NULL
);

CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	email varchar(100),
    password varchar(100),
    username TEXT,
    firebaseToken TEXT,
    create_at timestamp default NULL,
    update_at timestamp default NULL
);

CREATE TABLE IF NOT EXISTS api_keys (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	apikey text NOT NULL,
	scopes text NOT NULL,
	user_id INT NOT NULL,
	create_at timestamp default NULL,
	update_at timestamp default NULL
);

ALTER TABLE api_keys
ADD FOREIGN KEY (user_id) REFERENCES users (id)
on delete cascade on update cascade
DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE tasks
ADD FOREIGN KEY (user_id) REFERENCES users (id)
on delete cascade on update cascade
DEFERRABLE INITIALLY DEFERRED;