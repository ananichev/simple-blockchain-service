package store

import (
	"encoding/json"

	"github.com/ananichev/simple-blockchain-service/types"

	"github.com/boltdb/bolt"
)

var bucket = []byte("blocks")

// BoltDb represents boltdb store
type BoltDb struct {
	conn *bolt.DB
}

// NewBoltDb establishes connection to boltdb and returns BotlDb
func NewBoltDb() (Store, error) {
	db, err := bolt.Open("bolt.db", 0644, nil)
	if err != nil {
		return nil, err
	}

	return BoltDb{
		conn: db,
	}, nil
}

// AddBlock adds new block to store
func (b BoltDb) AddBlock(block types.Block) error {
	return b.conn.Update(func(tx *bolt.Tx) error {
		buc, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return ErrAddBlock{err}
		}
		data, err := json.Marshal(block)
		if err != nil {
			return ErrAddBlock{err}
		}
		err = buc.Put(block.BlockHash, data)
		if err != nil {
			return ErrAddBlock{err}
		}

		return buc.Put([]byte("lastHash"), block.BlockHash)
	})
}

// GetLast retrieves last n blocks
func (b BoltDb) GetLast(n int) ([]types.Block, error) {
	var blocks []types.Block
	err := b.conn.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(bucket)
		if buc == nil {
			return nil
		}

		var next, last []byte
		last = buc.Get([]byte("lastHash"))

		for i := 0; i < n; i++ {
			next = buc.Get(last)
			if next == nil {
				return nil
			}
			nextBl, e := decodeBlock(next)
			if e != nil {
				return ErrGetLast{e}
			}
			blocks = append(blocks, nextBl)
			last = nextBl.PreviousBlockHash
		}
		return nil
	})
	return blocks, err
}

// Close closes boltdb connection
func (b BoltDb) Close() error {
	return b.conn.Close()
}

func decodeBlock(b []byte) (types.Block, error) {
	var block types.Block
	err := json.Unmarshal(b, &block)
	return block, err
}
