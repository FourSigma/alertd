DROP SCHEMA IF EXISTS alertd;
CREATE SCHEMA alertd;
CREATE TABLE alertd.user
(
    id uuid PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL UNIQUE,
    password_salt text NOT NULL UNIQUE,
    password_hash text NOT NULL,
    state_id text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
CREATE TABLE alertd.user_token
(
    user_id uuid REFERENCES alertd.user(id) PRIMARY KEY,
    token text NOT NULL UNIQUE,
    state_id text NOT NULL
);
CREATE TABLE alertd.topic
(
    id uuid PRIMARY KEY,
    user_id uuid REFERENCES alertd.user(id) NOT NULL,
    name text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
CREATE TABLE alertd.topic_message
(
    topic_id uuid REFERENCES alertd.topic PRIMARY KEY,
    type_id text NOT NULL,
    msg text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
INSERT INTO alertd.user
    (id, first_name, last_name, email, password_salt, password_hash, state_id, created_at, updated_at)
VALUES
    ( uuid_generate_v4(), 'Siva', 'Manivannan', 'siva@alertd.com', 'hello', 'hello', 'Active', now(), now());
