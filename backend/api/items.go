package api

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

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
