package store

import (
	"io"

	"github.com/ananichev/simple-blockchain-service/types"
)

// Store represents storage interface
type Store interface {
	AddBlock(types.Block) error
	DeleteBlock(types.Block) error
	GetLastBlocks(int) ([]types.Block, error)
	StoreLink(link types.Link) error
	GetLinks() ([]types.Link, error)
	GetBlocks() ([]types.Block, error)

	io.Closer
}

// NewStore returns new store, for now boltdb is only supported
func NewStore() (Store, error) {
	return NewBoltDb()
}
