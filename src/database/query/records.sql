-- name: GetRecord :one
SELECT
    image_id,
    sight_id,
    sight_name,
    copywriting
FROM
    records
WHERE
    record_id = $1;

-- name: GetRandomRecord :one
SELECT
    record_id,
    image_id,
    sight_id,
    sight_name,
    copywriting
FROM
    records
WHERE
    record_id =(
        SELECT
            record_id
        FROM
            records OFFSET floor(random() *(
                SELECT
                    COUNT(*)
                FROM records))
        LIMIT 1);

-- name: CreateRecord :one
INSERT INTO records(image_id, sight_id, sight_name, copywriting)
    VALUES ($1, $2, $3, $4)
RETURNING
    record_id;

