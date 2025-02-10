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
        hotel_id INT NOT NULL REFERENCES hotels (id) ON DELETE CASCADE,
        room_type_id INT NOT NULL REFERENCES room_types (id) ON DELETE CASCADE,
        name VARCHAR(255) NOT NULL,
        number VARCHAR(50) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    payments (
        id SERIAL PRIMARY KEY,
        correlation_id INT NOT NULL, -- reservation id
        transaction_id VARCHAR(255) NOT NULL, -- transaction id from the payment gateway
        amount DECIMAL(10, 2) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        -- TODO: currency column (you can add a currency column if needed)
    );

CREATE TABLE
    reservations (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        room_type_id INT NOT NULL REFERENCES room_types (id) ON DELETE RESTRICT,
        quantity INT NOT NULL,
        email VARCHAR(255) NOT NULL,
        date DATE NOT NULL,
        time TIME NOT NULL,
        payment_id INT REFERENCES payments (id) ON DELETE SET NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    room_type_inventory (
        hotel_id INT NOT NULL REFERENCES hotels (id) ON DELETE CASCADE,
        room_type_id INT NOT NULL REFERENCES room_types (id) ON DELETE CASCADE,
        date DATE,
        total_inventory INT,
        total_reserved INT,
        PRIMARY KEY (hotel_id, room_type_id, date)
    );
