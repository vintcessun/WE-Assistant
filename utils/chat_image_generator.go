package utils

import (
	"bytes"
	"fmt"
	"image"

	"github.com/LagrangeDev/LagrangeGo/message"
)

// GenerateChatImage 生成聊天记录图片（使用Chrome内核渲染）
func GenerateChatImage(messages []*message.PrivateMessage, width int) (image.Image, error) {
	// 使用FakeQQ-UI服务器进行渲染
	server := GetFakeQQServer()

	// 发送渲染请求
	imageData, err := server.RenderPrivateMessages(messages, 3673183248, width)
	if err != nil {
		return nil, fmt.Errorf("渲染失败: %v", err)
	}

	// 将字节数据转换为image.Image
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("图片解码失败: %v", err)
	}

	return img, nil
}
