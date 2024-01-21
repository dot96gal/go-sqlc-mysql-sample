-- name: GetPublisher :one
SELECT * FROM publishers
WHERE id = ? LIMIT 1;

-- name: ListPublishers :many
SELECT * FROM publishers
ORDER BY id;

-- name: CreatePublisher :execresult
INSERT INTO publishers (
  name
) VALUES (
  ?
);

-- name: UpdatePublisher :exec
UPDATE publishers
SET name = ?
WHERE id = ?;

-- name: DeletePublisher :exec
DELETE FROM publishers
WHERE id = ?;
