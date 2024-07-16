package main

import (
	"fmt"
	"sort"
)

func main() {
	sc := []string{"c", "ca", "时间", "Dc", "ac", "ab"}
	sort.Strings(sc)
	fmt.Printf("%v\n\n", sc)

	for _, k := range sc {
		fmt.Printf("%s\n", k)
	}
}
