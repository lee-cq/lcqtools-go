package tencent

import "encoding/json"

// 返回API原始内容
func RawOcrBase(imageData string, language string, verifyTLS ...bool) (Response, error) {
	// Replace with your actual Tencent AI API credentials
	payload, err := json.Marshal(map[string]string{"ImageBase64": string(imageData)})
	if err != nil {
		return Response{}, err
	}
	req := RequestOptions{
		Method:    "POST",
		Url:       "https://ocr.tencentcloudapi.com",
		Service:   "ocr",
		Action:    "GeneralBasicOCR",
		Region:    "ap-guangzhou",
		Version:   "2018-11-19",
		SecretId:  getSearetId(),
		SecretKey: getSearetKey(),
		Token:     "",
		Body:      string(payload),
	}
	if len(verifyTLS) > 0 {
		return req.Request(verifyTLS[0]), nil
	} else {
		return req.Request(false), nil
	}
}
