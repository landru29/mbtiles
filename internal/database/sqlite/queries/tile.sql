-- name: TileCount :one
SELECT count(*) FROM tiles;

-- name: TileByIndex :one
SELECT
    zoom_level,
    tile_column,
    tile_row,
    tile_data
FROM tiles LIMIT 1 OFFSET @index;

-- name: TileByCoordinate :one
SELECT 
    zoom_level, 
    tile_column, 
    tile_row, 
    tile_data
FROM tiles
WHERE tile_column=@col AND tile_row=@row AND zoom_level=@zoomLevel;

-- name: Tiles :many
SELECT zoom_level, tile_column, tile_row, tile_data FROM tiles;

-- name: TileDataUpdate :exec
UPDATE tiles 
SET tile_data=? 
WHERE tile_column=@col AND tile_row=@row AND zoom_level=@zoomLevel;
