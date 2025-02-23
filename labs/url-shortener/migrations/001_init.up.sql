CREATE TABLE
    users (id BIGINT PRIMARY KEY);

CREATE TABLE
    urls_v0 (
        slug VARCHAR(11) NOT NULL PRIMARY KEY,
        long VARCHAR(2048) NOT NULL,
        user_id BIGINT REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMP
    );

CREATE TABLE
    urls_v1 (
        slug VARCHAR(11) NOT NULL PRIMARY KEY,
        long VARCHAR(2048) NOT NULL,
        user_id BIGINT REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMP
    );

CREATE TABLE
    urls_v2 (
        slug VARCHAR(11) NOT NULL PRIMARY KEY,
        long VARCHAR(2048) NOT NULL,
        user_id BIGINT REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMP
    );

CREATE TABLE
    url_stats (
        slug VARCHAR(11) REFERENCES urls_v1 (slug),
        view_count BIGINT
    );

CREATE SEQUENCE url_counter START 1;