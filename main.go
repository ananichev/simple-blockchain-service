package main

import (
	"log"
	"net/http"

	"github.com/ananichev/simple-blockchain-service/blockchain"
	"github.com/ananichev/simple-blockchain-service/handler"
	"github.com/ananichev/simple-blockchain-service/store"

	"github.com/gorilla/mux"
)

func main() {
	db, err := store.NewStore()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	svc := blockchain.NewBlockchainService(db)
	svc.Start()

	r := mux.NewRouter()
	r.HandleFunc("/add_data", handler.AddData(svc)).Methods(http.MethodPost)
	r.HandleFunc("/last_blocks/{amt}", handler.LastBlocks(svc)).Methods(http.MethodGet)

	log.Println("Listening on 3000...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
