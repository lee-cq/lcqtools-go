package snip_test

import (
	"fmt"
	"testing"

	"github.com/lee-cq/lcqtools-go/snip"
)

func TestSnip(t *testing.T) {
	err := snip.Snip("C:\\programGreen\\Snipaste-2.9.1-Beta-x64\\Snipaste.exe")
	if err != nil {
		fmt.Println(err)
	}
}
