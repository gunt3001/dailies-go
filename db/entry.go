// Interface implementation for Entries type
package db

import (
	"fmt"
	"time"
)

func (e Entry) FilterValue() string {
	return e.Content
}

func (e Entry) Title() string {
	date, err := e.getEntryDateAsTime()
	if err != nil {
		return "Invalid entry date"
	}
	dateString := date.Format("2 Jan '06")
	return fmt.Sprintf("%s (%s) - %s", dateString, e.getRelativeDateString(), e.Keyword)
}

func (e Entry) Description() string {
	return e.Content
}

func (e Entry) getEntryDateAsTime() (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", e.Date)
}

func (e Entry) getRelativeDateString() string {
	date, err := e.getEntryDateAsTime()
	if err != nil {
		return "?"
	}
	now := time.Now()
	diff := now.Sub(date)
	days := int(diff.Hours() / 24)
	return fmt.Sprintf("%dd ago", days)
}
