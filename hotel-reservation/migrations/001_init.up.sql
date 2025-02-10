CREATE TABLE
    hotels (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        location VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    room_types (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    rooms (
        id SERIAL PRIMARY KEY,
        hotel_id INT NOT NULL REFERENCES hotels (id) ON DELETE RESTRICT,
        room_type_id INT NOT NULL REFERENCES room_types (id) ON DELETE RESTRICT,
        name VARCHAR(255) NOT NULL,
        number VARCHAR(50) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    payments (
        id SERIAL PRIMARY KEY,
        correlation_id VARCHAR(128) NOT NULL, -- reservation id
        transaction_id VARCHAR(255) NOT NULL, -- transaction id from the payment gateway
        amount DECIMAL(10, 2) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT unique_payment UNIQUE (correlation_id)
    );

CREATE TABLE
    guests (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255),
        email VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    reservations (
        id SERIAL PRIMARY KEY,
        room_type_id INT NOT NULL REFERENCES room_types (id) ON DELETE RESTRICT,
        quantity INT NOT NULL,
        checkin DATE NOT NULL,
        checkout DATE NOT NULL,
        status VARCHAR(20) CHECK (status IN ('PENDING', 'CONFIRMED', 'CANCELLED')),
        guest_id INT NOT NULL REFERENCES guests (id) ON DELETE RESTRICT,
        hotel_id INT NOT NULL REFERENCES hotels (id) ON DELETE RESTRICT,
        payment_id INT REFERENCES payments (id) ON DELETE SET NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    room_type_inventory (
        hotel_id INT NOT NULL REFERENCES hotels (id) ON DELETE RESTRICT,
        room_type_id INT NOT NULL REFERENCES room_types (id) ON DELETE RESTRICT,
        date DATE,
        total_inventory INT,
        total_reserved INT,
        PRIMARY KEY (hotel_id, room_type_id, date)
    );
