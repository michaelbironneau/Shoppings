package main

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"net/url"
	"time"
)

type LogLevel int

const (
	LogLevelMin    = 1
	LogLevelMedium = 8
	LogLevelMax    = 32
)

func NewDB(server, username, password, database string, logLev LogLevel) (*sql.DB, error) {
	connStr := formatDBConnString(server, username, password, database, int(logLev))
	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return nil, err
	}
	// Fixes for idle connection drops and optimisation for connection re-use
	// https://github.com/denisenkom/go-mssqldb/issues/167
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(20 * time.Minute) //prevent dead connections as Azure SQL Database closes connections after 30mins
	return db, nil
}


// https://github.com/denisenkom/go-mssqldb#the-connection-string-can-be-specified-in-one-of-three-formats
func formatDBConnString(server, username, password, database string, logLev int) string {
	mssqlQueryParams := url.Values{}
	mssqlQueryParams.Add("database", database)

	mssqlUrl := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(username, password),
		Host:     server,
		RawQuery: mssqlQueryParams.Encode(),
	}

	return mssqlUrl.String()
}
