package clipboardm_test

import (
	"log"
	"testing"

	clipboardm "github.com/lee-cq/lcqtools-go/clipboard-m"
)

func TestReadText(t *testing.T) {
	s, err := clipboardm.ReadText()
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Cliob: %v\n", s)
}

func TestWriteAndReadText(t *testing.T) {
	testStr := "Test to Clipb"
	err := clipboardm.WriteText(testStr)
	if err != nil {
		log.Panic(err)
	}

	text, err := clipboardm.ReadText()
	if err != nil {
		log.Panic(err)
	}

	if text != testStr {
		log.Panicf("Failed: \"%s\" != \"%s\"", text, testStr)
	}
}

func TestReadImage(t *testing.T) {
	_, err := clipboardm.ReadImage()
	if err != nil {
		log.Panic(err)
	}
}
