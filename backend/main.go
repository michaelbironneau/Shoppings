package main
import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()
	app.Post("/token", func(c *fiber.Ctx) error {
		return nil
	})
	app.Get("/lists", func(c *fiber.Ctx) error {
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

	app.Listen(":3000")
}