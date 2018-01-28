package main

import (
	"log"
	"net/http"

	"github.com/ananichev/simple-blockchain-service/service"
	"github.com/ananichev/simple-blockchain-service/handler"
	"github.com/ananichev/simple-blockchain-service/store"
	"github.com/ananichev/simple-blockchain-service/state"

	"github.com/gorilla/mux"
)

func main() {
	db, err := store.NewStore()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	svc := service.NewBlockchainService(db, state.MyId, state.MyURL, state.Name)
	svc.Start()

	r := mux.NewRouter()
	r.HandleFunc("/management/add_transaction", handler.AddTransaction(svc)).Methods(http.MethodPost)
	r.HandleFunc("/management/add_link", handler.AddLink(svc)).Methods(http.MethodPost)
	r.HandleFunc("/management/sync", handler.Sync(svc)).Methods(http.MethodGet)
	r.HandleFunc("/management/status", handler.Status(svc)).Methods(http.MethodGet)

	r.HandleFunc("/blockchain/get_blocks/{num}", handler.LastBlocks(svc)).Methods(http.MethodGet)
	r.HandleFunc("/blockchain/receive_update", handler.ReceiveUpdate(svc)).Methods(http.MethodPost)


	log.Println("Listening on 3000...")
	log.Fatal(http.ListenAndServe(state.Port, r))
}
