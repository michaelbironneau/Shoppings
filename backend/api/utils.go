package api

import (
	"database/sql"
	"log"
)

var (
	Err401 = &Error{Code: 401, Message: "Unauthorized. Please check username and password."}
	Err500 = &Error{Code: 500, Message: "Unexpected error. Please try again later."}
)

type Error struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e *Error) Error() string {
	return e.Message
}

func getPrincipal(db *sql.DB, token string) (string, error) {
	if token == "" {
		return "", Err401
	}
	var res sql.NullString
	row := db.QueryRow("SELECT TOP 1 Username FROM Security.Principal WHERE Token = @InToken", sql.Named("InToken", token))
	if err := row.Scan(&res); err != nil && err != sql.ErrNoRows {
		log.Printf("Error retrieving principal from credentials: %v", err)
		return "", Err500
	} else if err == sql.ErrNoRows {
		log.Printf("Invalid token")
		return "", Err401
	}
	if !res.Valid {
		log.Printf("Invalid token: %v", token)
		return "", Err401
	}
	return res.String, nil
}

func dbError(err error) error {
	if err == nil {
		return nil
	}
	log.Printf("Database error: %v", err)
	return Err500
}
