package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/michaelbironneau/shoppings/backend/api"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func newTestDB() (*sql.DB, error) {
	return sql.Open("sqlserver", "database=opdb;log=32")
}

func resetTestDB(db *sql.DB) error {
	_, err := db.Exec(`
		TRUNCATE TABLE Security.Principal;
		EXECUTE [Security].[uspAddPrincipal] 
   		'TestUser'
  		,'TestPass';
		DELETE FROM App.[List];
		TRUNCATE TABLE App.StoreOrder;
		DELETE FROM App.Store;
		DELETE FROM App.Item `)
	return err
}

func SetupTestAPI() (*fiber.App, string, error) {
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: HandleError,
	})
	db, err := newTestDB()
	if err != nil {
		return nil, "", err
	}
	if err := db.Ping(); err != nil {
		return nil, "", err
	}
	if err := resetTestDB(db); err != nil {
		return nil, "", err
	}
	registerHandlers(app, db)
	token, err := getTestToken(db)
	return app, token, err
}

func unmarshalBody(body io.ReadCloser, i interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.Unmarshal(b, &i)
}

func makeTestRequest(app *fiber.App, method, route, token string, i interface{}) (*http.Response, error) {
	var r *http.Request
	if i == nil {
		r, _ := http.NewRequest(method, route, nil)
		r.Header.Add(api.TokenHeader, token)
		return app.Test(r, -1)
	}
	switch i.(type) {
	case string:
		r, _ = http.NewRequest(method, route, strings.NewReader(i.(string)))
	default:
		b, err := json.Marshal(i)
		if err != nil {
			return nil, err
		}
		r, _ = http.NewRequest(method, route, bytes.NewReader(b))
	}
	r.Header.Add(api.TokenHeader, token)
	return app.Test(r, -1)
}

func getTestToken(db *sql.DB) (string, error) {
	row := db.QueryRow(`SELECT TOP 1 Token From Security.Principal`)
	var token string
	if err := row.Scan(&token); err != nil {
		return "", err
	}
	return token, nil
}
