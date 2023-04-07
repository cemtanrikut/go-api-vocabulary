package main

import (
	"fmt"
	"log"
	"net/http"

	vocabulary "github.com/cemtanrikut/go-api-vocabulary/api/vocabulary"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/translate/{from}/{to}/{text}", TranslateHandler).Methods(http.MethodGet)

	log.Println("Listening ...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln("There's an error with the server,")
	}
}

func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	paramText := vars["text"]
	paramFrom := vars["from"]
	paramTo := vars["to"]
	result, err := vocabulary.TranslateText(paramFrom, paramTo, paramText)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result.Result))
		fmt.Println("Result is: ", result)
	}
}
