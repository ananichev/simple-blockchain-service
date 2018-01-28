package store

import (
	"encoding/json"
	"strconv"

	"github.com/ananichev/simple-blockchain-service/types"

	"github.com/boltdb/bolt"
)

var blockBucket = []byte("blocks")
var linkBucket = []byte("links")

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
		buc, err := tx.CreateBucketIfNotExists(blockBucket)
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

func (b BoltDb) DeleteBlock(block types.Block) error {
	return b.conn.Update(func(tx *bolt.Tx) error {
		buc, err := tx.CreateBucketIfNotExists(blockBucket)
		if err != nil {
			return err
		}

		return buc.Delete(block.BlockHash)
	})
}

// GetLastBlocks retrieves last n blocks
func (b BoltDb) GetLastBlocks(n int) ([]types.Block, error) {
	var blocks []types.Block
	err := b.conn.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(blockBucket)
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

func (b BoltDb) StoreLink(link types.Link) error {
	return b.conn.Update(func(tx *bolt.Tx) error {
		buc, err := tx.CreateBucketIfNotExists(linkBucket)
		if err != nil {
			return ErrAddLink{err}
		}
		data, err := json.Marshal(link)
		if err != nil {
			return ErrAddLink{err}
		}
		err = buc.Put([]byte(strconv.Itoa(link.Id)), data)
		if err != nil {
			return ErrAddLink{err}
		}
		return nil
	})
}

func (b BoltDb) GetLinks() ([]types.Link, error) {
	var links []types.Link
	err := b.conn.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(linkBucket)
		if buc == nil {
			return nil
		}

		var e error
		buc.ForEach(func(k, v []byte) error {
			var link types.Link
			e = json.Unmarshal(v, &link)
			if e != nil {
				return e
			}
			links = append(links, link)
			return nil
		})
		return nil
	})
	return links, err
}

func (b BoltDb) GetBlocks() ([]types.Block, error) {
	var blocks []types.Block
	err := b.conn.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(blockBucket)
		if buc == nil {
			return nil
		}

		var e error
		buc.ForEach(func(k, v []byte) error {
			var block types.Block
			e = json.Unmarshal(v, &block)
			if e != nil {
				return e
			}
			blocks = append(blocks, block)
			return nil
		})
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
