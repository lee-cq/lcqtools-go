package tencent

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func PrintRawRequest(req *http.Request, body string) {

	rp := fmt.Sprintf("> %s %s %s\n", req.Method, req.RequestURI, req.Proto)
	hl := make([]string, 0)
	for k, vs := range req.Header {
		if k == "" {
			continue
		}
		hl = append(hl, fmt.Sprintf("%s: %s", k, strings.Join(vs, ";")))
	}
	rh := ">> " + strings.Join(hl, "\n>> ") + "\n"
	rb := ">> " + strings.Join(strings.Split(body, "\n"), "\n>> ")

	// 输出原始的 HTTP 请求报文
	fmt.Println(rp + rh + "\n" + rb)
}

func getSearetId() string {
	return os.Getenv("TC_LCQTOOLS_SECRET_ID")
}

func getSearetKey() string {
	return os.Getenv("TC_LCQTOOLS_SECRET_Key")
}
