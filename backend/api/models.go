package api

type List struct {
	ID       string  `json:"id"`
	StoreID  *string `json:"storeId"`
	Name     string  `json:"name"`
	Summary  string  `json:"summary"`
	Archived bool    `json:"archived"`
}

type ListItem struct {
	ItemID     *string `json:"id"`
	ListID     string `json:"listId"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	Checked    bool   `json:"checked"`
	StoreOrder int    `json:"storeOrder"`
}

type ListUpdate struct {
	UpdateTime int64      `json:"updatedAt"`
	Updates    []ListItem `json:"updates"`
}

type Store struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Item struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	StoreOrders []StoreOrder `json:"storeOrders"`
}

type StoreOrder struct {
	StoreName string `json:"store"`
	Order     int    `json:"order"`
}
