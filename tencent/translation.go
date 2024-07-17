package tencent

import "encoding/json"

type TranslateTextRequestOps struct {
	SourceText       string
	Source           string
	Target           string
	ProjectId        int
	UntranslatedText string
	TermRepoIDList   []string
	SentRepoIDList   []string
}

type TransloateResponseSuccess struct {
	Response struct {
		TargetText string `json:"SourceText"`
		Source     string `json:"Source"`
		Target     string `json:"Target"`
		RequestId  string `json:"RequestId"`
	} `json:"Response"`
}

func (rOps TranslateTextRequestOps) Request(verifyTLS ...bool) (string, error) {
	// Replace with your actual Tencent AI API credentials
	payload, err := json.Marshal(rOps)
	if err != nil {
		return "", err
	}

	req := RequestOptions{
		Method:    "POST",
		Url:       "https://ocr.tencentcloudapi.com",
		Service:   "tmt",
		Action:    "TextTranslate",
		Region:    "ap-guangzhou",
		Version:   "2018-03-21",
		SecretId:  "SecretId",
		SecretKey: "SecretKey",
		Token:     "",
		Body:      payload,
	}
	if len(verifyTLS) > 0 {
		return req.Request(verifyTLS[0])
	} else {
		return req.Request(false)
	}
}

func RawTextTranslate(text, fromLang, toLang string, verify ...bool) (string, error) {
	return TranslateTextRequestOps{
		SourceText: text,
		Source:     fromLang,
		Target:     toLang,
		ProjectId:  0,
	}.Request(verify...)
}
