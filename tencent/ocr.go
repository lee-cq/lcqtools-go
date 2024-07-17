package tencent

import "encoding/json"

// 返回API原始内容
func RawOcrBase(imageData string, language string, verifyTLS ...bool) (string, error) {
	// Replace with your actual Tencent AI API credentials
	payload, err := json.Marshal(map[string]string{"ImageBase64": string(imageData)})
	if err != nil {
		return "", nil
	}
	req := RequestOptions{
		Method:    "POST",
		Url:       "https://ocr.tencentcloudapi.com",
		Service:   "ocr",
		Action:    "GeneralBasicOCR",
		Region:    "ap-guangzhou",
		Version:   "2018-11-19",
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
