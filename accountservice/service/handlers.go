package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ahussein/goblog/accountservice/model"
	"github.com/gorilla/mux"
)

type DBClientManager interface {
	OpenDB()
	Seed()
	QueryAccount(accountId string) (model.Account, error)
}

var DBClient DBClientManager

func GetAccount(w http.ResponseWriter, r *http.Request) {
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
