ALTER TABLE players ADD COLUMN team_id INT;

ALTER TABLE players
ADD CONSTRAINT fk_team
FOREIGN KEY (team_id) REFERENCES teams(id);
