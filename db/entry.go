// Interface implementation for Entries type
package db

import (
	"time"
)

func (e Entry) Title() string {
	if e.Keyword.Valid {
		return e.Keyword.String
	}
	return ""
}

func (e Entry) Description() string {
	if e.Content.Valid {
		return e.Content.String
	}
	return ""
}

func (e Entry) getEntryDateAsTime() (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", e.Date)
}
