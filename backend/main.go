package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/michaelbironneau/shoppings/backend/api"
	"github.com/spf13/viper"
	"log"
	"time"
)

func registerHandlers(app *fiber.App, db *sql.DB) {
	app.Post("/token", api.GetToken)
	app.Get("/lists", func(c *fiber.Ctx) error {
		return api.GetLists(c, db)
	})
	app.Post("/lists", func(c *fiber.Ctx) error {
		return nil
	})
	app.Post("/lists/:id/archive", func(c *fiber.Ctx) error {
		return api.SetListArchived(c, db, 1)
	})
	app.Post("/lists/:id/unarchive", func(c *fiber.Ctx) error {
		return api.SetListArchived(c, db, 0)
	})
	app.Get("/lists/:id", func(c *fiber.Ctx) error {
		return api.GetList(c, db)
	})
	app.Put("/lists/:id", func(c *fiber.Ctx) error {
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
