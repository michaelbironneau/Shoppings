package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func GetStores(c *fiber.Ctx, db *sql.DB) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek("X-Token")))
	if err != nil {
		return err
	}
	var stores []Store
	rows, err := db.Query("SELECT CAST(StoreID AS VARCHAR(255)), [Name] FROM App.Store")
	if err != nil {
		return dbError(err)
	}
	defer rows.Close()
	for rows.Next() {
		var store Store
		if err := rows.Scan(&store.ID, &store.Name); err != nil {
			return dbError(err)
		}
		stores = append(stores, store)
	}
	return c.JSON(stores)
}

func AddItem(c *fiber.Ctx, db *sql.DB) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek("X-Token")))
	if err != nil {
		return err
	}
	var item Item
	if err := json.Unmarshal(c.Body(), &item); err != nil {
		log.Printf("Error unmarshalling body: %v", err)
		return &Error{Code: 400, Message: "Invalid request body, expecting a JSON-encoded Item"}
	}
	if item.Name == "" {
		return &Error{Code: 400, Message: "Empty item name"}
	}
	var existingID int
	row := db.QueryRow("SELECT TOP 1 1 FROM App.Item WHERE [Name] LIKE @InName", sql.Named("InName", item.Name))
	if err := row.Scan(&existingID); err == nil {
		// existing item
		return &Error{Code: 409, Message: fmt.Sprintf("There is already an item with that name, with ID %v", existingID)}
	}
	_, err = db.Exec("INSERT INTO App.Item ([Name]) VALUES (@InName)", sql.Named("InName", item.Name))
	return err
}

func SearchItem(c *fiber.Ctx, db *sql.DB) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek("X-Token")))
	if err != nil {
		return err
	}
	needle := c.Params("needle")
	needle = "%" + needle + "%"
	s := `
	SELECT TOP 10 CAST(ItemId AS VARCHAR(255)) AS [Id], [Name] FROM APP.Item WHERE [Name] LIKE @InNeedle
	ORDER BY [Name]
	`
	rows, err := db.Query(s, sql.Named("InNeedle", needle))
	if err != nil {
		return dbError(err)
	}
	defer rows.Close()
	var ret []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return dbError(err)
		}
		ret = append(ret, i)
	}
	return c.JSON(ret)
}

func CheckItem(c *fiber.Ctx, db *sql.DB) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek("X-Token")))
	if err != nil {
		return err
	}
	listID, err := c.ParamsInt("id")
	if err != nil {
		return &Error{Code: 400, Message: "List ID should be an integer"}
	}
	itemID, err := c.ParamsInt("item")
	if err != nil {
		return &Error{Code: 400, Message: "Item ID should be an integer"}
	}
	if _, err := db.Exec(`
	UPDATE App.ListItem SET Checked = 1
	WHERE ListId = @InListId AND ItemId = @InItemId
	`, sql.Named("InListId", listID), sql.Named("InItemId", itemID)); err != nil {
		return dbError(err)
	}
	return nil
}
