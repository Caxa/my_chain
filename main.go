package main

import (
	"fmt"
	"log"
	"net/http"

	"mychain/handlers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	handlers.RegisterHandlers(router)

	fmt.Println("🚀 Server started at http://localhost:8085")
	log.Fatal(http.ListenAndServe(":8085", router))
}
