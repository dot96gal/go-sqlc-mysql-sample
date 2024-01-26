-- name: GetBook :one
SELECT
  *
FROM
  books
WHERE
  uuid = ?
LIMIT
  1;

-- name: ListBooks :many
SELECT
  *
FROM
  books
ORDER BY
  uuid;

-- name: CreateBook :exec
INSERT INTO
  books (uuid, title, publisher_uuid)
VALUES
  (?, ?, ?);

-- name: UpdateBook :exec
UPDATE books
SET
  title = ?
WHERE
  uuid = ?;

-- name: DeleteBook :exec
DELETE FROM books
WHERE
  uuid = ?;

-- name: GetBookPublisher :one
SELECT
  b.uuid AS book_uuid,
  b.title AS book_title,
  p.uuid AS publisher_uuid,
  p.name AS publisher_name
FROM
  books AS b
  INNER JOIN publishers AS p ON b.publisher_uuid = p.uuid
WHERE
  b.uuid = ?
LIMIT
  1;
