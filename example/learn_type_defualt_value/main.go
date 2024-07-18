package main

import "fmt"

type D struct {
	s string
	i int
	m map[string]int
	l []string
}

func main() {
	d := D{}
	if d.s == "" {
		fmt.Println("String 默认空字符串")
	}
	if d.i == 0 {
		fmt.Println("Int 默认 0")
	}
	if d.l == nil {
		fmt.Println("切片或数组默认 nil")
	}
	if d.m == nil {
		fmt.Println("mapping 默认 Nil")
	}
}
