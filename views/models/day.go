// Represents a day in the calendar to be displayed in a list, with or without an associated journal entry
package models

import (
	"dailies-go/db"
	"fmt"
	"time"
)

type Day struct {
	date  time.Time
	entry *db.Entry
}

// Constructor
func NewDay(date time.Time, entry *db.Entry) Day {
	return Day{
		date:  date,
		entry: entry,
	}
}

// Implement list item interface

func (d Day) FilterValue() string {
	if d.entry != nil && d.entry.Content.Valid {
		return d.entry.Content.String
	}
	return ""
}

func (d Day) Title() string {
	var caption string
	if d.entry == nil {
		caption = "No entry"
	} else {
		caption = d.entry.Title()
	}

	dateString := d.date.Format("2 Jan '06")
	return fmt.Sprintf("%s (%s) - %s", dateString, d.getRelativeDateString(), caption)
}

func (d Day) Description() string {
	if d.entry == nil {
		return "No entry"
	} else {
		return d.entry.Description()
	}
}

func (d Day) getRelativeDateString() string {
	now := time.Now()
	diff := now.Sub(d.date)
	days := int(diff.Hours() / 24)
	if days == 0 {
		return "Today"
	} else if days == 1 {
		return "Yesterday"
	} else if days < 0 {
		return fmt.Sprintf("in %dd", days)
	} else {
		return fmt.Sprintf("%dd ago", days)
	}
}
