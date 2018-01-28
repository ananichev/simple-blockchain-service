package types

// Transaction represents transaction data
type Transaction struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int 		`json:"amount"`
}

// Block represents block
type Block struct {
	PreviousBlockHash []byte   			`json:"prev_hash"`
	Tx	              []Transaction `json:"tx"`
	Timestamp         int64    			`json:"ts"`
	BlockHash         []byte   			`json:"hash"`
}

type Link struct {
	Id int		 `json:"id"`
	URL string `json:"url"`
}

type Status struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	LastHash   []byte   `json:"last_hash"`
	Neighbours []string `json:"neighbours"`
	URL 			 string   `json:"url"`
}

type Update struct {
	SenderId string `json:"sender_id"`
	Block    Block  `json:"block"`
}
