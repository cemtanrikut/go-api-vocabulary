package api

import (
	"fmt"
	"time"

	gt "github.com/bas24/googletranslatefree"
)

type VocabularyData struct {
	ID         string
	Response   Response
	CreateDate string
	IsActive   bool
}

type Response struct {
	LangFrom string
	LangTo   string
	Text     string
	Result   string
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

func SaveTextToDB(data Response) error {
	resp := VocabularyData{
		ID:         "",
		Response:   data,
		CreateDate: time.Now().GoString(),
		IsActive:   true,
	}

	fmt.Println(resp)

	// TODO: Should save to DB from here
	// ...

	return nil
}
