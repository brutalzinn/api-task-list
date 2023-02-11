CREATE TABLE IF NOT EXISTS tasks (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	title TEXT,
    description TEXT,
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