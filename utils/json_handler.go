package utils

import (
	"encoding/json"
	"log"
)

func HandleAuthorWorkJSON(imgURLCh chan Work, data []byte, author_name string, img_quality string) error {
	var result map[string]json.RawMessage
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var body map[string]interface{}
	json.Unmarshal(result["body"], &body)
	//fmt.Println(body["illusts"])
	illusts, ok := body["illusts"].(map[string]interface{})
	if !ok {
		log.Fatal("Error: illusts is not of type map[string]interface{}")
	}
	for work_code, _ := range illusts {
		var tmp = Work{WorkID: work_code, Author: author_name, ImgFiness: img_quality}
		imgURLCh <- tmp
	}
	return nil
}

func HandleWorkTypeJSON(data []byte, img_quality string) ([]string, error) {
	var picurls []string
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Fatal(err)
		return picurls, err
	}

	body, ok := result["body"].([]interface{})
	if !ok || len(body) == 0 {
		log.Fatal("Error: parsing image urls failed - type 1 ")
	}

	for i := 0; i < len(body); i += 1 {
		urls, ok := body[i].(map[string]interface{})["urls"].(map[string]interface{})
		if !ok {
			log.Fatal("Error: parsing image urls failed - type 2 ")
		}
		selectedURL, ok := urls[img_quality].(string)
		if !ok {
			log.Fatal("Warning: cannot find image quality of: ", img_quality, " Using 'regular' instead")
		} else {
			picurls = append(picurls, selectedURL)
		}
	}

	return picurls, err
}
