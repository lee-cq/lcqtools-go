package toast_test

import (
	"testing"

	"github.com/lee-cq/lcqtools-go/toast"
)

func TestToastPush(t *testing.T) {
	// 定义测试数据
	tests := []struct {
		name    string
		notif   toast.Notification
		wantErr bool
	}{
		{
			name: "Valid notification",
			notif: toast.Notification{
				AppID:   "Appid test",
				Title:   "Title test",
				Message: "Message Test",
				Actions: []toast.Action{
					{"protocol", "Button", ""},
					{"protocol", "Me too!", ""},
				},
			},
			wantErr: false,
		},
		// 可以添加更多的测试用例，例如测试错误的情况
	}

	// 遍历测试数据
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用 Push 方法
			if err := tt.notif.Push(); (err != nil) != tt.wantErr {
				// 使用 t.Errorf 报告错误
				t.Errorf("Notification.Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestToolsNotify(t *testing.T) {
	toast.ToolsNotify("Test Tools Notify", "Msg")
}
