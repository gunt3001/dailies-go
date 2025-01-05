-- name: GetEntriesByDateRange :many
SELECT *
FROM Entries
WHERE Date >= @date_start AND Date <= @date_end
ORDER BY Date DESC;

-- name: GetRandomEntry :one
SELECT *
FROM Entries
WHERE Date IN (
    SELECT Date 
    FROM Entries 
    WHERE Content IS NOT NULL AND Remarks != ''
    ORDER BY RANDOM() LIMIT 1
);