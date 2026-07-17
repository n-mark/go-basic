-- +goose Up
-- +goose StatementBegin
CREATE TABLE players (
    id            BIGSERIAL PRIMARY KEY,
    name          TEXT   NOT NULL,
    figures_color TEXT   NOT NULL CHECK (figures_color IN ('white', 'black'))
);

CREATE TABLE figures (
    id         BIGSERIAL PRIMARY KEY,
    symbol     TEXT   NOT NULL,
    piece_color TEXT  NOT NULL
);

CREATE TABLE moves (
    id            BIGSERIAL PRIMARY KEY,
    time_took  BIGINT NOT NULL,
    position_from TEXT   NOT NULL,
    position_to   TEXT   NOT NULL,
    figure        TEXT   NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS moves;
DROP TABLE IF EXISTS figures_taken;
DROP TABLE IF EXISTS players;
-- +goose StatementEnd
