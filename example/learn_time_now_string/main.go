package main

import (
	"fmt"
	"time"
)

func main() {
	timer := time.Now().UnixMicro()
	fmt.Printf("%s", fmt.Sprintf("%d", timer))
}
