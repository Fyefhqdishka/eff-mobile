-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS songs (
    id serial PRIMARY KEY,
    group_name varchar(255) NOT NULL REFERENCES groups(name) ON DELETE CASCADE,
    song varchar(255) NOT NULL,
    text text NOT NULL,
    link varchar(255) NOT NULL,
    releaseDate varchar(255) NOT NULL
);

CREATE INDEX idx_group_name ON songs(group_name);
CREATE INDEX idx_song_name ON songs(song);
CREATE INDEX idx_release_date ON songs(releaseDate);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS songs;

DROP INDEX idx_group_name;
DROP INDEX idx_song_name;
DROP INDEX idx_release_date;

-- +goose StatementEnd
