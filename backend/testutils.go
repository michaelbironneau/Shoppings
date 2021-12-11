package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"io"
	"io/ioutil"
)

func newTestDB() (*sql.DB, error) {
	return sql.Open("sqlserver", "database=opdb;log=32")
}

func resetTestDB(db *sql.DB) error {
	_, err := db.Exec(`
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
		DELETE FROM App.Item `)
	return err
}

func SetupTestAPI() (*fiber.App, error) {
	app := 	fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: HandleError,
	})
	db, err := newTestDB()
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err 
	}
	if err := resetTestDB(db); err != nil {
		return nil, err
	}
	registerHandlers(app, db)
	return app, nil
}

func unmarshalBody(body io.ReadCloser, i interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.Unmarshal(b, &i)
}