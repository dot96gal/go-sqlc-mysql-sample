-- name: GetAuthorBook :one
SELECT * FROM author_books
WHERE author_id = ? AND book_id = ? LIMIT 1;

-- name: ListAuthorBooks :many
SELECT
  a.id AS author_id,
  a.name AS author_name,
  a.bio AS author_bio,
  b.id AS book_id,
  b.title AS book_title
FROM authors AS a
INNER JOIN author_books AS ab
ON a.id = ab.author_id
INNER JOIN books AS b
ON ab.book_id = b.id
ORDER BY a.id, b.id;

-- name: CreateAuthorBook :exec
INSERT INTO author_books (
  author_id, book_id 
) VALUES (
  ?, ?
);

-- name: DeleteAuthorBook :exec
DELETE FROM author_books
WHERE author_id = ? AND book_id = ?;
