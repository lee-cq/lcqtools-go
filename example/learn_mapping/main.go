package main

import "fmt"

type TMapping struct {
	m map[string]string
}

func main() {
	m := TMapping{}
	m.m["a"] = "a"
	fmt.Printf("%#v", m)

}
