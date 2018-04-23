DROP SCHEMA IF EXISTS alertd;
CREATE SCHEMA alertd;
CREATE TYPE user_state_id AS ENUM
('Active','EmailVerificationSent', 'PasswordResetRequested', 'Suspended', 'Flagged');
CREATE TABLE alertd.users
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
CREATE TABLE alertd.tokens
(
    user_id uuid REFERENCES alertd.users(id) ON DELETE CASCADE NOT NULL,
    token text NOT NULL UNIQUE,
    state_id text NOT NULL,
    PRIMARY KEY(user_id, token)
);
CREATE TABLE alertd.topics
(
    id uuid PRIMARY KEY,
    user_id uuid REFERENCES alertd.users(id) NOT NULL,
    name text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    UNIQUE(user_id, name)
);
CREATE TABLE alertd.messages
(
    id uuid,
    topic_id uuid REFERENCES alertd.topics(id) ON DELETE CASCADE,
    msg text NOT NULL,
    type_id text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    PRIMARY KEY(id)
);
INSERT INTO alertd.users
    (id, first_name, last_name, email, password_salt, password_hash, state_id, created_at, updated_at)
VALUES
    ( uuid_generate_v4(), 'TestFirstName', 'TestLastName', 'test@alertd.com', 'hello', 'hello', 'Active', now(), now());
