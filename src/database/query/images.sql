-- name: GetImage :one
SELECT
    image_path
FROM
    images
WHERE
    image_id = $1;

-- name: CreateImage :one
INSERT INTO images(image_path)
    VALUES ($1)
RETURNING
    image_id;

-- name: UpdateImage :exec
UPDATE
    images
SET
    image_path = $2
WHERE
    image_id = $1;

