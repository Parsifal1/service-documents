package main

import (
	"context"
	"document-service/handler"
	"document-service/models"
	"document-service/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	_, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	//services.InitCache()
	services.Session = make(map[string]struct{})
	services.Documents = make(map[uint64]models.Document)
	handlers.InitRouter(router)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
