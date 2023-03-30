package api

import (
	"fmt"

	gt "github.com/bas24/googletranslatefree"
)

type VocabularyData struct {
	ID         string
	Word       string
	Meaning    string
	CreateDate string
	IsActive   bool
}

func TranslateText(langFrom, langTo, text string) (string, error) {
	result, err := gt.Translate(text, langFrom, langTo) //en, tr
	if err != nil {
		fmt.Println("Translation error")
		return "", fmt.Errorf("Translation Error")
	}
	fmt.Println(result)

	return result, nil
}
