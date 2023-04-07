package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	gt "github.com/bas24/googletranslatefree"
)

type VocabularyData struct {
	ID         string   `json:"_id" bson:"_id"`
	DeviceID   string   `json:"device_id" bson:"device_id"`
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
		DeviceID:   "1",
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

func GetData(id string, resp http.ResponseWriter, req *http.Request, client *mongo.Client, collection *mongo.Collection) string {
	resp.Header().Set("Content-Type", "application/json")
	var data VocabularyData

	userData := collection.FindOne(context.Background(), bson.M{"id": id})
	err := userData.Decode(&data)

	if err != nil {
		return err.Error()
	}

	jsonResult, jsonError := json.Marshal(data)
	if jsonError != nil {
		return jsonError.Error()
	}

	return string(jsonResult)
}

func GetDatasByDeviceID(deviceID string, resp http.ResponseWriter, req *http.Request, client *mongo.Client, collection *mongo.Collection) string {
	resp.Header().Set("Content-Type", "application/json")
	var vocabMList []primitive.M

	cursor, err := collection.Find(context.Background(), bson.M{"is_active": true, "device_id": deviceID})
	if err != nil {
		return err.Error()
	}

	for cursor.Next(context.Background()) {
		var user bson.M
		if err = cursor.Decode(&user); err != nil {
			return err.Error()
		}
		vocabMList = append(vocabMList, user)
	}
	defer cursor.Close(context.Background())

	jsonResult, err := json.Marshal(vocabMList)
	if err != nil {
		return err.Error()
	}

	return string(jsonResult)

}

func DeleteFromDB(resp http.ResponseWriter, req *http.Request, client *mongo.Client, collection *mongo.Collection, id string) string {
	resp.Header().Set("Content-Type", "application/json")

	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": id, "is_active": true}, bson.D{{"$set",
		bson.D{
			{"is_active", false},
		},
	}})
	if err != nil {
		fmt.Println("not found")
		return "Error"
	}
	return "Success"
}
