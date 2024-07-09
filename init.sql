CREATE TABLE games (
    id SERIAL primary key,
    code VARCHAR(100) NOT NULL,
    round INTEGER NULL,
    player_id INTEGER NOT NULL,
    game_type_id INTEGER NULL,
    hand_choice VARCHAR(100) NULL,
    created_at TIMESTAMP default now() NOT null updated_at TIMESTAMP default now() NOT null
);

CREATE TABLE game_types (
    id SERIAL primary key,
    name VARCHAR(100) not NULL,
);

CREATE TABLE players (
    id SERIAL primary key,
    username VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP default now() not null
);

INSERT INTO
    public.game_types("name")
VALUES
('hompimpa');

INSERT INTO
    public.game_types("name")
VALUES
('rock paper scissors');

CREATE TABLE codes (
    code VARCHAR(100) not NULL UNIQUE,
    created_at TIMESTAMP default now() not null,
    is_finished BOOLEAN default false host_id INTEGER --from players
    started_at VARCHAR(100) NULL,
    current_round INTEGER NULL
);

CREATE TABLE results (
    id SERIAL primary key,
    code VARCHAR(100) NOT NULL,
    hand_choice VARCHAR(100) NOT NULL,
    winner_player_ids JSON NULL,
    loser_player_ids JSON NULL,
    -- disqualification_player_ids JSON NULL,
    round INTEGER NULL,
    created_at TIMESTAMP default now() not NULL
);