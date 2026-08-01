-- test users
INSERT INTO users (username, password_hash, trophies, base_level, xp, last_played)
VALUES 
('testuser1', '$2a$10$z07GAxr.3lxkDylbyDiQ7.WU2aoo/BrgZDMEm4/LfwVd7s8W40JzW', 0, 4, 0, NOW()),
('testuser2', '$2a$10$sZ9I541mbOtx6DCpo./PeujEz9D2ywWGuSt25w5fyFLqI/Xk9c4ie', 0, 4, 0, NOW()),
('testuser3', '$2a$10$7jZCR5EeuOnUQh0e8A6NgeOyColn4AK5wZmf/ZRF411zQgkviFrNq', 0, 4, 0, NOW()),
('testuser4', '$2a$10$Z8sxwu4o.pwIcFwuVkOrqe8EAGOAwQb0H1Yu73UyrbdVJ/xFFe1NK', 0, 4, 0, NOW()),
('testuser5', '$2a$10$95/Jwo341FbyS5j0lzrHbeSEumaWiPeQoSVYaWn8dicHrcYN3Tweu', 0, 4, 0, NOW());

-- get userIDs for test users (assuming they get IDs after your seeded troops/buildings)
-- city details with max resources
INSERT INTO city_details (user_id, resource1, resource2, max_resource1, max_resource2, max_troop_army_size, last_updated, max_defence_buildings, max_resource_buildings)
SELECT user_id, 10000, 10000, 10000, 10000, 100, NOW(), 10, 10
FROM users WHERE username LIKE 'testuser%';

-- battle history
INSERT INTO user_battle_history (user_id, number_of_battles, battles_won, battles_lost, trophies)
SELECT user_id, 0, 0, 0, 0
FROM users WHERE username LIKE 'testuser%';

-- army details
INSERT INTO army_details (user_id, troop_units_used, army_composition, created_on)
SELECT user_id, 0, '[]', NOW()
FROM users WHERE username LIKE 'testuser%';

-- unlock all troops for test users
INSERT INTO user_troop_details (user_id, troop_id, troop_level)
SELECT u.user_id, td.troop_id, 3
FROM users u
CROSS JOIN troop_details td
WHERE u.username LIKE 'testuser%';

-- place all defense buildings for test users in a grid pattern
INSERT INTO user_buildings (user_id, building_id, level, grid_x, grid_y)
SELECT u.user_id, bd.building_id, 3,
    2 + (ROW_NUMBER() OVER (PARTITION BY u.user_id ORDER BY bd.building_id) * 5),
    10
FROM users u
CROSS JOIN building_details bd
WHERE u.username LIKE 'testuser%'
AND bd.building_type = 'defense';

-- place town hall for test users
INSERT INTO user_buildings (user_id, building_id, level, grid_x, grid_y)
SELECT u.user_id, bd.building_id, 4, 25, 25
FROM users u
CROSS JOIN building_details bd
WHERE u.username LIKE 'testuser%'
AND bd.name = 'Town Hall';