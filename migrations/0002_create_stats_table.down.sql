CREATE TABLE stats (
    id SERIAL PRIMARY KEY,
    player_id INTEGER NOT NULL REFERENCES players(id),
    points INTEGER NOT NULL CHECK (points >= 0),
    rebounds INTEGER NOT NULL CHECK (rebounds >= 0),
    assists INTEGER NOT NULL CHECK (assists >= 0),
    steals INTEGER NOT NULL CHECK (steals >= 0),
    blocks INTEGER NOT NULL CHECK (blocks >= 0),
    fouls INTEGER NOT NULL CHECK (fouls >= 0 AND fouls <= 6),
    turnovers INTEGER NOT NULL CHECK (turnovers >= 0),
    minutes_played FLOAT NOT NULL CHECK (minutes_played >= 0 AND minutes_played <= 48.0),
    game_date DATE NOT NULL
);

DROP TABLE stats;
