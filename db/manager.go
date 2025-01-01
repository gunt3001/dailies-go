package db

import (
	"database/sql"
	"fmt"
)

type JournalManager struct {
}

func InitJournalManager(fileName string) (*Queries, error) {

	dbConn, err := sql.Open("sqlite", fileName)
	if err != nil {
		fmt.Println("Error opening database connection: ", err)
		return nil, err
	}
	queries := New(dbConn)
	return queries, nil
}
