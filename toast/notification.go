package toast

import "log"

func ToolsNotify(title, msg string) {
	notification := Notification{
		AppID:    "LCQ Tools",
		Title:    title,
		Message:  msg,
		Duration: Short,
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}
