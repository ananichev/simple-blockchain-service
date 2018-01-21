package types

// Row represents data row
type Row struct {
	Data string `json:"data"`
}

// Block represents block
type Block struct {
	PreviousBlockHash []byte   `json:"previous_block_hash"`
	Rows              []string `json:"rows"`
	Timestamp         int64    `json:timestamp`
	BlockHash         []byte   `json:"block_hash"`
}
