package store

import (
	"io"

	"github.com/ananichev/simple-blockchain-service/types"
)

// Store represents storage interface
type Store interface {
	AddBlock(types.Block) error
	GetLast(int) ([]types.Block, error)

	io.Closer
}

// NewStore returns new store, for now boltdb is only supported
func NewStore() (Store, error) {
	return NewBoltDb()
}
