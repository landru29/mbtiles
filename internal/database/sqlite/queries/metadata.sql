-- name: InsertMetadata :exec
INSERT INTO metadata(name, value) VALUES (?, ?);

-- name: Metadata :many
SELECT name, value FROM metadata;

-- name: WipeAllMetadata :exec
DELETE FROM metadata;