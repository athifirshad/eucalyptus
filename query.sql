-- name: GetHealthRecord :one
SELECT *
FROM health_record
WHERE patient_id = $1;