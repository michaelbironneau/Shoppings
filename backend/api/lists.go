package api

import (
	"database/sql"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"time"
)

const mergeSQL = `
MERGE App.ListItem AS TARGET
USING (VALUES (@InListId,
			   @InItemId,
			   @inName, 
			   @InQuantity,
			   @InUsername)) AS up(ListId, ItemId, [Name], Quantity, Username)
ON TARGET.ListId = up.ListId AND (up.ItemId = TARGET.ItemId OR up.[Name] = TARGET.[Name])
WHEN MATCHED THEN
	UPDATE SET TARGET.Quantity = up.Quantity, TARGET.Username = up.Username
WHEN NOT MATCHED THEN 
	INSERT ([ListId], [ItemId], [Name], Quantity, Username)
		VALUES (up.ListId, up.ItemId, up.[Name], up.Quantity, up.Username);`

func updateItem(db *sql.DB, listID int, item ListItem, principal string) error {
	// update or create item
	var (
		itemID *int
		err    error
	)
	// 1. If there's no ItemID, try to create one
	if item.ItemID == "" {
		if item.Name == "" {
			return &Error{Code: 400, Message: "Item ID and Name cannot both be blank"}
		}
		// try to convert the "custom" item over to an App.Item by looking up in App.Item
		// see if we have any matching items in App.Item
		row := db.QueryRow("SELECT TOP 1 CAST(ItemId AS VARCHAR(255)) FROM App.Item WHERE [Name] LIKE @InItemName", sql.Named("InItemName", item.Name))
		var newItemID string
		if err := row.Scan(&newItemID); err != nil && err != sql.ErrNoRows {
			return dbError(err)
		} else if err == nil {
			itemIDConv, err := strconv.Atoi(newItemID)
			if err != nil {
				// should never get reached
				log.Printf("Got non-numeric ID from DB: %s", newItemID)
				return Err500
			}
			itemID = &itemIDConv
			// update the list with the new Item ID, if it's there, to avoid inconsistencies
			if _, err := db.Exec("UPDATE App.ListItem SET ItemId = @InItemId WHERE ListId = @InListId AND [Name] LIKE @InName",
				sql.Named("InItemId", itemID),
				sql.Named("InListId", listID),
				sql.Named("InName", item.Name)); err != nil {
				return dbError(err)
			}
		} else {
			//fine if we don't have a match, ignore
		}
	} else {
		itemIDConv, err := strconv.Atoi(item.ItemID)
		if err != nil {
			return &Error{Code: 400, Message: "The item ID should be a number"}
		}
		itemID = &itemIDConv
	}
	//2. We made sure the item is a fully-fledged App.Item, if possible, and now we can just execute a merge statement
	_, err = db.Exec(mergeSQL,
		sql.Named("InListId", listID),
		sql.Named("InItemId", itemID),
		sql.Named("InName", item.Name),
		sql.Named("InQuantity", item.Quantity),
		sql.Named("InUsername", principal))
	return dbError(err)
}

func UpdateItems(c *fiber.Ctx, db *sql.DB) error {
	principal, err := getPrincipal(db, string(c.Request().Header.Peek(TokenHeader)))
	if err != nil {
		return err
	}
	listID, err := c.ParamsInt("id")
	if err != nil {
		return &Error{Code: 400, Message: "List ID should be an integer"}
	}
	var items ListUpdate
	if err := json.Unmarshal(c.Body(), &items); err != nil {
		log.Printf("Error unmarshalling update body: %v", err)
		return &Error{Code: 400, Message: "Invalid body: should have `updates` with list of updates"}
	}
	for i := range items.Updates {
		if err := updateItem(db, listID, items.Updates[i], principal); err != nil {
			return Err500
		}
	}
	return nil
}

func GetItems(c *fiber.Ctx, db *sql.DB, since int) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek(TokenHeader)))
	if err != nil {
		return err
	}
	listID, err := c.ParamsInt("id")
	if err != nil {
		return &Error{Code: 400, Message: "List ID should be an integer"}
	}
	s := `
SELECT 
	LI.ListId, 
	ISNULL(LI.ItemId, ''), 
	ISNULL(ISNULL(I.Name, LI.Name),'') AS [Name], 
	Quantity, 
	Checked, 
	ISNULL(SO.[Order], 0) AS [StoreOrder],
	ValidFrom AS [LastModified]
FROM App.ListItem LI
INNER JOIN App.List L ON L.ListId = LI.ListId
LEFT JOIN App.Item I ON I.ItemId = LI.ItemID
LEFT JOIN App.StoreOrder SO ON SO.ItemId = I.ItemId AND SO.StoreId = L.StoreId
WHERE LI.ListId = @InListId AND ValidFrom > @InSince AND LI.Quantity > 0
	`
	var ret ListUpdate
	rows, err := db.Query(s, sql.Named("InListId", listID), sql.Named("InSince", time.Unix(int64(since), 0)))
	if err != nil {
		return dbError(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			li        ListItem
			updatedAt time.Time
		)
		if err := rows.Scan(&li.ListID, &li.ItemID, &li.Name, &li.Quantity, &li.Checked, &li.StoreOrder, &updatedAt); err != nil {
			return dbError(err)
		}
		ret.Updates = append(ret.Updates, li)
		if updatedAt.Unix() > ret.UpdateTime {
			ret.UpdateTime = updatedAt.Unix()
		}
	}
	return c.JSON(ret)
}

