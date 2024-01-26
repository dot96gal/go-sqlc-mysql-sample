-- name: GetPublisher :one
SELECT
  *
FROM
  publishers
WHERE
  uuid = ?
LIMIT
  1;

-- name: ListPublishers :many
SELECT
  *
FROM
  publishers
ORDER BY
  uuid;

-- name: CreatePublisher :exec
INSERT INTO
  publishers (uuid, name)
VALUES
  (?, ?);

-- name: UpdatePublisher :exec
UPDATE publishers
SET
  name = ?
WHERE
  uuid = ?;

-- name: DeletePublisher :exec
DELETE FROM publishers
WHERE
  uuid = ?;

-- name: GetPublisherBooks :many
SELECT
  p.uuid AS publisher_uuid,
  p.name AS publisher_name,
  b.uuid AS book_uuid,
  b.title AS book_title
FROM
  publishers AS p
  INNER JOIN books AS b ON p.uuid = b.publisher_uuid
WHERE
  p.uuid = ?
ORDER BY
  p.uuid,
  b.uuid;
