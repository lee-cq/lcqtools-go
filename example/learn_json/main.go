package main

import (
	"encoding/json"
	"fmt"
)

type RespError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

type Response struct {
	Error     RespError `json:"Error"`
	RequestId string    `json:"RequestId"`

	// TextBase
	TargetText string `json:"SourceText"`
	Source     string `json:"Source"`
	Target     string `json:"Target"`

	// OCR GinBase
}

type jResponse struct {
	Response Response `json:"Response"`
}

func main() {
	// respError := `{"Response":{"Error":{"Code":"AuthFailure.SecretIdNotFound","Message":"The SecretId is not found, please ensure that your SecretId is correct."},"RequestId":"b253d23f-1ecb-460e-b794-89d1893334de"}}`
	respSuccess := `{"Response": {"RequestId": "asdad-adssad-adsasd", "TargetText": "我爱你", "Source": "en", "Target": "cn"}}`
	var j jResponse
	json.Unmarshal([]byte(respSuccess), &j)
	fmt.Printf("%+v", j)
	// fmt.Println(j.Response["Error"])

}
