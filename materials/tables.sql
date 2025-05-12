CREATE TABLE users (
    id UUID PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE games (
	id UUID PRIMARY KEY,
	field TEXT
);

ALTER TABLE games ADD COLUMN vs_computer Boolean DEFAULT TRUE;
ALTER TABLE games ADD COLUMN status INTEGER DEFAULT 0;
ALTER TAble games ADD COLUMN Playerx text;
ALTER TAble games ADD COLUMN Playero text;
