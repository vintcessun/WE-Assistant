package logic

import (
	"image/png"
	"os"

	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/message"
	"github.com/vintcessun/WE-Assistant/utils"
)

func HandleGroupMessage(ctx *MessageContext) error {
	if groupMsg, ok := ctx.GetGroupMessage(); ok {
		formattedMsg := utils.FormatGroupMessageJSON(groupMsg)
		utils.Infof("收到群聊消息: %s", formattedMsg)
	}
	return nil
}

var messages = make([]*MessageContext, 0)

func HandlePrivateMessage(ctx *MessageContext) error {
	if privateMsg, ok := ctx.GetPrivateMessage(); ok {
		formattedMsg := utils.FormatPrivateMessageJSON(privateMsg)
		utils.Infof("收到私聊消息: %s", formattedMsg)
		messages = append(messages, ctx)
		privateMessages := make([]*message.PrivateMessage, 0, len(messages))
		for _, m := range messages {
			if pm, ok := m.GetPrivateMessage(); ok {
				privateMessages = append(privateMessages, pm)
			}
		}
		if img, err := utils.GenerateChatImage(privateMessages, 1024); err != nil {
			file, err := os.Create("output.png")
			if err != nil {
				return err
			}
			defer file.Close()

			err = png.Encode(file, img)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type PrivateMessageHandler struct{}

func (h *PrivateMessageHandler) Handle(client *client.QQClient, msg interface{}) error {
	if event, ok := msg.(*message.PrivateMessage); ok {
		client.SendPrivateMessage(event.Sender.Uin, []message.IMessageElement{
			message.NewText("Hello World!"),
		})
	}
	return nil
}

type GroupMessageHandler struct{}

func (h *GroupMessageHandler) Handle(client *client.QQClient, msg interface{}) error {
	if event, ok := msg.(*message.GroupMessage); ok {
		client.SendGroupMessage(event.GroupUin, []message.IMessageElement{
			message.NewText("Hello World!"),
		})
	}
	return nil
}
