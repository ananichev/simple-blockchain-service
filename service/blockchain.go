package service

import (
	"crypto/sha256"
	"encoding/json"
	"strconv"
	// "fmt"
	"log"
	"time"
	// "io"
	"bytes"
	"net/http"

	"github.com/ananichev/simple-blockchain-service/store"
	"github.com/ananichev/simple-blockchain-service/types"
)

// Service represents simple blockchain service
type Service struct {
	DataCh  chan types.Transaction
	storage store.Store
	txs    []types.Transaction
	myId   string
	myName string
	myURL  string
}

// NewBlockchainService returns blockchain service
func NewBlockchainService(db store.Store, myId, myURL, myName string) *Service {
	return &Service{
		DataCh:  make(chan types.Transaction, 10),
		storage: db,
		myId: myId,
		myName: myName,
		myURL: myURL,
	}
}

// Start runs listener for adding txs
func (s *Service) Start() {
	go func() {
		for {
			s.addRow(<-s.DataCh)
		}
	}()
}

// LastBlocks returns last n blocks
func (s *Service) LastBlocks(n int) ([]types.Block, error) {
	return s.storage.GetLastBlocks(n)
}

func (s *Service) AllBlocks() ([]types.Block, error) {
	return s.storage.GetBlocks()
}

func (s *Service) StoreUpdate(upd types.Update) (bool, error) {
	blocks , err := s.LastBlocks(1)
	if err != nil {
		return false, err
	}

	if string(blocks[0].BlockHash) == string(upd.Block.PreviousBlockHash) {
		err = s.addBlock(upd.Block)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	if string(blocks[0].PreviousBlockHash) == string(upd.Block.PreviousBlockHash) {
		if blocks[0].Timestamp < upd.Block.Timestamp {
			return false, nil
		} else {
			err := s.storage.DeleteBlock(blocks[0])
			if err != nil {
				return false, err
			}
			err = s.addBlock(upd.Block)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}

	return false, nil
}

func (s *Service) MyNeighbour(senderId string) bool {
	links, err := s.storage.GetLinks()
	if err != nil {
		return false
	}

	for _, l := range links {
		if strconv.Itoa(l.Id) == senderId {
			return true
		}
	}
	return false
}

func (s *Service) addRow(r types.Transaction) {
	s.txs = append(s.txs, r)
	if len(s.txs) == 5 {
		s.createBlock()
		s.txs = nil
	}
}

func (s *Service) createBlock() {
	newBlock := types.Block{
		PreviousBlockHash: s.prevHash(),
		Tx:            	   s.txs,
		Timestamp:         time.Now().Unix(),
	}

	// h := sha256.New()
	// str := []byte(hex.EncodeToString([]byte(fmt.Sprintf("%v", newBlock))))
	// h := sha256.Sum256([]byte(str))
	// sum := sha256.Sum256([]byte())

	// io.WriteString(h, fmt.Sprintf("%v", newBlock))


	str, err := json.Marshal(newBlock)
	if err != nil {
		log.Fatal(err)
	}

	h := sha256.New()
	h.Write(str)
	log.Println(string(str))
	hh := h.Sum(nil)
	log.Println(string(hh))
	// sum := sha256.Sum256(str)

	// str := []byte(hex.EncodeToString(sum[:]))

	newBlock.BlockHash = hh

	err = s.addBlock(newBlock)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service)addBlock(block types.Block) error {
	err := s.storage.AddBlock(block)
	if err != nil {
		return err
	}
	links, err := s.storage.GetLinks()
	if err != nil {
		return err
	}

	for _, l := range links {
		u := types.Update{
			SenderId: s.myId,
			Block: block,
		}


		b, e := json.Marshal(u)
		if e != nil {
			log.Printf("Error marshaling update: id - %v", l.Id)
			continue
		}
		_, e = http.Post(l.URL + "/blockchain/receive_update", "application/json", bytes.NewBuffer(b))
		if e != nil {
			log.Printf("Error sending update: id - %v", l.Id)
			continue
		}
	}
	return nil
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
