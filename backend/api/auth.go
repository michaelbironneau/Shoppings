package api

import (
	"database/sql"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
)

func GetToken(ctx *fiber.Ctx) error {
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
}
