package main

import (
	"log"
	"sync"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

type Func func()

func listenHotkey(fn Func, key hotkey.Key, mods ...hotkey.Modifier) (err error) {
	ms := []hotkey.Modifier{}
	ms = append(ms, mods...)
	log.Printf("mods=%v, key=%v", ms, key)
	hk := hotkey.New(ms, key)
	hk.Register()

	err = hk.Register()
	if err != nil {
		return err
	}

	for {
		<-hk.Keydown()
		<-hk.Keyup()
		fn()
	}
}

func main() { mainthread.Init(fn) }
func fn() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := listenHotkey(func() { log.Println("Ctrl+Alt+S Done. ") }, hotkey.KeyS, hotkey.ModCtrl, hotkey.ModAlt)
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer wg.Done()

		err := listenHotkey(func() { log.Println("Ctrl+Alt+O Done. ") }, hotkey.KeyO, hotkey.ModCtrl, hotkey.ModAlt)
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()

}
