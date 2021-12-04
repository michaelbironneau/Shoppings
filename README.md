# Shoppings

Taking shopping lists into the 22nd century.

## MVP

A List is a named collection of Items. List name defaults to "{DATE} list".

An Item is a tuple of <Name, Quantity>. This will eventually become <ItemId, Quantity>, with a lookup table of Items and we create new Item as required.

- Create list (REST API)
- Archive list (REST API)
- View list (shopping mode - check off items) (Websocket API)
- Update list (Websocket API)
  - Add item
  - Update item name/quantity
  - Delete item

There should be some sort of feedback when other user is typing/on app so user knows they're potentially going to be seeing updates. Websocket updates are tuples of <ItemId, QuantityDiff, UserId> or, initially, <ListItemName, QuantityDiff, UserId>.

## Feature wishlist

- ListItems can have a sort order per store, so that the entire list can be sorted for that store (the order should be the order the user typically traverses the store in).
- Recommendations - have you forgotten X?
- Recommendations - you have X in your list, should you be adding Y?
- Recipes (collections of ListItem)
- Clone from previous list
- In Shopping Mode, screen stays on the whole time
- Notifications when lists are first created, and badges when updates are made to a list
- Import loads of Items from Tesco etc or by learning from previous lists
- Learn store order from which order items are checked off list in shopping mode
