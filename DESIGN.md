# Shoppings

Taking shopping lists into the 22nd century.

# Auth

Users specify a name and password which is locally stored on device and stored in hashed version in db.

There is an API endpoint to generate a token given name/password. Each token is stored in DB and validated against all requests. It's valid until revoked. This is probably fine while it's just the two of us using it; will need work at a later stage.

# Database

## "App" Schema

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

## Security Schema

Principal

- PrincipalId
- Name
- Token
- Passhash
