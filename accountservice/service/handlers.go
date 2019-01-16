package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ahussein/goblog/accountservice/model"
	"github.com/gorilla/mux"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
)

type DBClientManager interface {
	OpenDB()
	Seed()
	QueryAccount(accountId string) (model.Account, error)
}

var DBClient DBClientManager

func GetAccount(w http.ResponseWriter, r *http.Request) {

	// create http client to check the matric
	client := &http.Client{
		Transport: &ochttp.Transport{
			Propagation: &b3.HTTPFormat{},
		},
	}

	client.Get("http://google.com")

	var accountId = mux.Vars(r)["accountId"]
	account, err := DBClient.QueryAccount(accountId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, _ := json.Marshal(account)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
