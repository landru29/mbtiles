-- +migrate Up

CREATE TABLE metadata (
    name text NOT NULL,
    value text NOT NULL
);

CREATE TABLE android_metadata (
    locale TEXT NOT NULL
);

CREATE TABLE tiles (
    zoom_level integer NOT NULL,
    tile_column integer NOT NULL,
    tile_row integer NOT NULL,
    tile_data blob
);

CREATE UNIQUE INDEX tile_index on tiles (zoom_level, tile_column, tile_row);

-- +migrate Down