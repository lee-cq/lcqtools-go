package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	clipboardm "github.com/lee-cq/lcqtools-go/clipboard-m"
	"github.com/lee-cq/lcqtools-go/tencent"
	"github.com/lee-cq/lcqtools-go/toast"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

type Func func()

// MapStr2Key :=  map[string]uint32 {
// 	"ctrl": hotkey.ModCtrl,
// }

func listenHotkey(fnExec Func, key hotkey.Key, mods ...hotkey.Modifier) (err error) {
	ms := []hotkey.Modifier{}
	ms = append(ms, mods...)
	log.Printf("Listen Hotkey: mods=%v, key=%v", ms, key)
	hk := hotkey.New(ms, key)
	// log.Printf(" %s", hk.String())

	if err := hk.Unregister(); err != nil {
		log.Printf("Unregister Error: %s : %v", hk.String(), err)
	}

	if err := hk.Register(); err != nil {
		return err
	}

	for {
		<-hk.Keydown()
		<-hk.Keyup()
		fnExec()
	}
}

func main() {
	fmt.Println(os.Getwd())
	mainthread.Init(fn)
}

func fn() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := listenHotkey(func() {
			text, err := clipboardm.ReadText()
			if err != nil {
				log.Panicf("[Error]从剪切板获取数据失败：%s", err)
				toast.ToolsNotify("从剪切板获取数据失败：", err.Error())
				return
			}
			log.Printf("从剪切板获取数据: %s", text)

			resp, err := tencent.TextTranslate(text, "en", "zh")
			if err != nil {
				log.Printf("[Error]翻译API返回错误: %s", err)
				toast.ToolsNotify("翻译API返回错误", err.Error())
				return
			}

			err = clipboardm.WriteText(resp.TargetText)
			if err != nil && resp.TargetText != "" {
				log.Printf("[Error]剪切板写入失败: %s", err)
				toast.ToolsNotify("剪切板写入失败", err.Error())
				return
			}
			log.Printf("翻译成功: %s", resp.TargetText)
			toast.ToolsNotify("翻译成功", resp.TargetText)
			log.Println("Ctrl+Alt+S Done. ")
		}, hotkey.KeyD, hotkey.ModCtrl, hotkey.ModAlt)
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer wg.Done()

		err := listenHotkey(func() { log.Println("Ctrl+Alt+O Done. "); time.Sleep(5 * time.Second); log.Println("SLEEP 5") }, hotkey.KeyO, hotkey.ModCtrl, hotkey.ModAlt)
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()

}
