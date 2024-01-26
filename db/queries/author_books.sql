-- name: GetAuthorBook :one
SELECT
  *
FROM
  author_books
WHERE
  author_uuid = ?
  AND book_uuid = ?
LIMIT
  1;

-- name: ListAuthorBooks :many
SELECT
  a.uuid AS author_uuid,
  a.name AS author_name,
  a.bio AS author_bio,
  b.uuid AS book_uuid,
  b.title AS book_title
FROM
  authors AS a
  INNER JOIN author_books AS ab ON a.uuid = ab.author_uuid
  INNER JOIN books AS b ON ab.book_uuid = b.uuid
ORDER BY
  a.uuid,
  b.uuid;

-- name: CreateAuthorBook :exec
INSERT INTO
  author_books (author_uuid, book_uuid)
VALUES
  (?, ?);

-- name: DeleteAuthorBook :exec
DELETE FROM author_books
WHERE
  author_uuid = ?
  AND book_uuid = ?;
