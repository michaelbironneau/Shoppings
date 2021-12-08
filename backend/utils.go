package main

import (
	"database/sql"
	"log"
)

func getPrincipal(db *sql.DB, token string) (string, error) {
	var res sql.NullString
	row := db.QueryRow("SELECT TOP 1 Username FROM Security.Principal WHERE Token = @InToken", sql.Named("InToken", token))
	if err := row.Scan(&res); err != nil {
		log.Printf("Error retrieving principal from credentials: %v", err)
		return "", Err500
	}
	if !res.Valid {
		log.Printf("Invalid token: %v", token)
		return "", Err401
	}
	return res.String, nil
}
