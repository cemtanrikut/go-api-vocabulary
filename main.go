package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	vocabulary "github.com/cemtanrikut/go-api-vocabulary/api/vocabulary"
	db "github.com/cemtanrikut/go-api-vocabulary/db"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var ctx context.Context
var collection *mongo.Collection
var router *mux.Router

func main() {
	// client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://admin:CEMt1994@sandbox.0sac2.mongodb.net/?retryWrites=true&w=majority"))

	// router = mux.NewRouter()
	// vocabCollection := client.Database("vocup").Collection("vocabulary")
	router, ctx, client, collection = db.MongoClient()
	//router, ctx, client, vocabCollection = db.MongoClient("vocabulary")
	fmt.Println("***", &router, ctx, &client, &collection)
	router.HandleFunc("/translate/{from}/{to}/{text}", TranslateHandler).Methods(http.MethodGet)
	router.HandleFunc("/translateDB/{from}/{to}/{text}", TranslateDBHandler).Methods(http.MethodGet)

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

func TranslateDBHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paramText := vars["text"]
	paramFrom := vars["from"]
	paramTo := vars["to"]
	result := vocabulary.SaveTextToDB(w, r, client, collection, paramFrom, paramTo, paramText)

	w.WriteHeader(http.StatusOK)
	fmt.Println("Result is: ", result)

}
