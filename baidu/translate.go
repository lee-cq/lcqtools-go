package baidu

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

func md5hex(data string) (string, error) {
	hasher := md5.New()
	if _, err := hasher.Write([]byte(data)); err != nil {
		return "", err
	}
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString, nil
}

func TranslateApi(text string, fromLang string, toLang string) (*Response, error) {
	apiUrl := "http://api.fanyi.baidu.com/api/trans/vip/translate"
	salt := rand.Intn(32768) + 32768 // 产生随机数
	appid := `keys["baidu_translate_appid"]`
	appkey := `keys["baidu_translate_appkey"]`
	sign, err := md5hex(appid + text + strconv.Itoa(salt) + appkey)
	if err != nil {
		return &Response{}, err
	}

	// params := map[string]string

	resp, err := http.PostForm(apiUrl,
		url.Values{
			"appid": []string{appid},
			"q":     []string{text},
			"from":  []string{fromLang},
			"to":    []string{toLang},
			"salt":  []string{strconv.Itoa(salt)},
			"sign":  []string{sign},
		},
	)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, json.Unmarshal(body, &struct{ Error string }{})
	}

	var response Response
	json.Unmarshal(body, &response)
	return &response, nil
}

type Response struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`

	From        string `json:"from"`
	To          string `json:"to"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}
