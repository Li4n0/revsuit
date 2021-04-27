package server

import (
	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/file"
	"github.com/li4n0/revsuit/internal/notice"
	"github.com/li4n0/revsuit/pkg/dns"
	"github.com/li4n0/revsuit/pkg/ftp"
	"github.com/li4n0/revsuit/pkg/mysql"
	http "github.com/li4n0/revsuit/pkg/rhttp"
	"github.com/li4n0/revsuit/pkg/rmi"
	"gorm.io/gorm/logger"
	log "unknwon.dev/clog/v2"
)

type Revsuit struct {
	logLevel log.Level

	http  *http.Server
	dns   *dns.Server
	mysql *mysql.Server
	rmi   *rmi.Server
	ftp   *ftp.Server
}

func initDatabase(dsn string) {
	err := database.InitDB("sqlite", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = database.DB.AutoMigrate(&http.Record{})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = database.DB.AutoMigrate(&dns.Record{})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = database.DB.AutoMigrate(&mysql.Record{})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = database.DB.AutoMigrate(&http.Rule{})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = database.DB.AutoMigrate(&dns.Rule{})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = database.DB.AutoMigrate(&mysql.Rule{})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = database.DB.AutoMigrate(&file.MySQLFile{})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = database.DB.AutoMigrate(&rmi.Record{})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = database.DB.AutoMigrate(&rmi.Rule{})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = database.DB.AutoMigrate(&ftp.Record{})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = database.DB.AutoMigrate(&ftp.Rule{})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = database.DB.AutoMigrate(&file.FTPFile{})
	if err != nil {
		log.Fatal(err.Error())
	}

}

func initLog(level string) (logLevel log.Level) {

	switch level {
	case "debug", "trace":
		gin.SetMode(gin.DebugMode)
		database.DB.Logger.LogMode(logger.Info)
		logLevel = log.LevelTrace
	case "info":
		gin.SetMode(gin.DebugMode)
		database.DB.Logger.LogMode(logger.Info)
		logLevel = log.LevelInfo
	case "warning", "warn":
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
	return logLevel
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

	logLevel := initLog(c.LogLevel)
	initDatabase(c.Database)
	initNotice(c.Notice)

	s := &Revsuit{
		logLevel: logLevel,
		http:     http.GetServer(),
	}
	if c.DNS.Enable {
		s.dns = dns.GetServer()
	}
	if c.MySQL.Enable {
		s.mysql = mysql.GetServer()
		s.mysql.Config = c.MySQL
	}
	if c.RMI.Enable {
		s.rmi = rmi.GetServer()
		s.rmi.Config = c.RMI
	}
	if c.FTP.Enable {
		s.ftp = ftp.GetServer()
		s.ftp.Config = c.FTP
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
	if revsuit.rmi != nil {
		go revsuit.rmi.Run()
	}
	if revsuit.ftp != nil {
		go revsuit.ftp.Run()
	}

	revsuit.http.Run()
}
