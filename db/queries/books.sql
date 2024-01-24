-- name: GetBook :one
SELECT * FROM books 
WHERE id = ? LIMIT 1;

-- name: ListBooks :many
SELECT * FROM books
ORDER BY id;

-- name: CreateBook :execresult
INSERT INTO books (
  title,
  publisher_id
) VALUES (
  ?, ?
);

-- name: UpdateBook :exec
UPDATE books
SET title = ?
WHERE id = ?;

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = ?;

-- name: GetBookPublisher :one
SELECT
  b.id AS book_id,
  b.title AS book_title,
  p.id AS publisher_id,
  p.name AS publisher_name
FROM books AS b
INNER JOIN publishers AS p
ON b.publisher_id = p.id
WHERE b.id = ? LIMIT 1;
