-- name: GetAuthor :one
SELECT
  *
FROM
  authors
WHERE
  uuid = ?
LIMIT
  1;

-- name: ListAuthors :many
SELECT
  *
FROM
  authors
ORDER BY
  uuid;

-- name: CreateAuthor :exec
INSERT INTO
  authors (uuid, name, bio)
VALUES
  (?, ?, ?);

-- name: UpdateAuthor :exec
UPDATE authors
SET
  name = ?,
  bio = ?
WHERE
  uuid = ?;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE
  uuid = ?;
