package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Error struct {
	Message string `json:"error"`
}

func (e *Error) Error() string {
	return e.Message
}

var (
	Err401 = &Error{"Unauthorized. Please check username and password."}
	Err500 = &Error{"Unexpected error. Please try again later."}
)

func registerHandlers(app *fiber.App, db *sql.DB) {
	app.Post("/token", func(c *fiber.Ctx) error {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return Err401
		}
		row := db.QueryRow("SELECT [Security].[udf_GetToken] (@InUsername, @InPassword)", sql.Named("InUsername", req.Username), sql.Named("InPassword", req.Password))
		var token sql.NullString
		if err := row.Scan(&token); err != nil {
			log.Printf("Error retrieving token: %v\n")
			return Err401
		}
		if !token.Valid {
			log.Printf("Invalid login for user %s\n", req.Username)
			return Err401
		}
		return c.JSON(struct {
			Token string `json:"token"`
		}{token.String})
	})
	app.Get("/lists", func(c *fiber.Ctx) error {
		_, err := getPrincipal(db, string(c.Request().Header.Peek("X-Token")))
		if err != nil {
			return err
		}

		return nil
	})
	app.Post("/lists", func(c *fiber.Ctx) error {
		return nil
	})
	app.Post("/lists/:id/archive", func(c *fiber.Ctx) error {
		return nil
	})
	app.Post("/lists/:id/unarchive", func(c *fiber.Ctx) error {
		return nil
	})
	app.Get("/lists/:id", func(c *fiber.Ctx) error {
		return nil
	})
	app.Get("/lists/:id/updates/:since", func(c *fiber.Ctx) error {
		return nil
	})
	app.Post("/lists/:id/updates", func(c *fiber.Ctx) error {
		return nil
	})
	app.Post("/lists/:id/items/complete/:item", func(c *fiber.Ctx) error {
		return nil
	})
	app.Get("/item-search/:needle", func(c *fiber.Ctx) error {
		return nil
	})
}

func main() {
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	app := fiber.New()
	var db *sql.DB
	for {
		db, err = NewDB(viper.GetString("server"), viper.GetString("username"), viper.GetString("password"), viper.GetString("database"), LogLevelMedium)
		if err != nil {
			log.Printf("Error connecting to DB: %v", err)
		} else {
			break
		}
		time.Sleep(time.Second * 5)
	}
	app.Use(cors.New())
	registerHandlers(app, db)
	app.Listen(":8080")
}
