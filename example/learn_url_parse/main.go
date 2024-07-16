package main

import (
	"fmt"
	"net/url"
)

func main() {
	up, err := url.Parse("https://qq.com/qq?q=2&b=3#title")
	upq, _ := url.ParseQuery(up.RawQuery)
	if err != nil {
		fmt.Println("Error")
	}
	upq.Add("时间", "2024年")
	fmt.Printf("%#v", upq.Encode())
}
