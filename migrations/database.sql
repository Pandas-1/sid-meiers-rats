CREATE TABLE users (
    user_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    trophies INTEGER NOT NULL DEFAULT 0,
    base_level INTEGER NOT NULL DEFAULT 1,
    xp INTEGER NOT NULL DEFAULT 0,
    last_played TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    password_hash TEXT NOT NULL
);

CREATE TABLE city_details (
    user_id INTEGER NOT NULL,
    resource1 BIGINT NOT NULL DEFAULT 0,
    resource2 BIGINT NOT NULL DEFAULT 0,
    max_resource1 BIGINT NOT NULL DEFAULT 500,
    max_resource2 BIGINT NOT NULL DEFAULT 500,
    max_troop_army_size BIGINT NOT NULL DEFAULT 20,
    last_updated TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    max_defence_buildings BIGINT NOT NULL DEFAULT 3,
    max_resource_buildings BIGINT NOT NULL DEFAULT 3,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE troop_details (
    troop_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    base_cost INTEGER NOT NULL,
    troop_attack_power INTEGER NOT NULL,
    building_attack_power INTEGER NOT NULL,
    defence INTEGER NOT NULL,
    range INTEGER NOT NULL,
    attribute_strength INTEGER NOT NULL,
    attribute_weakness INTEGER NOT NULL,
    troop_army_space INTEGER NOT NULL,
    movement_speed INTEGER NOT NULL,
    max_level INTEGER[] NOT NULL
);

CREATE TABLE building_details (
    building_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    building_type VARCHAR(50) NOT NULL,
    production INTEGER NOT NULL DEFAULT 0,
    scaling INTEGER NOT NULL DEFAULT 0,
    health_bar INTEGER NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    defence_attack INTEGER NOT NULL DEFAULT 0,
    defence_range INTEGER NOT NULL DEFAULT 0,
    max_level INTEGER[] NOT NULL,
    cost_resource1 INTEGER NOT NULL DEFAULT 0,
    cost_resource2 INTEGER NOT NULL DEFAULT 0,
    element_type INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE user_buildings (
    instance_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER NOT NULL,
    building_id INTEGER NOT NULL,
    level BIGINT NOT NULL DEFAULT 1,
    grid_x INTEGER NOT NULL,
    grid_y INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (building_id) REFERENCES building_details(building_id)
);

CREATE TABLE user_troop_details (
    user_id INTEGER NOT NULL,
    troop_id INTEGER NOT NULL,
    troop_level INTEGER NOT NULL DEFAULT 1,
    PRIMARY KEY (user_id, troop_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (troop_id) REFERENCES troop_details(troop_id)
);

CREATE TABLE army_details (
    user_id INTEGER NOT NULL PRIMARY KEY,
    troop_units_used INTEGER NOT NULL DEFAULT 0,
    army_composition jsonb NOT NULL DEFAULT '[]',
    created_on TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE user_battle_history (
    user_id INTEGER NOT NULL PRIMARY KEY,
    number_of_battles INTEGER NOT NULL DEFAULT 0,
    battles_won INTEGER NOT NULL DEFAULT 0,
    battles_lost INTEGER NOT NULL DEFAULT 0,
    trophies INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE battles (
    battle_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    attacker_id INTEGER NOT NULL,
    defender_id INTEGER NOT NULL,
    resource1_won INTEGER NOT NULL DEFAULT 0,
    resource2_won INTEGER NOT NULL DEFAULT 0,
    victory_percentage INTEGER NOT NULL,
    fought_at TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (attacker_id) REFERENCES users(user_id),
    FOREIGN KEY (defender_id) REFERENCES users(user_id)
);