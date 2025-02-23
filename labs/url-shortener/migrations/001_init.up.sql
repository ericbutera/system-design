CREATE TABLE
    users (id BIGINT PRIMARY KEY);

CREATE TABLE
    urls (
        slug CHAR(11) NOT NULL PRIMARY KEY,
        long VARCHAR(2048) NOT NULL,
        user_id BIGINT REFERENCES users (id),
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMP
    );

CREATE TABLE
    counter (value BIGINT NOT NULL);

CREATE INDEX idx_slug ON urls (slug);

CREATE TABLE
    url_stats (
        slug CHAR(11) REFERENCES urls (slug),
        view_count BIGINT
    );

CREATE SEQUENCE url_counter START 1;