package api

import (
	"database/sql"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
)

func UpdateList(c *fiber.Ctx, db *sql.DB) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek("X-Token")))
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
	_, err := getPrincipal(db, string(c.Request().Header.Peek("X-Token")))
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
	_, err := getPrincipal(db, string(c.Request().Header.Peek("X-Token")))
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
	_, err := getPrincipal(db, string(c.Request().Header.Peek("X-Token")))
	if err != nil {
		return err
	}
	query := `
		SELECT 
	CAST(L.ListId AS VARCHAR(255)) AS 'ID',
	ISNULL(L.Name, '') As 'ListName',
	ISNULL(S.Name, '') AS 'StoreName',
	STRING_AGG(LL.Name, ', ') AS 'Summary'
FROM App.List L
INNER JOIN App.Store S ON L.StoreId = S.StoreID
OUTER APPLY (
	SELECT TOP 3 I.Name 
	FROM App.ListItem LI 
	INNER JOIN App.Item I ON I.ItemId = LI.ItemId
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
	_, err := getPrincipal(db, string(c.Request().Header.Peek("X-Token")))
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
	STRING_AGG(LL.Name, ', ') AS 'Summary'
FROM App.List L
INNER JOIN App.Store S ON L.StoreId = S.StoreID
OUTER APPLY (
	SELECT TOP 3 I.Name 
	FROM App.ListItem LI 
	INNER JOIN App.Item I ON I.ItemId = LI.ItemId
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
