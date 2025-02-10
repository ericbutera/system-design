INSERT INTO hotels (name, location) VALUES ('Hotel A', 'Location A');
INSERT INTO hotels (name, location) VALUES ('Hotel B', 'Location B');

INSERT INTO room_types (id, name) VALUES (1, 'Single');
INSERT INTO room_types (id, name) VALUES (2, 'Double');

INSERT INTO rooms (hotel_id, name, number, room_type_id) VALUES (1, 'Room A', '101',1);
INSERT INTO rooms (hotel_id, name, number, room_type_id) VALUES (1, 'Room B', '102',1);
INSERT INTO rooms (hotel_id, name, number, room_type_id) VALUES (2, 'Room A', '101',1);
INSERT INTO rooms (hotel_id, name, number, room_type_id) VALUES (2, 'Room B', '102',2);

INSERT INTO room_type_inventory (hotel_id, room_type_id, date, total_inventory, total_reserved) VALUES (1, 1, NOW(), 2, 0);
INSERT INTO room_type_inventory (hotel_id, room_type_id, date, total_inventory, total_reserved) VALUES (2, 1, NOW(), 1, 0);
INSERT INTO room_type_inventory (hotel_id, room_type_id, date, total_inventory, total_reserved) VALUES (2, 2, NOW(), 1, 0);
