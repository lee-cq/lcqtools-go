package tencent

import (
	"encoding/json"
	"fmt"
	"os"
)

type TranslateTextRequestOps struct {
	SourceText       string
	Source           string
	Target           string
	ProjectId        int
	UntranslatedText string
	TermRepoIDList   []string `json:"TermRepoIDList,omitempty"`
	SentRepoIDList   []string `json:"TermRepoIDList,omitempty"`
}

func (rOps TranslateTextRequestOps) Request(verifyTLS ...bool) (Response, error) {
	// Replace with your actual Tencent AI API credentials
	payload, err := json.Marshal(rOps)
	if err != nil {
		return Response{}, err
	}

	req := RequestOptions{
		Method:    "POST",
		Url:       "https://tmt.tencentcloudapi.com",
		Service:   "tmt",
		Action:    "TextTranslate",
		Region:    "ap-chengdu",
		Version:   "2018-03-21",
		SecretId:  os.Getenv("TC_LCQTOOLS_SECRET_ID"),
		SecretKey: os.Getenv("TC_LCQTOOLS_SECRET_KEY"),
		Token:     "",
		Body:      string(payload),
	}
	if len(verifyTLS) > 0 {
		return req.Request(verifyTLS[0]), nil
	} else {
		return req.Request(false), nil
	}
}

func TextTranslate(text, fromLang, toLang string, verify ...bool) (Response, error) {
	resp, err := TranslateTextRequestOps{
		SourceText:     text,
		Source:         fromLang,
		Target:         toLang,
		ProjectId:      0,
		TermRepoIDList: []string{},
		SentRepoIDList: []string{},
	}.Request(verify...)
	if err != nil {
		return Response{}, err
	}
	if resp.Error.Code != "" {
		return resp, fmt.Errorf("ResponseError: %s: %s", resp.Error.Code, resp.Error.Message)
	}

	return resp, nil
}
