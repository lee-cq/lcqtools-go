package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func sha256hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

func hmacsha256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}

type RequestOptions struct {
	Method    string
	Url       string
	Service   string
	Action    string
	Region    string
	Version   string
	SecretId  string
	SecretKey string
	Token     string
	Query     map[string]string
	Headers   map[string]string
	Body      map[string]string
}

func TencentRequest(ops RequestOptions) (string, error) {
	urlp, err := url.Parse(ops.Url)
	if err != nil {
		fmt.Printf("Failed to parse URL: %v\n", err)
		return "", err
	}
	host := urlp.Host

	algorithm := "TC3-HMAC-SHA256"
	var timestamp = time.Now().Unix()

	// ************* 步骤 1：拼接规范请求串 *************
	httpRequestMethod := "POST"
	canonicalURI := "/"
	canonicalQueryString := ""
	contentType := "application/json; charset=utf-8"
	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s\nx-tc-action:%s\n",
		contentType, host, strings.ToLower(ops.Action))
	signedHeaders := "content-type;host;x-tc-action"
	payload := "{\"ImageBase64\":\"asd\"}"
	hashedRequestPayload := sha256hex(payload)
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		httpRequestMethod,
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		hashedRequestPayload)
	log.Println(canonicalRequest)

	// ************* 步骤 2：拼接待签名字符串 *************
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, ops.Service)
	hashedCanonicalRequest := sha256hex(canonicalRequest)
	string2sign := fmt.Sprintf("%s\n%d\n%s\n%s",
		algorithm,
		timestamp,
		credentialScope,
		hashedCanonicalRequest)
	log.Println(string2sign)

	// ************* 步骤 3：计算签名 *************
	secretDate := hmacsha256(date, "TC3"+ops.SecretKey)
	secretService := hmacsha256(ops.Service, secretDate)
	secretSigning := hmacsha256("tc3_request", secretService)
	signature := hex.EncodeToString([]byte(hmacsha256(string2sign, secretSigning)))
	log.Println(signature)

	// ************* 步骤 4：拼接 Authorization *************
	authorization := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm,
		ops.SecretId,
		credentialScope,
		signedHeaders,
		signature)
	log.Println(authorization)

	// ************* 步骤 5：构造并发起请求 *************

	httpRequest, _ := http.NewRequest("POST", ops.Url, strings.NewReader(payload))
	httpRequest.Header = map[string][]string{
		"Host":           {host},
		"X-TC-Action":    {ops.Action},
		"X-TC-Version":   {ops.Version},
		"X-TC-Timestamp": {strconv.FormatInt(timestamp, 10)},
		"Content-Type":   {contentType},
		"Authorization":  {authorization},
	}
	if ops.Region != "" {
		httpRequest.Header["X-TC-Region"] = []string{ops.Region}
	}
	if ops.Token != "" {
		httpRequest.Header["X-TC-Token"] = []string{ops.Token}
	}
	httpClient := http.Client{}
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println(body.String())
	return "", nil
}

func TencentOcrBase(imageData string) (string, error) {
	// Replace with your actual Tencent AI API credentials
	req := RequestOptions{
		Method:    "POST",
		Url:       "https://ocr.tencentcloudapi.com",
		Service:   "ocr",
		Action:    "GeneralBasicOCR",
		Region:    "",
		Version:   "2018-11-19",
		SecretId:  "SecretId",
		SecretKey: "SecretKey",
		Token:     "",
		Body:      map[string]string{"ImageBase64": string(imageData)},
	}
	return TencentRequest(req)
}
