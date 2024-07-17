package tencent

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
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
	Body      []byte

	Timestamp     int64
	authorization string
	algorithm     string
}

func NewRequest() *RequestOptions {
	return &RequestOptions{
		Timestamp: time.Now().Unix(),
		algorithm: "TC3-HMAC-SHA256",
	}
}

type StringMapping map[string]string

func (ops *RequestOptions) getAlgorithm() string {
	if ops.algorithm == "" {
		ops.algorithm = "TC3-HMAC-SHA256"
		return "TC3-HMAC-SHA256"
	} else {
		return ops.algorithm
	}

}

func (ops *RequestOptions) Host() (string, error) {
	urlp, err := url.Parse(ops.Url)
	if err != nil {
		return "", err
	}
	return urlp.Host, nil
}

func (ops *RequestOptions) canonicalQueries() (string, error) {
	urlp, err := url.Parse(ops.Url)
	if err != nil {
		return "", err
	}
	urlPQ, err := url.ParseQuery(urlp.RawQuery)
	if err != nil {
		return "", err
	}

	for k, v := range ops.Query {
		urlPQ.Add(k, v)
	}
	return urlPQ.Encode(), nil
}

func (ops *RequestOptions) canonicalHeadersAndSignedHeaders() (string, string, error) {

	if ops.Headers == nil {
		ops.Headers = make(map[string]string)
	}
	host, err := ops.Host()
	if err != nil {
		return "", "", err
	}
	delete(ops.Headers, "Content-Type")
	delete(ops.Headers, "Host")
	ops.Headers["content-type"] = "application/json; charset=utf-8"
	ops.Headers["host"] = host
	ops.Headers["x-tc-action"] = ops.Action

	keys := make([]string, 0, len(ops.Headers))
	lowerHeaders := make(map[string]string, len(ops.Headers))
	for k, v := range ops.Headers {
		k = strings.ToLower(k)
		v = strings.ToLower(v)
		keys = append(keys, k)
		lowerHeaders[k] = v
	}

	sort.Strings(keys)
	kvs := make([]string, 0, len(ops.Headers))

	for _, k := range keys {
		kvs = append(kvs, fmt.Sprintf("%s:%s", k, lowerHeaders[k]))
	}

	return strings.Join(kvs, "\n") + "\n", strings.Join(keys, ";"), nil
}

func (ops *RequestOptions) Signature() (string, error) {

	if ops.Timestamp == 0 {
		ops.Timestamp = time.Now().Unix()
	}

	day := time.Unix(ops.Timestamp, 0).UTC().Format("2006-01-02")

	canonicalQueries, err := ops.canonicalQueries()
	if err != nil {
		return "", err
	}

	canonicalHeaders, signedHeaders, err := ops.canonicalHeadersAndSignedHeaders()
	if err != nil {
		return "", err
	}
	fmt.Printf("%s\n\n%s", canonicalHeaders, signedHeaders)

	canonicalRequest := strings.Join([]string{
		ops.Method,
		"/",
		canonicalQueries,
		canonicalHeaders,
		signedHeaders,
		sha256hex(string(ops.Body)),
	}, "\n")

	str2sgin := strings.Join([]string{
		ops.getAlgorithm(),                                 // Algorithm
		strconv.FormatInt(ops.Timestamp, 10),               // RequestTimestamp
		fmt.Sprintf("%s/%s/tc3_request", day, ops.Service), // CredentialScope
		sha256hex(canonicalRequest),
	}, "\n")

	secretDate := hmacsha256(day, "TC3"+ops.SecretKey)
	secretService := hmacsha256(ops.Service, secretDate)
	secretSigning := hmacsha256("tc3_request", secretService)
	signature := hex.EncodeToString([]byte(hmacsha256(str2sgin, secretSigning)))

	ops.authorization = fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		ops.getAlgorithm(),
		ops.SecretId,
		fmt.Sprintf("%s/%s/tc3_request", day, ops.Service),
		signedHeaders,
		signature,
	)

	ops.Headers["Authorization"] = ops.authorization
	return ops.authorization, nil
}

func (ops *RequestOptions) Request(verifyTLS bool) (string, error) {
	if _, err := ops.Signature(); err != nil {
		return "", err
	}

	httpRequest, _ := http.NewRequest("POST", ops.Url, strings.NewReader(string(ops.Body)))
	httpRequest.Header = map[string][]string{
		// "Host":           {ops.H},
		"X-TC-Action":    {ops.Action},
		"X-TC-Version":   {ops.Version},
		"X-TC-Timestamp": {strconv.FormatInt(ops.Timestamp, 10)},
	}
	for k, v := range ops.Headers {
		httpRequest.Header[k] = []string{v}
	}
	if ops.Region != "" {
		httpRequest.Header["X-TC-Region"] = []string{ops.Region}
	}
	if ops.Token != "" {
		httpRequest.Header["X-TC-Token"] = []string{ops.Token}
	}

	httpClient := http.Client{}
	if !verifyTLS {
		httpClient = http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: verifyTLS,
				},
			},
		}
	}

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
