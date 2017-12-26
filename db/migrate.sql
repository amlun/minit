CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64),
    UNIQUE (name)
);

CREATE TABLE relationships (
    id SERIAL PRIMARY KEY,
    owner_id INT,
    user_id INT,
    state VARCHAR(16),
    re_state VARCHAR (16)
    UNIQUE (owner_id, user_id)
);