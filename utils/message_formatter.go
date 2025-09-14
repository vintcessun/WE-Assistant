package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/LagrangeDev/LagrangeGo/message"
)

// MessageElementInfo 消息元素信息结构
type MessageElementInfo struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content,omitempty"`
	URL     string      `json:"url,omitempty"`
	Target  uint32      `json:"target,omitempty"`
	ReplyID uint32      `json:"reply_id,omitempty"`
}

// MessageInfo 完整消息信息结构
type MessageInfo struct {
	Type     string               `json:"type"`
	Sender   string               `json:"sender"`
	SenderID uint32               `json:"sender_id"`
	GroupID  uint32               `json:"group_id,omitempty"`
	Elements []MessageElementInfo `json:"elements"`
	Summary  string               `json:"summary"`
}

// FormatMessageElementsJSON 格式化消息元素为JSON格式
func FormatMessageElementsJSON(elements []message.IMessageElement) string {
	if len(elements) == 0 {
		return `{"elements": [], "summary": "空消息"}`
	}

	var elementInfos []MessageElementInfo
	var summaryParts []string

	for _, element := range elements {
		info := formatElementToInfo(element)
		if info != nil {
			elementInfos = append(elementInfos, *info)

			// 构建摘要（完整输出，不省略）
			switch info.Type {
			case "text":
				summaryParts = append(summaryParts, info.Content.(string))
			case "image":
				summaryParts = append(summaryParts, "[图片]")
			case "at":
				summaryParts = append(summaryParts, fmt.Sprintf("@%d", info.Target))
			case "reply":
				summaryParts = append(summaryParts, fmt.Sprintf("[回复消息ID:%d]", info.ReplyID))
			default:
				summaryParts = append(summaryParts, fmt.Sprintf("[%s]", info.Type))
			}
		}
	}

	summary := strings.Join(summaryParts, " ")
	if summary == "" {
		summary = "多媒体消息"
	}

	result := map[string]interface{}{
		"elements": elementInfos,
		"summary":  summary,
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf(`{"error": "JSON序列化失败: %v"}`, err)
	}

	return string(jsonBytes)
}

// formatElementToInfo 格式化单个消息元素为信息结构
func formatElementToInfo(element message.IMessageElement) *MessageElementInfo {
	switch elem := element.(type) {
	case *message.TextElement:
		return &MessageElementInfo{
			Type:    "text",
			Content: elem.Content,
		}
	case *message.AtElement:
		return &MessageElementInfo{
			Type:   "at",
			Target: elem.TargetUin,
		}
	case *message.ImageElement:
		url := ""
		if elem.URL != "" {
			url = elem.URL
		}
		return &MessageElementInfo{
			Type: "image",
			URL:  url,
		}
	case *message.ReplyElement:
		return &MessageElementInfo{
			Type:    "reply",
			ReplyID: elem.ReplySeq,
		}
	default:
		// 只处理指定的类型，其他类型返回nil
		return nil
	}
}

// FormatPrivateMessageJSON 格式化私聊消息为JSON
func FormatPrivateMessageJSON(privateMsg *message.PrivateMessage) string {
	info := MessageInfo{
		Type:     "private",
		Sender:   privateMsg.Sender.Nickname,
		SenderID: privateMsg.Sender.Uin,
		Elements: []MessageElementInfo{},
	}

	// 格式化消息元素
	var elementInfos []MessageElementInfo
	var summaryParts []string

	for _, element := range privateMsg.Elements {
		elemInfo := formatElementToInfo(element)
		if elemInfo != nil {
			elementInfos = append(elementInfos, *elemInfo)

			switch elemInfo.Type {
			case "text":
				summaryParts = append(summaryParts, elemInfo.Content.(string))
			case "image":
				summaryParts = append(summaryParts, "[图片]")
			case "at":
				summaryParts = append(summaryParts, fmt.Sprintf("@%d", elemInfo.Target))
			case "reply":
				summaryParts = append(summaryParts, fmt.Sprintf("[回复消息ID:%d]", elemInfo.ReplyID))
			}
		}
	}

	info.Elements = elementInfos
	info.Summary = strings.Join(summaryParts, " ")
	if info.Summary == "" {
		info.Summary = "多媒体消息"
	}

	jsonBytes, err := json.Marshal(info)
	if err != nil {
		return fmt.Sprintf(`{"error": "JSON序列化失败: %v"}`, err)
	}

	return string(jsonBytes)
}

// FormatGroupMessageJSON 格式化群聊消息为JSON
func FormatGroupMessageJSON(groupMsg *message.GroupMessage) string {
	info := MessageInfo{
		Type:     "group",
		Sender:   groupMsg.Sender.Nickname,
		SenderID: groupMsg.Sender.Uin,
		GroupID:  groupMsg.GroupUin,
		Elements: []MessageElementInfo{},
	}

	// 格式化消息元素
	var elementInfos []MessageElementInfo
	var summaryParts []string

	for _, element := range groupMsg.Elements {
		elemInfo := formatElementToInfo(element)
		if elemInfo != nil {
			elementInfos = append(elementInfos, *elemInfo)

			switch elemInfo.Type {
			case "text":
				summaryParts = append(summaryParts, elemInfo.Content.(string))
			case "image":
				summaryParts = append(summaryParts, "[图片]")
			case "at":
				summaryParts = append(summaryParts, fmt.Sprintf("@%d", elemInfo.Target))
			case "reply":
				summaryParts = append(summaryParts, fmt.Sprintf("[回复消息ID:%d]", elemInfo.ReplyID))
			}
		}
	}

	info.Elements = elementInfos
	info.Summary = strings.Join(summaryParts, " ")
	if info.Summary == "" {
		info.Summary = "多媒体消息"
	}

	jsonBytes, err := json.Marshal(info)
	if err != nil {
		return fmt.Sprintf(`{"error": "JSON序列化失败: %v"}`, err)
	}

	return string(jsonBytes)
}

// GetMessageSummary 获取消息的简要摘要（完整输出）
func GetMessageSummary(elements []message.IMessageElement) string {
	var summaryParts []string

	for _, element := range elements {
		switch elem := element.(type) {
		case *message.TextElement:
			summaryParts = append(summaryParts, elem.Content)
		case *message.AtElement:
			summaryParts = append(summaryParts, fmt.Sprintf("@%d", elem.TargetUin))
		case *message.ImageElement:
			summaryParts = append(summaryParts, "[图片]")
		case *message.ReplyElement:
			summaryParts = append(summaryParts, fmt.Sprintf("[回复消息ID:%d]", elem.ReplySeq))
		default:
			// 跳过其他类型
		}
	}

	if len(summaryParts) == 0 {
		return "[多媒体消息]"
	}

	return strings.Join(summaryParts, " ")
}
