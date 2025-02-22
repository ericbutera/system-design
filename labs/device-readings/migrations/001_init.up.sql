-- Create the Device table
CREATE TABLE devices (
    id BIGINT PRIMARY KEY,
    device_id VARCHAR(255) UNIQUE NOT NULL,
    type VARCHAR(50),
    location VARCHAR(50)
);

-- Create the Reading table
CREATE TABLE readings (
    timestamp TIMESTAMPTZ NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    reading_type VARCHAR(50) NOT NULL,
    value FLOAT NOT NULL,
    CONSTRAINT readings_unique UNIQUE (timestamp, device_id, reading_type)
);

-- Convert the readings table into a hypertable for timescale
SELECT create_hypertable('readings', 'timestamp');
