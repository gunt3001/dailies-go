-- name: GetEntriesByDateRange :many
SELECT *
FROM Entries
WHERE Date >= @date_start AND Date <= @date_end
ORDER BY Date DESC;