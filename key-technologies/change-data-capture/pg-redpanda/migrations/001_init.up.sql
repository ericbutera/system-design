CREATE TABLE
    users (id BIGINT PRIMARY KEY);

CREATE TABLE
    urls (
        slug VARCHAR(11) NOT NULL PRIMARY KEY,
        long VARCHAR(2048) NOT NULL,
        user_id BIGINT REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMP
    );