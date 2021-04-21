package server

import (
	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/notice"
	"github.com/li4n0/revsuit/pkg/dns"
	"github.com/li4n0/revsuit/pkg/mysql"
	http "github.com/li4n0/revsuit/pkg/rhttp"
	"gorm.io/gorm/logger"
	log "unknwon.dev/clog/v2"
)

type Revsuit struct {
	http  *http.Server
	dns   *dns.Server
	mysql *mysql.Server
}

func initDatabase(dsn string) {
	err := database.InitDB("sqlite", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = database.DB.AutoMigrate(&http.Record{})
	err = database.DB.AutoMigrate(&dns.Record{})
	err = database.DB.AutoMigrate(&mysql.Record{})
	err = database.DB.AutoMigrate(&http.Rule{})
	err = database.DB.AutoMigrate(&dns.Rule{})
	err = database.DB.AutoMigrate(&mysql.Rule{})
	err = database.DB.AutoMigrate(&mysql.File{})
	if err != nil {
		log.Fatal(err.Error())
	}

}

func initLog(level string) {
	var logLevel log.Level

	switch level {
	case "debug":
		gin.SetMode(gin.DebugMode)
		database.DB.Logger.LogMode(logger.Info)
		logLevel = log.LevelTrace
	case "info":
		gin.SetMode(gin.DebugMode)
		database.DB.Logger.LogMode(logger.Info)
		logLevel = log.LevelInfo
	case "warning":
		gin.SetMode(gin.ReleaseMode)
		database.DB.Logger.LogMode(logger.Warn)
		logLevel = log.LevelWarn
	case "error":
		gin.SetMode(gin.ReleaseMode)
		database.DB.Logger.LogMode(logger.Error)
		logLevel = log.LevelError
	case "fatal":
		gin.SetMode(gin.ReleaseMode)
		database.DB.Logger.LogMode(logger.Error)
		logLevel = log.LevelFatal
	}
	_ = log.NewConsole(100,
		log.ConsoleConfig{
			Level: logLevel,
		})

}

func initNotice(nc noticeConfig) {
	n := notice.New()
	if nc.DingTalk != "" {
		n.AddBot(&notice.DingTalk{
			URL: nc.DingTalk,
		})
	}
	if nc.Lark != "" {
		n.AddBot(&notice.Lark{
			URL: nc.Lark,
		})
	}
	if nc.WeiXin != "" {
		n.AddBot(&notice.Weixin{
			URL: nc.WeiXin,
		})
	}
	if nc.Slack != "" {
		n.AddBot(&notice.Slack{
			URL: nc.Slack,
		})
	}
}

func New(c *Config) *Revsuit {

	initDatabase(c.Database)
	initLog(c.LogLevel)
	initNotice(c.Notice)

	s := &Revsuit{
		http: http.GetServer(),
	}
	if c.DNS.Enable {
		s.dns = dns.GetServer()
	}
	if c.Mysql.Enable {
		s.mysql = mysql.GetServer()
		s.mysql.Config = c.Mysql
	}

	if c.Addr != "" {
		s.http.SetAddr(c.Addr)
	}
	if c.Token != "" {
		s.http.SetToken(c.Token)
	}
	if c.IpHeader != "" {
		s.http.SetIpHeader(c.IpHeader)
	}
	return s
}

func (revsuit *Revsuit) Run() {
	defer log.Stop()
	revsuit.registerRouter()

	if revsuit.dns != nil {
		go revsuit.dns.Run()
	}
	if revsuit.mysql != nil {
		go revsuit.mysql.Run()
	}

	revsuit.http.Run()
}
