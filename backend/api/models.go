package api

type List struct {
	ID       string  `json:"id"`
	StoreID  *string `json:"storeId"`
	Name     string  `json:"name"`
	Summary  string  `json:"summary"`
	Archived bool    `json:"archived"`
}

type ListItem struct {
	ID         string `json:"id"`
	ListID     string `json:"listId"`
	Name       string `json:"name"`
	Quantity   string `json:"quantity"`
	Checked    bool   `json:"checked"`
	StoreOrder int    `json:"storeOrder"`
}

type ListUpdate struct {
	ItemName     string `json:"name"`
	QuantityDiff int    `json:"quantityDiff"`
	User         string `json:"user"`
}

type Store struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
