package blockchain

import (
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/ananichev/simple-blockchain-service/store"
	"github.com/ananichev/simple-blockchain-service/types"
)

// Service represents simple blockchain service
type Service struct {
	DataCh  chan types.Row
	storage store.Store
	rows    []types.Row
}

// NewBlockchainService returns blockchain service
func NewBlockchainService(db store.Store) *Service {
	return &Service{
		DataCh:  make(chan types.Row, 10),
		storage: db,
	}
}

// Start runs listener for adding rows
func (s *Service) Start() {
	go func() {
		for {
			s.addRow(<-s.DataCh)
		}
	}()
}

// LastBlocks returns last n blocks
func (s *Service) LastBlocks(n int) ([]types.Block, error) {
	return s.storage.GetLast(n)
}

func (s *Service) addRow(r types.Row) {
	s.rows = append(s.rows, r)
	if len(s.rows) == 5 {
		s.createBlock()
		s.rows = nil
	}
}

func (s *Service) createBlock() {
	var rows []string
	for _, r := range s.rows {
		rows = append(rows, r.Data)
	}

	newBlock := types.Block{
		PreviousBlockHash: s.prevHash(),
		Rows:              rows,
		Timestamp:         time.Now().Unix(),
	}

	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", newBlock)))
	newBlock.BlockHash = h.Sum(nil)

	err := s.storage.AddBlock(newBlock)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) prevHash() []byte {
	blocks, err := s.LastBlocks(1)
	if err != nil {
		log.Fatal(err)
	}

	if len(blocks) > 0 {
		return blocks[0].BlockHash
	}
	return []byte("0")
}
