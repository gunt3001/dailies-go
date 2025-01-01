// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: sqlc-query.sql

package db

import (
	"context"
)

const getEntriesByDateRange = `-- name: GetEntriesByDateRange :many
SELECT date, content, keyword, mood, remarks
FROM Entries
WHERE Date >= ?1 AND Date <= ?2
ORDER BY Date DESC
`

type GetEntriesByDateRangeParams struct {
	DateStart string
	DateEnd   string
}

func (q *Queries) GetEntriesByDateRange(ctx context.Context, arg GetEntriesByDateRangeParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, getEntriesByDateRange, arg.DateStart, arg.DateEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.Date,
			&i.Content,
			&i.Keyword,
			&i.Mood,
			&i.Remarks,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
