// Interface implementation for Entries type
package db

func (e Entry) FilterValue() string {
	return e.Content
}

func (e Entry) Title() string {
	return e.Date
}

func (e Entry) Description() string {
	return e.Content
}
