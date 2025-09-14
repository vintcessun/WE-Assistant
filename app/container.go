package app

import (
	"github.com/vintcessun/WE-Assistant/bot"
	"github.com/vintcessun/WE-Assistant/config"
	"github.com/vintcessun/WE-Assistant/logic"
	"github.com/vintcessun/WE-Assistant/utils"
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/auth"
)

// Container 依赖注入容器
type Container struct {
	config       *config.Config
	logger       *utils.ProtocolLogger
	client       *client.QQClient
	bot          *bot.Bot
	logicManager *logic.LogicManager
}

// NewContainer 创建新的容器实例
func NewContainer() *Container {
	return &Container{}
}

// Initialize 初始化所有依赖
func (c *Container) Initialize() error {
	c.config = &config.Config{}
	c.logger = utils.GetProtocolLogger()

	// 初始化配置
	config.Init()
	c.config = config.GlobalConfig

	// 初始化日志
	utils.Init()

	// 创建客户端
	appInfo := auth.AppList["linux"]["3.2.15-30366"]
	c.client = client.NewClient(c.config.Bot.Account, c.config.Bot.Password)
	c.client.SetLogger(c.logger)
	c.client.UseVersion(appInfo)
	c.client.AddSignServer(c.config.Bot.SignServer)
	c.client.UseDevice(auth.NewDeviceInfo(114514))

	// 创建Bot
	c.bot = bot.NewBot(c.client)
	
	// 加载签名文件
	c.bot.GetAuthManager().LoadSig()

	// 创建逻辑管理器
	c.logicManager = logic.NewLogicManager(c.client)

	return nil
}

// GetBot 获取Bot实例
func (c *Container) GetBot() *bot.Bot {
	return c.bot
}

// GetLogicManager 获取逻辑管理器实例
func (c *Container) GetLogicManager() *logic.LogicManager {
	return c.logicManager
}

// GetClient 获取客户端实例
func (c *Container) GetClient() *client.QQClient {
	return c.client
}

// GetConfig 获取配置实例
func (c *Container) GetConfig() *config.Config {
	return c.config
}
