package queries

const (
	ListAttributeValueGraphByAttributeIDsSQL = `
SELECT
	a.id AS attribute_id,
	v.id AS value_id,
	v.slug AS value_slug,
	v.value AS value_value,
	v.created_at AS value_created_at,
	v.updated_at AS value_updated_at,
	r.from_attribute_id,
	r.from_value_id,
	r.to_attribute_id,
	r.to_value_id,
	p_to.slug AS to_attribute_slug,
	v_to.slug AS to_value_slug
FROM strategy_attribute a
LEFT JOIN strategy_attribute_value v
	ON v.strategy_attribute_id = a.id
LEFT JOIN strategy_attribute_value_relation r
	ON r.from_value_id = v.id
LEFT JOIN strategy_attribute p_to
	ON p_to.id = r.to_attribute_id
	AND p_to.strategy_id = a.strategy_id
LEFT JOIN strategy_attribute_value v_to
	ON v_to.id = r.to_value_id
WHERE a.strategy_id = $1
  AND a.id = ANY($2::uuid[])
ORDER BY a.created_at ASC, v.created_at ASC, r.created_at ASC`

	ListRelationsByFromValueIDsSQL = `
SELECT
	r.from_attribute_id, r.from_value_id, r.to_attribute_id, r.to_value_id,
	p_to.slug AS to_attribute_slug, v_to.slug AS to_value_slug
FROM strategy_attribute_value_relation r
JOIN strategy_attribute p_from ON p_from.id = r.from_attribute_id
JOIN strategy_attribute p_to ON p_to.id = r.to_attribute_id
JOIN strategy_attribute_value v_to ON v_to.id = r.to_value_id
WHERE p_from.strategy_id = $1
  AND p_to.strategy_id = $1
  AND r.from_value_id = ANY($2::uuid[])
ORDER BY r.from_value_id ASC, r.created_at ASC`

	DeleteRelationsByFromValueSQL = `
DELETE FROM strategy_attribute_value_relation r
USING strategy_attribute p_from
WHERE r.from_attribute_id = p_from.id
  AND p_from.strategy_id = $1
  AND r.from_attribute_id = $2
  AND r.from_value_id = $3`

	InsertRelationSQL = `
INSERT INTO strategy_attribute_value_relation (
	id, from_attribute_id, from_value_id, to_attribute_id, to_value_id, created_at, updated_at
) VALUES (gen_random_uuid(), $1, $2, $3, $4, NOW(), NOW())
ON CONFLICT (from_value_id, to_value_id) DO NOTHING`
)
