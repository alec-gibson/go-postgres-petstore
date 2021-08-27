-- name: FindPetByID :one
SELECT * FROM petstore.pets
WHERE id = $1 LIMIT 1;

-- name: ListPets :many
SELECT * FROM petstore.pets
WHERE cardinality(@tags::varchar[]) = 0 OR tag = ANY(@tags::varchar[])
ORDER BY name;

-- name: ListPetsWithLimit :many
SELECT * FROM petstore.pets
WHERE cardinality(@tags::varchar[]) = 0 OR tag = ANY(@tags::varchar[])
ORDER BY name LIMIT @max_records;

-- name: CreatePet :one
INSERT INTO petstore.pets (
  name, tag
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeletePet :exec
DELETE FROM petstore.pets
WHERE id = $1;
