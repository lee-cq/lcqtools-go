package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("%#v", os.Environ())
}
