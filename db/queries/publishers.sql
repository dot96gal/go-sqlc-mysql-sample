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

-- name: GetPublisherBooks :many
SELECT
  p.id AS publisher_id,
  p.name AS publisher_name,
  b.id AS book_id,
  b.title AS book_title
FROM publishers AS p
INNER JOIN books AS b
ON p.id = b.publisher_id
WHERE p.id = ?;
