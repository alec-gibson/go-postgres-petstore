-- name: FindPetByID :one
SELECT * FROM pets
WHERE id = $1 LIMIT 1;

-- name: ListPets :many
SELECT * FROM pets
WHERE cardinality($1::varchar[]) = 0 OR tag = ANY($1::varchar[])
ORDER BY name;

-- name: ListPetsWithLimit :many
SELECT * FROM pets
WHERE cardinality($1::varchar[]) = 0 OR tag = ANY($1::varchar[])
ORDER BY name LIMIT $2;

-- name: CreatePet :one
INSERT INTO pets (
  name, tag
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeletePet :exec
DELETE FROM pets
WHERE id = $1;