func UpdateList(c *fiber.Ctx, db *sql.DB) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek(TokenHeader)))
	if err != nil {
		return err
	}
	var list List
	if err := json.Unmarshal(c.Body(), &list); err != nil {
		log.Printf("Error parsing request body for InsertList: %v", err)
		return &Error{Code: 400, Message: "Unexpected request body"}
	}
	listID, err := c.ParamsInt("id")
	if err != nil {
		return &Error{Code: 400, Message: "List ID should be an integer"}
	}
	s := `
	DECLARE @StoreId SMALLINT = (SELECT TOP 1 StoreId FROM App.Store WHERE [Name] LIKE @InStoreName);
	UPDATE App.List 
	SET [Name] = @InName, StoreId = @StoreId, Archived = @InArchived
	WHERE ListId = @InListId
	`
	if _, err := db.Exec(s, sql.Named("InListId", listID), sql.Named("InStoreName", list.StoreID), sql.Named("InName", list.Name), sql.Named("InArchived", list.Archived)); err != nil {
		return dbError(err)
	}
	return nil
}

func InsertList(c *fiber.Ctx, db *sql.DB) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek(TokenHeader)))
	if err != nil {
		return err
	}
	var list List
	if err := json.Unmarshal(c.Body(), &list); err != nil {
		log.Printf("Error parsing request body for InsertList: %v", err)
		return &Error{Code: 400, Message: "Unexpected request body"}
	}
	s := `
	DECLARE @StoreId SMALLINT = (SELECT TOP 1 StoreId FROM App.Store WHERE [Name] LIKE @InStoreName);
	INSERT INTO App.List  ([Name], StoreId, Archived)
	OUTPUT Inserted.ListId
		VALUES (@InName, @StoreId , @InArchived)
	`
	var retID string
	row := db.QueryRow(s, sql.Named("InStoreName", list.StoreID), sql.Named("InName", list.Name), sql.Named("InArchived", list.Archived))
	if err := row.Scan(&retID); err != nil {

		log.Printf("Error adding list: %v", err)
		return Err500
	}
	return c.JSON(struct {
		ID string `json:"id"`
	}{retID})
}

func SetListArchived(c *fiber.Ctx, db *sql.DB, archived int) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek(TokenHeader)))
	if err != nil {
		return err
	}
	listId, err := c.ParamsInt("id", 0)
	if err != nil {
		return &Error{Code: 400, Message: "The `id` parameter should be an integer"}
	}
	_, err = db.Exec(`
		UPDATE App.List
			SET Archived = @InArchived
		WHERE ListId = @InListId
	`, sql.Named("InListId", listId), sql.Named("InArchived", archived))
	return dbError(err)
}

func GetLists(c *fiber.Ctx, db *sql.DB) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek(TokenHeader)))
	if err != nil {
		return err
	}
	query := `
		SELECT 
	CAST(L.ListId AS VARCHAR(255)) AS 'ID',
	ISNULL(L.Name, '') As 'ListName',
	ISNULL(S.Name, '') AS 'StoreName',
	ISNULL(STRING_AGG(LL.Name, ', '), '') AS 'Summary'
FROM App.List L
LEFT JOIN App.Store S ON L.StoreId = S.StoreID
OUTER APPLY (
	SELECT TOP 3 ISNULL(I.Name, LI.Name) As [Name] 
	FROM App.ListItem LI 
	LEFT JOIN App.Item I ON I.ItemId = LI.ItemId
	WHERE LI.ListId = L.ListId AND LI.Checked = 0
) LL
WHERE L.Archived = 0
GROUP BY L.ListId, CAST(L.ListId AS VARCHAR(255)), ISNULL(L.Name, ''), ISNULL(S.Name, '')
ORDER BY L.ListId
		`
	var ret []List
	rows, err := db.Query(query)
	if err != nil {
		return dbError(err)
	}
	defer rows.Close()
	for rows.Next() {
		var l List
		if err := rows.Scan(&l.ID, &l.Name, &l.StoreID, &l.Summary); err != nil {
			return dbError(err)
		}
		ret = append(ret, l)
	}
	return c.JSON(ret)
}

func GetList(c *fiber.Ctx, db *sql.DB) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek(TokenHeader)))
	if err != nil {
		return err
	}
	listId, err := c.ParamsInt("id", 0)
	if err != nil {
		return &Error{Message: "The `id` parameter should be an integer"}
	}
	query := `
		SELECT 
	CAST(L.ListId AS VARCHAR(255)) AS 'ID',
	ISNULL(L.Name, '') As 'ListName',
	ISNULL(S.Name, '') AS 'StoreName',
	ISNULL(STRING_AGG(LL.Name, ', '), '') AS 'Summary'
FROM App.List L
LEFT JOIN App.Store S ON L.StoreId = S.StoreID
OUTER APPLY (
	SELECT TOP 3 ISNULL(I.Name, LI.Name) As [Name] 
	FROM App.ListItem LI 
	LEFT JOIN App.Item I ON I.ItemId = LI.ItemId
	WHERE LI.ListId = L.ListId AND LI.Checked = 0
) LL
WHERE L.ListId = @InListId
GROUP BY L.ListId, CAST(L.ListId AS VARCHAR(255)), ISNULL(L.Name, ''), ISNULL(S.Name, '')
ORDER BY L.ListId
`
	var ret List
	rows, err := db.Query(query, sql.Named("InListId", listId))
	if err != nil {
		return dbError(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ret.ID, &ret.Name, &ret.StoreID, &ret.Summary); err != nil {
			return dbError(err)
		}
		break // we only want one result, if there are somehow many (this should not be possible)
	}
	return c.JSON(ret)
}
