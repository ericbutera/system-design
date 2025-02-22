INSERT INTO hotels (id, name, location) VALUES (1, 'Hotel A', 'Location A');
INSERT INTO hotels (id, name, location) VALUES (2, 'Hotel B', 'Location B');

INSERT INTO room_types (id, name) VALUES (1, 'Single');
INSERT INTO room_types (id, name) VALUES (2, 'Double');

INSERT INTO rooms (hotel_id, name, number, room_type_id) VALUES (1, 'Room A', '101',1);
INSERT INTO rooms (hotel_id, name, number, room_type_id) VALUES (1, 'Room B', '102',1);
INSERT INTO rooms (hotel_id, name, number, room_type_id) VALUES (1, 'Room C', '103',1);
INSERT INTO rooms (hotel_id, name, number, room_type_id) VALUES (2, 'Room A', '101',1);
INSERT INTO rooms (hotel_id, name, number, room_type_id) VALUES (2, 'Room B', '102',2);
INSERT INTO rooms (hotel_id, name, number, room_type_id) VALUES (2, 'Room C', '103',2);

INSERT INTO guests (id, name, email) VALUES (1, 'Example Person', 'example@example.com');
INSERT INTO payments (id, correlation_id, transaction_id, amount) VALUES (1, '1', '123456', 100.00);
INSERT INTO reservations (room_type_id, quantity, checkin, checkout, status, guest_id, hotel_id, payment_id) VALUES (1, 1, '2021-01-01', '2021-01-02', 'CONFIRMED', 1, 1, 1);
