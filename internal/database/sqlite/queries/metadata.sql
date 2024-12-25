-- name: InsertMetadata :exec
INSERT INTO metadata(name, value) VALUES (?, ?);

-- name: Metadata :many
SELECT name, value FROM metadata;

-- name: WipeAllMetadata :exec
DELETE FROM metadata;


-- name: UpdateMetadata :exec
UPDATE metadata SET value=? WHERE name=:name;