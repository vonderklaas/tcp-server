package models

// Stock
type Stock struct {
	StockID int64  `json:"stockid"`
	Name 		string `json:"name"`
	Price 	uint64 `json:"price"`
	Company string `json:"company"`
}