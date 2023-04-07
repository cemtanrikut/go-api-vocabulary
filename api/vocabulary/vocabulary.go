package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	gt "github.com/bas24/googletranslatefree"
)

type VocabularyData struct {
	ID         string   `json:"_id" bson:"_id"`
	Response   Response `json:"response" bson:"response"`
	CreateDate string   `json:"create_date" bson:"create_date"`
	IsActive   bool     `json:"is_active" bson:"is_active"`
}

type Response struct {
	LangFrom string `json:"lang_from" bson:"lang_from"`
	LangTo   string `json:"lang_to" bson:"lang_to"`
	Text     string `json:"text" bson:"text"`
	Result   string `json:"result" bson:"result"`
}

func TranslateText(langFrom, langTo, text string) (Response, error) {
	result, err := gt.Translate(text, langFrom, langTo) //en, tr
	if err != nil {
		fmt.Println("Translation error")
		return Response{}, fmt.Errorf("Translation Error")
	}

	resp := Response{
		LangFrom: langFrom,
		LangTo:   langTo,
		Text:     text,
		Result:   result,
	}

	return resp, nil
}

func SaveTextToDB(resp http.ResponseWriter, req *http.Request, client *mongo.Client, collection *mongo.Collection, from string, to string, text string) string {
	translator, err := TranslateText(from, to, text)
	if err != nil {
		fmt.Println("Translate err")
	}

	//Generate UUID
	newUUID, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	var request = VocabularyData{
		ID:         string(newUUID),
		Response:   translator,
		CreateDate: time.Now().GoString(),
		IsActive:   true,
	}

	insertRes, insertErr := collection.InsertOne(context.Background(), request)
	fmt.Println("Insert response :", insertRes)
	if insertErr != nil {
		fmt.Println("Insert err")
	}

	jsonData, jsonError := json.Marshal(request)
	fmt.Println("Json data : ", jsonData)
	if jsonError != nil {
		fmt.Println("Json err")
	}

	return string(jsonData)
}
