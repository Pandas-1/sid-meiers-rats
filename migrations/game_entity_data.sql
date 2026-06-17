-- clear existing seed data
TRUNCATE troop_details CASCADE;
TRUNCATE building_details CASCADE;

-- troops
INSERT INTO troop_details (name, base_cost, troop_attack_power, building_attack_power, defence, range, attribute_strength, attribute_weakness, troop_army_space, movement_speed, max_level)
VALUES 
('Burning Rat', 20, 10, 10, 10, 2, 4, 3, 1, 2, ARRAY[3,4,5,7]),
('Wet Rat',     20, 10, 10,  7, 2, 2, 5, 1, 2, ARRAY[3,4,5,7]),
('Flying Rat',  20, 10, 10, 10, 5, 5, 2, 5, 3, ARRAY[3,4,5,7]),
('Ground Rat',  20, 10, 10, 13, 1, 3, 4, 5, 1, ARRAY[3,4,5,7]),
('Bright Rat',  35, 15,  8, 10, 1, 0, 0, 1, 4, ARRAY[3,4,5,7]),
('Dark Rat',    35,  8, 15, 13, 3, 0, 0, 4, 1, ARRAY[3,4,5,7]);

-- buildings
INSERT INTO building_details (name, building_type, production, scaling, health_bar, width, height, defence_attack, defence_range, max_level, cost_resource1, cost_resource2)
VALUES
('Town Hall',        'base',     10, 0, 500, 4, 4,   0,  0, ARRAY[4,4,4,4], 2000, 2000),
('Fire Tower',       'defense',  0, 2, 100, 2, 2, 100,  8, ARRAY[3,5,7,9],    0,  200),
('Water Faucet',     'defense',  0, 2, 100, 2, 2, 100,  5, ARRAY[3,5,7,9],    0,  200),
('Air Missiles',     'defense',  0, 2, 100, 2, 2, 100, 10, ARRAY[3,5,7,9],    0,  200),
('Ground Cannons',   'defense',  0, 2, 100, 2, 2, 100,  5, ARRAY[3,5,7,9],    0,  200),
('Gold Mine',        'resource', 10, 2,  75, 2, 2,   0,  0, ARRAY[3,5,7,9],    0,  500),
('Elixir Collector', 'resource', 10, 2,  75, 2, 2,   0,  0, ARRAY[3,5,7,9],  500,    0);
