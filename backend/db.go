package main

import (
	"database/sql"
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

func newTestDB() (*sql.DB, error) {
	return sql.Open("sqlserver", "database=opdb;log=16")
}

func resetTestDB(db *sql.DB) error {
	_, err := db.Exec(`USE OpDb
		GO
		TRUNCATE TABLE Security.Principal 
		GO
		EXECUTE [Security].[uspAddPrincipal] 
   		'TestUser'
  		,'TestPass'
		GO
		DELETE FROM App.[List]
		GO
		TRUNCATE TABLE App.StoreOrder 
		GO
		DELETE FROM App.Store 
		GO
		DELETE FROM App.Item 
		GO `)
	return err
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
