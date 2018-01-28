package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"io/ioutil"

	"github.com/ananichev/simple-blockchain-service/service"
	"github.com/ananichev/simple-blockchain-service/types"

	"github.com/gorilla/mux"
)

// LastBlocks is a wrapper for adding blockchain service to handler,
// returns http handler function for retrieving last N blocks
func LastBlocks(svc *service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		amtStr := mux.Vars(r)["num"]
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

func ReceiveUpdate(svc *service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errorHandler(w, err)
			return
		}

		var update types.Update
		err = json.Unmarshal(body, &update)
		if err != nil {
			errorHandler(w, err)
			return
		}


		if !svc.MyNeighbour(update.SenderId) {
			w.Write([]byte(`{"succeess":false, err_code: 1, message: "Not my neighbour"}`))
			return
		}

		stored, err := svc.StoreUpdate(update)
		if !stored {
			w.Write([]byte(`{"succeess":false, err_code: 1, message: "Is not stored"}`))
			return
		}
		if err != nil {
			errorHandler(w, err)
			return
		}
		w.Write([]byte(`{"succeess":true, err_code: 0, message: "stored"}`))
	}
}

func errorHandler(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
