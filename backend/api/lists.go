package api

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

func SetListArchived(c *fiber.Ctx, db *sql.DB, archived int) error {
	_, err := getPrincipal(db, string(c.Request().Header.Peek("X-Token")))
	if err != nil {
		return err
	}
	listId, err := c.ParamsInt("id", 0)
	if err != nil {
		return &Error{Message: "The `id` parameter should be an integer"}
	}
	_, err = db.Exec(`
		UPDATE App.List
			SET Archived = 1
		WHERE ListId = @InListId
	`, sql.Named("InListId", listId))
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
