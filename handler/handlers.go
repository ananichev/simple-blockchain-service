package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ananichev/simple-blockchain-service/blockchain"
	"github.com/ananichev/simple-blockchain-service/types"

	"github.com/gorilla/mux"
)

// AddData is a wrapper for adding blockchain service to handler,
// returns http handler function for adding row
func AddData(svc *blockchain.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errorHandler(w, err)
			return
		}

		var row types.Row
		err = json.Unmarshal(body, &row)
		if err != nil {
			errorHandler(w, err)
			return
		}
		svc.DataCh <- row
	}
}

// LastBlocks is a wrapper for adding blockchain service to handler,
// returns http handler function for retrieving last N blocks
func LastBlocks(svc *blockchain.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		amtStr := mux.Vars(r)["amt"]
		amt, err := strconv.Atoi(amtStr)
		if err != nil {
			errorHandler(w, err)
			return
		}

		blocks, err := svc.LastBlocks(amt)
		if err != nil {
			errorHandler(w, err)
			return
		}

		data, err := json.Marshal(blocks)
		if err != nil {
			errorHandler(w, err)
			return
		}
		w.Write(data)
	}
}

func errorHandler(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
