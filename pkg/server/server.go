package server

import (
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/file"
	"github.com/li4n0/revsuit/internal/ipinfo"
	"github.com/li4n0/revsuit/internal/notice"
	"github.com/li4n0/revsuit/internal/record"
	"github.com/li4n0/revsuit/pkg/dns"
	"github.com/li4n0/revsuit/pkg/ftp"
	"github.com/li4n0/revsuit/pkg/mysql"
	http "github.com/li4n0/revsuit/pkg/rhttp"
	"github.com/li4n0/revsuit/pkg/rmi"
	"gorm.io/gorm/logger"
	log "unknwon.dev/clog/v2"
)

const VERSION = "0.3.0"

type Revsuit struct {
	config   *Config
	logLevel log.Level

	http  *http.Server
	dns   *dns.Server
	mysql *mysql.Server
	rmi   *rmi.Server
	ftp   *ftp.Server

	clients     map[int]*gin.Context
	clientID    int
	clientsLock sync.RWMutex
	clientsNum  chan struct{}
}

func (revsuit *Revsuit) addClient(c *gin.Context) int {
	revsuit.clientsLock.Lock()
	defer revsuit.clientsLock.Unlock()

	revsuit.clientID++
	revsuit.clients[revsuit.clientID] = c
	revsuit.clientsNum <- struct{}{}
	return revsuit.clientID
}

func (revsuit *Revsuit) removeClient(id int) {
	revsuit.clientsLock.Lock()
	defer revsuit.clientsLock.Unlock()

	delete(revsuit.clients, id)
	<-revsuit.clientsNum
}

func initDatabase(dsn string) {
	_ = log.NewConsole(100,
		log.ConsoleConfig{
			Level: log.LevelInfo,
		})

	var err error
	if strings.Contains(dsn, "@tcp") {
		err = database.InitDB(database.Mysql, dsn)
	} else if strings.Contains(dsn, ".db") {
		err = database.InitDB(database.Sqlite, dsn)
	} else {
		err = errors.New("unsupported database")
	}
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
		gin.SetMode(gin.ReleaseMode)
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
	default:
		gin.SetMode(gin.DebugMode)
		database.DB.Logger.LogMode(logger.Info)
		logLevel = log.LevelInfo
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

	initDatabase(c.Database)
	logLevel := initLog(c.LogLevel)
	ipinfo.Init(c.IpLocationDatabase)
	initNotice(c.Notice)

	s := &Revsuit{
		config:   c,
		logLevel: logLevel,
		http:     http.GetServer(),
	}

	s.dns = dns.GetServer()
	s.dns.Config = c.DNS

	s.mysql = mysql.GetServer()
	s.mysql.Config = c.MySQL

	s.rmi = rmi.GetServer()
	s.rmi.Config = c.RMI

	s.ftp = ftp.GetServer()
	s.ftp.Config = c.FTP

	if c.Addr != "" {
		s.http.SetAddr(c.Addr)
	}
	if c.Token != "" {
		s.http.SetToken(c.Token)
	}
	if c.HTTP.IpHeader != "" {
		s.http.SetIpHeader(c.HTTP.IpHeader)
	}
	if c.Domain != "" {
		s.dns.SetServerDomain(c.Domain)
	}
	if c.ExternalIP != "" {
		s.dns.SetServerIP(c.ExternalIP)
		s.ftp.SetPasvIP(c.ExternalIP)
	}
	s.clients = make(map[int]*gin.Context)
	s.clientsNum = make(chan struct{}, 100)
	return s
}

func (revsuit *Revsuit) Run() {
	defer log.Stop()
	revsuit.registerRouter()

	if revsuit.dns != nil && revsuit.dns.Enable {
		go revsuit.dns.Run()
	}
	if revsuit.rmi != nil && revsuit.rmi.Enable {
		go revsuit.rmi.Run()
	}
	if revsuit.mysql != nil && revsuit.mysql.Enable {
		go revsuit.mysql.Run()
	}
	if revsuit.ftp != nil && revsuit.ftp.Enable {
		go revsuit.ftp.Run()
	}
	go func() {
		for r := range record.Channel() {
			<-revsuit.clientsNum
			revsuit.clientsLock.RLock()
			for _, client := range revsuit.clients {
				client.SSEvent("message", r.GetFlag())
				client.Writer.Flush()
			}
			revsuit.clientsNum <- struct{}{}
			revsuit.clientsLock.RUnlock()
		}
	}()
	revsuit.http.Run()
}
