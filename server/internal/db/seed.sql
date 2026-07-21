-- Pet Types
INSERT INTO pet_types (id, name) VALUES (1, 'Cat');
INSERT INTO pet_types (id, name) VALUES (2, 'Dog');
INSERT INTO pet_types (id, name) VALUES (3, 'Lizard');
INSERT INTO pet_types (id, name) VALUES (4, 'Snake');
INSERT INTO pet_types (id, name) VALUES (5, 'Bird');
INSERT INTO pet_types (id, name) VALUES (6, 'Hamster');
INSERT INTO pet_types (id, name) VALUES (7, 'Camel');

-- Specialties
INSERT INTO specialties (id, name) VALUES (1, 'Radiology');
INSERT INTO specialties (id, name) VALUES (2, 'Surgery');
INSERT INTO specialties (id, name) VALUES (3, 'Dentistry');

-- Vets
INSERT INTO vets (id, first_name, last_name) VALUES (1, 'James', 'Carter');
INSERT INTO vets (id, first_name, last_name) VALUES (2, 'Helen', 'Leary');
INSERT INTO vets (id, first_name, last_name) VALUES (3, 'Linda', 'Douglas');
INSERT INTO vets (id, first_name, last_name) VALUES (4, 'Rafael', 'Ortega');
INSERT INTO vets (id, first_name, last_name) VALUES (5, 'Henry', 'Stevens');
INSERT INTO vets (id, first_name, last_name) VALUES (6, 'Sharon', 'Jenkins');

-- Vet Specialties
INSERT INTO vet_specialties (vet_id, specialty_id) VALUES (2, 1);
INSERT INTO vet_specialties (vet_id, specialty_id) VALUES (3, 2);
INSERT INTO vet_specialties (vet_id, specialty_id) VALUES (3, 3);
INSERT INTO vet_specialties (vet_id, specialty_id) VALUES (4, 2);
INSERT INTO vet_specialties (vet_id, specialty_id) VALUES (5, 1);

-- Owners
INSERT INTO owners (id, first_name, last_name, address, city, telephone, email) VALUES (1, 'George', 'Franklin', '110 W. Liberty St.', 'Madison', '6085551023', 'george.franklin@example.com');
INSERT INTO owners (id, first_name, last_name, address, city, telephone, email) VALUES (2, 'Betty', 'Davis', '638 Cardinal Ave.', 'Sun Prairie', '6085551749', 'betty.davis@example.com');
INSERT INTO owners (id, first_name, last_name, address, city, telephone, email) VALUES (3, 'Eduardo', 'Rodriquez', '2693 Commerce St.', 'McFarland', '6085558763', 'eduardo.rodriquez@example.com');
INSERT INTO owners (id, first_name, last_name, address, city, telephone, email) VALUES (4, 'Harold', 'Davis', '563 Friendly St.', 'Windsor', '6085553198', 'harold.davis@example.com');
INSERT INTO owners (id, first_name, last_name, address, city, telephone, email) VALUES (5, 'Peter', 'McTavish', '2387 S. Fair Way', 'Madison', '6085552765', 'peter.mctavish@example.com');
INSERT INTO owners (id, first_name, last_name, address, city, telephone, email) VALUES (6, 'Jean', 'Coleman', '105 N. Lake St.', 'Monona', '6085552654', 'jean.coleman@example.com');
INSERT INTO owners (id, first_name, last_name, address, city, telephone, email) VALUES (7, 'Jeff', 'Black', '1450 Oak Blvd.', 'Monona', '6085555387', 'jeff.black@example.com');
INSERT INTO owners (id, first_name, last_name, address, city, telephone, email) VALUES (8, 'Maria', 'Escobito', '345 Maple St.', 'Madison', '6085557683', 'maria.escobito@example.com');
INSERT INTO owners (id, first_name, last_name, address, city, telephone, email) VALUES (9, 'David', 'Schroeder', '2749 Blackhawk Trail', 'Madison', '6085559435', 'david.schroeder@example.com');
INSERT INTO owners (id, first_name, last_name, address, city, telephone, email) VALUES (10, 'Carlos', 'Estaban', '2335 Independence La.', 'Waunakee', '6085555487', 'carlos.estaban@example.com');
INSERT INTO owners (id, first_name, last_name, address, city, telephone, email) VALUES (11, 'William', 'Clown', '17 Circus Lane', 'Madison', '6085558899', 'william.clown@example.com');

-- Pets
INSERT INTO pets (id, name, birth_date, type_id, owner_id) VALUES (1, 'Leo', '2020-09-07', 1, 1);
INSERT INTO pets (id, name, birth_date, type_id, owner_id) VALUES (2, 'Basil', '2022-08-06', 6, 2);
INSERT INTO pets (id, name, birth_date, type_id, owner_id) VALUES (3, 'Rosy', '2021-04-17', 2, 3);
INSERT INTO pets (id, name, birth_date, type_id, owner_id) VALUES (4, 'Jewel', '2020-03-07', 2, 3);
INSERT INTO pets (id, name, birth_date, type_id, owner_id) VALUES (5, 'Iggy', '2020-11-30', 3, 4);
INSERT INTO pets (id, name, birth_date, type_id, owner_id) VALUES (6, 'George', '2020-01-20', 4, 5);
INSERT INTO pets (id, name, birth_date, type_id, owner_id) VALUES (7, 'Samantha', '2022-09-04', 1, 6);
INSERT INTO pets (id, name, birth_date, type_id, owner_id) VALUES (8, 'Max', '2022-09-04', 1, 6);
INSERT INTO pets (id, name, birth_date, type_id, owner_id) VALUES (9, 'Lucky', '2021-08-06', 5, 7);
INSERT INTO pets (id, name, birth_date, type_id, owner_id) VALUES (10, 'Mulligan', '2007-02-24', 2, 8);
INSERT INTO pets (id, name, birth_date, type_id, owner_id) VALUES (11, 'Camel', '2021-05-12', 7, 9);

-- Visits
INSERT INTO visits (id, pet_id, date, description) VALUES (1, 7, '2023-01-01', 'Rabies shot');
INSERT INTO visits (id, pet_id, date, description) VALUES (2, 8, '2023-03-04', 'Rabies shot');
INSERT INTO visits (id, pet_id, date, description) VALUES (3, 8, '2023-06-04', 'Neutered');
INSERT INTO visits (id, pet_id, date, description) VALUES (4, 7, '2023-09-04', 'Spayed');
