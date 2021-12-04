# Shoppings

Taking shopping lists into the 22nd century.

# Database

List

- ListId
- Name

ListItem (System-Versioned)

- ListItemId
- ListId (FK)
- Name [Eventually, ItemId]
- Quantity (>= 0)
- Checked (boolean)
- UserName

Item

- ItemId
- Name

Store

- StoreId
- Name

StoreOrder

- StoreId
- ItemId
- Order (integer, default to 0, higher means put it towards top of list so earlier in journey)
