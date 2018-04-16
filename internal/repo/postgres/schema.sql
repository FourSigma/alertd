CREATE TABLE user(
       id             uuid PRIMARY KEY,
       first_name     text NOT NULL,
       last_name      text NOT NULL, 
       email          text NOT NULL UNIQUE,
       password_salt  text NOT NULL,
       password_hash  text NOT NULL,
       state_id        text NOT NULL,
       created_at     timestampz NOT NULL,
       updated_at     timestampz NOT NULL

);


CREATE TABLE user_token (
    user_id        uuid PRIMARY KEY,
    token          text NOT NULL UNIQUE,
    state_id        text NOT NULL,
);
