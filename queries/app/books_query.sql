-- name: InsertBook :one
INSERT INTO books(id, name, description)
VALUES (@id, @name, @description)
RETURNING *;
-- name: SelectBookWhereID :one
SELECT *
FROM books
WHERE id = @id;
-- name: SelectBooksCount :one
SELECT COUNT(*)
FROM books
WHERE name LIKE @name
    AND description LIKE @description;
-- name: SelectBooks :many
SELECT *
FROM books
WHERE name LIKE @name
    AND description LIKE @description
ORDER BY created_at DESC
LIMIT @lim OFFSET @ofst;
-- name: UpdateBookWhereID :exec
UPDATE books
SET name = @name,
    description = @description,
    updated_at = now()
WHERE id = @id;
-- name: DeleteBookWhereID :exec
DELETE FROM books
WHERE id = @id;