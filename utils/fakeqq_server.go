package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/LagrangeDev/LagrangeGo/message"
)

// FakeQQServer 管理FakeQQ-UI服务器
type FakeQQServer struct {
	port      int
	cmd       *exec.Cmd
	mu        sync.Mutex
	isRunning bool
	baseURL   string
}

var (
	serverInstance *FakeQQServer
	serverOnce     sync.Once
)

// GetFakeQQServer 获取FakeQQ服务器实例
func GetFakeQQServer() *FakeQQServer {
	serverOnce.Do(func() {
		serverInstance = &FakeQQServer{
			port: 3001, // 默认端口
		}
	})
	return serverInstance
}

// Start 启动FakeQQ-UI服务器
func (s *FakeQQServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return nil
	}

	// 检查pnpm是否可用
	if err := s.checkPnpm(); err != nil {
		return fmt.Errorf("pnpm检查失败: %v", err)
	}

	// 查找可用端口
	port, err := s.findAvailablePort()
	if err != nil {
		return fmt.Errorf("查找可用端口失败: %v", err)
	}
	s.port = port
	s.baseURL = fmt.Sprintf("http://localhost:%d", port)

	// 启动服务器
	cmd := exec.Command("pnpm", "run", "dev", "--", "--port", strconv.Itoa(port))
	cmd.Dir = "FakeQQ-UI"

	// 重定向输出到日志文件
	logFile, err := os.Create("fakeqq_server.log")
	if err != nil {
		return fmt.Errorf("创建日志文件失败: %v", err)
	}

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动服务器失败: %v", err)
	}

	s.cmd = cmd
	s.isRunning = true

	// 等待服务器启动
	if err := s.waitForServerReady(); err != nil {
		s.Stop()
		return fmt.Errorf("服务器启动超时: %v", err)
	}

	log.Printf("FakeQQ-UI服务器已启动在端口 %d", port)
	return nil
}

// Stop 停止服务器
func (s *FakeQQServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning || s.cmd == nil {
		return nil
	}

	if err := s.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("停止服务器失败: %v", err)
	}

	s.cmd.Wait()
	s.isRunning = false
	log.Println("FakeQQ-UI服务器已停止")
	return nil
}

// IsRunning 检查服务器是否运行
func (s *FakeQQServer) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isRunning
}

// GetBaseURL 获取服务器基础URL
func (s *FakeQQServer) GetBaseURL() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.baseURL
}

// checkPnpm 检查pnpm是否可用
func (s *FakeQQServer) checkPnpm() error {
	cmd := exec.Command("pnpm", "--version")
	cmd.Dir = "FakeQQ-UI"

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("pnpm不可用: %v", err)
	}

	log.Printf("检测到pnpm版本: %s", strings.TrimSpace(string(output)))
	return nil
}

// findAvailablePort 查找可用端口
func (s *FakeQQServer) findAvailablePort() (int, error) {
	// 从3001开始尝试
	for port := 3001; port <= 3010; port++ {
		if s.isPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("没有找到可用端口")
}

// isPortAvailable 检查端口是否可用
func (s *FakeQQServer) isPortAvailable(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 100*time.Millisecond)
	if err != nil {
		return true // 连接失败说明端口可用
	}
	conn.Close()
	return false
}

// waitForServerReady 等待服务器就绪
func (s *FakeQQServer) waitForServerReady() error {
	timeout := 30 * time.Second
	start := time.Now()

	for time.Since(start) < timeout {
		if s.checkServerHealth() {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("服务器启动超时")
}

// checkServerHealth 检查服务器健康状态
func (s *FakeQQServer) checkServerHealth() bool {
	resp, err := http.Get(fmt.Sprintf("%s/index.html", s.baseURL))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// RenderPrivateMessages 渲染私聊消息到图片
func (s *FakeQQServer) RenderPrivateMessages(privateMsg []*message.PrivateMessage, selfID uint32, width int) ([]byte, error) {
	if !s.IsRunning() {
		if err := s.Start(); err != nil {
			return nil, fmt.Errorf("启动服务器失败: %v", err)
		}
	}
	var otherID uint32
	for _, elem := range privateMsg {
		if elem.Sender.Uin != selfID {
			otherID = elem.Sender.Uin
		}
	}

	// 转换消息格式
	requestData := s.convertPrivateMessageToFakeQQFormat(privateMsg, selfID, otherID)

	// 发送渲染请求
	url := fmt.Sprintf("%s/render?width=%d", s.baseURL, width)
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("JSON编码失败: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("渲染请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("渲染失败: %s - %s", resp.Status, string(body))
	}

	// 读取图片数据
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取图片数据失败: %v", err)
	}

	return imageData, nil
}

// convertPrivateMessageToFakeQQFormat 转换私聊消息格式为FakeQQ-UI格式
func (s *FakeQQServer) convertPrivateMessageToFakeQQFormat(privateMsg []*message.PrivateMessage, selfID uint32, otherID uint32) map[string]interface{} {
	messages := make([]map[string]interface{}, 0)

	for _, elem := range privateMsg {
		messageType := "other"
		if elem.Sender.Uin == selfID {
			messageType = "self"
		}
		messageData := make([]map[string]interface{}, 0)

		for _, ele := range elem.Elements {
			switch e := ele.(type) {

			case *message.TextElement:
				messageData = append(messageData, map[string]interface{}{
					"type": "text",
					"message": []map[string]interface{}{
						{
							"type": "text",
							"data": map[string]interface{}{
								"text": e.Content,
							},
						},
					},
				})
			case *message.ImageElement:
				messageData = append(messageData, map[string]interface{}{
					"type": "image",
					"message": []map[string]interface{}{
						{
							"type": "image",
							"data": map[string]interface{}{
								"url": e.URL,
							},
						},
					},
				})
			case *message.ReplyElement:
				messageData = append(messageData, map[string]interface{}{
					"type": messageType,
					"message": []map[string]interface{}{
						{
							"type": "reply",
							"data": map[string]interface{}{
								"head":    elem.Sender.Nickname,
								"message": s.extractReplyText(e.Elements),
							},
						},
					},
				})
			default:
				// 其他类型消息转换为文本
				messageData = append(messageData, map[string]interface{}{
					"type": messageType,
					"message": []map[string]interface{}{
						{
							"type": "text",
							"data": map[string]interface{}{
								"text": fmt.Sprintf("[%d消息]", ele.Type()),
							},
						},
					},
				})
			}
		}

		messages = append(messages, map[string]interface{}{"type": messageType, "messages": messageData})
	}

	return map[string]interface{}{
		"id": map[string]string{
			"other": fmt.Sprintf("%d", otherID),
			"self":  fmt.Sprintf("%d", selfID),
		},
		"messages": messages,
	}
}

// extractReplyText 从回复元素中提取文本内容
func (s *FakeQQServer) extractReplyText(elements []message.IMessageElement) string {
	for _, elem := range elements {
		if textElem, ok := elem.(*message.TextElement); ok {
			return textElem.Content
		}
	}
	return "[回复消息]"
}
