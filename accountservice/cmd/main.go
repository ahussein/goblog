package main

import (
	"fmt"

	"github.com/ahussein/goblog/accountservice/dbclient"
	"github.com/ahussein/goblog/accountservice/service"
)

var appName = "accountservice"

func main() {
	fmt.Printf("Starting %v\n", appName)
	initializeDBClient()
	service.StartWebServer("6767")
}

func initializeDBClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenDB()
	service.DBClient.Seed()
}
