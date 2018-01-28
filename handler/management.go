package handler

import (
	"encoding/json"
	"net/http"
  "io/ioutil"

  "github.com/ananichev/simple-blockchain-service/service"
	"github.com/ananichev/simple-blockchain-service/types"

	// "github.com/gorilla/mux"
)

// AddTransaction is a wrapper for adding blockchain service to handler,
// returns http handler function for adding row
func AddTransaction(svc *service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errorHandler(w, err)
			return
		}

		var tx types.Transaction
		err = json.Unmarshal(body, &tx)
		if err != nil {
			errorHandler(w, err)
			return
		}
		svc.DataCh <- tx
	}
}

func AddLink(svc *service.Service) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errorHandler(w, err)
			return
		}

    var link types.Link
    err = json.Unmarshal(body, &link)
		if err != nil {
			errorHandler(w, err)
			return
		}

    svc.StoreLink(link)
  }
}

func Status(svc *service.Service) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    sts, err := svc.Status()
		if err != nil {
			errorHandler(w, err)
			return
		}

		b, err := json.Marshal(sts)
		if err != nil {
			errorHandler(w, err)
			return
		}

		w.Write(b)
  }
}

func Sync(svc *service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		blocks, err := svc.AllBlocks()

		if err != nil {
			errorHandler(w, err)
			return
		}

		b, err := json.Marshal(blocks)
		if err != nil {
			errorHandler(w, err)
			return
		}

		w.Write(b)
	}
}
