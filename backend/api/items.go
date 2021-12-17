package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"strings"
)

func GetStores(c *fiber.Ctx, db *sql.DB) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek(TokenHeader)))
	if err != nil {
		return err
	}
	stores := make([]Store, 0)
	rows, err := db.Query("SELECT CAST(StoreID AS VARCHAR(255)), [Name] FROM App.Store ORDER BY StoreId")
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

func GetAllItems(c *fiber.Ctx, db *sql.DB) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek(TokenHeader)))
	if err != nil {
		return err
	}

	rows, err := db.Query(`SELECT CAST(ItemId AS VARCHAR(255)), [Name] FROM App.Item ORDER BY [Name]`)
	if err != nil {
		return dbError(err)
	}
	defer rows.Close()
	items := make([]Item, 0)
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return dbError(err)
		}
		items = append(items, item)
	}
	return c.JSON(items)
}

func AddItem(c *fiber.Ctx, db *sql.DB) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek(TokenHeader)))
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
		log.Printf("Found existing item when trying to add item")
		return &Error{Code: 409, Message: fmt.Sprintf("There is already an item with that name, with ID %v", existingID)}
	}
	row = db.QueryRow("INSERT INTO App.Item ([Name]) OUTPUT Inserted.ItemId VALUES (@InName)", sql.Named("InName", strings.Title(item.Name)))
	var newID int
	if err := row.Scan(&newID); err != nil {
		return dbError(err)
	}
	go AsyncUpdateItems(db, newID, strings.Title(item.Name)) // update list items to use the new ID if they previously had a NULL id
	return c.JSON(struct {
		ID string `json:"id"`
	}{strconv.Itoa(newID)})
}

func AsyncUpdateItems(db *sql.DB, ItemId int, ItemName string) {
	s := `UPDATE App.ListItem SET ItemId = @InItemId, [Name] = @InItemName WHERE Name LIKE @InItemName`
	_, err := db.Exec(s, sql.Named("InItemId", ItemId), sql.Named("InItemName", ItemName))
	if err != nil {
		log.Printf("Error async updating Item IDs: %v", err)
	}
}

func SearchItem(c *fiber.Ctx, db *sql.DB) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek(TokenHeader)))
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
	ret := make([]Item, 0)
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
	_, err := getPrincipal(db, string(c.Request().Header.Peek(TokenHeader)))
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
