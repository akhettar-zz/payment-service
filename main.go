package main

import (
	"fmt"
	"payment-service/api"
	_ "payment-service/docs"
	"payment-service/repository"
	"log"
	"net/http"
	"os"
)

// default mongo url
var mongoUrl = "mongodb://localhost:27017/payment-db"

const (
	port  = ":8080"
	fxUrl = "http://localhost:9090/fx"
	chUrl = "http://localhost:9090/ch"
)

func init() {
	url, exists := os.LookupEnv("MONGO_URL")
	if exists {
		mongoUrl = url
	}
}

// @BasePath /
// @title Payment Service API
// @version 1.0

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {
	fmt.Println("Starting main")
	fmt.Printf("Connecting to mongo on %s", mongoUrl)
	repo := repository.NewRepository(mongoUrl)
	router := api.NewPaymentHandler(repo, fxUrl, chUrl).NewRouter()
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}
	log.Printf("Starting the paymet server on %s", srv.Addr)
	log.Fatalln(srv.ListenAndServe())
}
