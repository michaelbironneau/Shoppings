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
	app.Post("/token", func(c *fiber.Ctx) error {
		return api.GetToken(c, db)
	})
	app.Get("/lists", func(c *fiber.Ctx) error {
		return api.GetLists(c, db)
	})
	app.Post("/lists", func(c *fiber.Ctx) error {
		return api.InsertList(c, db)
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
		return api.UpdateList(c, db)
	})
	app.Get("/lists/:id/updates/:since", func(c *fiber.Ctx) error {
		since, err := c.ParamsInt("since")
		if err != nil {
			return &api.Error{Code: 400, Message: "`since` must be integer"}
		}
		return api.GetItems(c, db, since)
	})
	app.Patch("/lists/:id/updates", func(c *fiber.Ctx) error {
		return api.UpdateItems(c, db)
	})
	app.Get("/lists/:id/items", func(c *fiber.Ctx) error {
		return api.GetItems(c, db, 0)
	})
	app.Post("/lists/:id/items/:item/complete", func(c *fiber.Ctx) error {
		return api.CheckItem(c, db)
	})
	app.Get("/item-search/:needle", func(c *fiber.Ctx) error {
		return api.SearchItem(c, db)
	})
	app.Post("/items", func(c *fiber.Ctx) error {
		return api.AddItem(c, db)
	})
	app.Get("/stores", func(c *fiber.Ctx) error {
		return api.GetStores(c, db)
	})
}

func main() {
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*api.Error); ok {
				code = e.Code
			}
			err = ctx.Status(code).JSON(err.(*api.Error))
			return nil
		},
	})
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
