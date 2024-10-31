-- name: GetSight :one
SELECT
    sight_id,
    sight_name,
    sight_description
FROM
    sights
WHERE
    sight_id = $1;

-- name: GetSights :many
SELECT
    sight_id,
    sight_name,
    sight_description
FROM
    sights;

