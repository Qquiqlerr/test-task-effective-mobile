-- +goose Up
-- +goose StatementBegin
CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    "group" VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    text TEXT NOT NULL,
    link TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE song;
-- +goose StatementEnd
