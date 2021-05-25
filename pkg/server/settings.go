package server

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/pkg/dns"
	"github.com/li4n0/revsuit/pkg/ftp"
	"github.com/li4n0/revsuit/pkg/mysql"
	"github.com/li4n0/revsuit/pkg/rhttp"
	"github.com/li4n0/revsuit/pkg/rmi"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	log "unknwon.dev/clog/v2"
)

type Rules struct {
	Http  []rhttp.Rule
	Dns   []dns.Rule
	Mysql []mysql.Rule
	Rmi   []rmi.Rule
	Ftp   []ftp.Rule
}

func exportRules(c *gin.Context) {
	var (
		db    = database.DB
		rules Rules
	)

	db.Model(&rhttp.Rule{}).Find(&rules.Http)
	db.Model(&dns.Rule{}).Find(&rules.Dns)
	db.Model(&mysql.Rule{}).Find(&rules.Mysql)
	db.Model(&rmi.Rule{}).Find(&rules.Rmi)
	db.Model(&ftp.Rule{}).Find(&rules.Ftp)

	out, err := yaml.Marshal(rules)
	if err != nil {
		log.Warn("export rules error: %s", err)
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=revsuit_rules_%s.yaml", time.Now().Format("20060102150405")))
	c.String(200, string(out))
}

func (revsuit *Revsuit) importRules(c *gin.Context) {
	var (
		db    = database.DB
		rules Rules
		count int
		errs  []string
	)

	f, err := c.FormFile("rules")
	if err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		log.Trace("%v", err)
		return
	}

	file, _ := f.Open()
	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		log.Trace("%v", err)
		return
	}

	err = yaml.Unmarshal(content, &rules)
	if err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	for _, rule := range rules.Http {
		if err := db.Model(&rhttp.Rule{}).Create(&rule).Error; err != nil {
			errs = append(errs, errors.Wrap(err, fmt.Sprintf("http rule[%s]", rule.Name)).Error())
			continue
		}
		count++
	}
	if err := revsuit.http.UpdateRules(); err != nil {
		errs = append(errs, err.Error())
	}

	for _, rule := range rules.Dns {
		if err := db.Model(&dns.Rule{}).Create(&rule).Error; err != nil {
			errs = append(errs, errors.Wrap(err, fmt.Sprintf("dns rule[%s]", rule.Name)).Error())
			continue
		}
		count++
	}
	if err := revsuit.dns.UpdateRules(); err != nil {
		errs = append(errs, err.Error())
	}

	for _, rule := range rules.Mysql {
		if err := db.Model(&mysql.Rule{}).Create(&rule).Error; err != nil {
			errs = append(errs, errors.Wrap(err, fmt.Sprintf("mysql rule[%s]", rule.Name)).Error())
			continue
		}
		count++
	}
	if err := revsuit.mysql.UpdateRules(); err != nil {
		errs = append(errs, err.Error())
	}

	for _, rule := range rules.Rmi {
		if err := db.Model(&rmi.Rule{}).Create(&rule).Error; err != nil {
			errs = append(errs, errors.Wrap(err, fmt.Sprintf("rmi rule[%s]", rule.Name)).Error())
			continue
		}
		count++
	}
	if err := revsuit.rmi.UpdateRules(); err != nil {
		errs = append(errs, err.Error())
	}

	for _, rule := range rules.Ftp {
		if err := db.Model(&ftp.Rule{}).Create(&rule).Error; err != nil {
			errs = append(errs, errors.Wrap(err, fmt.Sprintf("ftp rule[%s]", rule.Name)).Error())
			continue
		}
		count++
	}
	if err := revsuit.ftp.UpdateRules(); err != nil {
		errs = append(errs, err.Error())
	}

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  errs,
		"result": fmt.Sprintf("%d rules were imported successfully, %d failed.", count, len(errs)),
	})
}

func (revsuit *Revsuit) getPlatformConfig(c *gin.Context) {
	var res = make(map[string]string)
	res["Addr"] = revsuit.config.Addr
	res["Token"] = revsuit.config.Token
	res["Domain"] = revsuit.config.Domain
	res["ExternalIP"] = revsuit.config.ExternalIP
	res["Database"] = revsuit.config.Database
	res["LogLevel"] = revsuit.config.LogLevel
	res["IpHeader"] = revsuit.config.IpHeader

	c.JSON(200, res)
}

func (revsuit *Revsuit) updatePlatformConfig(c *gin.Context) {
	var form = make(map[string]string)

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if form["LogLevel"] != revsuit.config.LogLevel {
		revsuit.logLevel = initLog(form["LogLevel"])
		revsuit.config.LogLevel = form["LogLevel"]
		if revsuit.logLevel == log.LevelTrace {
			revsuit.http.Router.Use(gin.Logger())
		}
		log.Info("Update platform config [log_level] to %s", form["LogLevel"])
	}

	if form["Token"] != revsuit.config.Token {
		revsuit.config.Token = form["Token"]
		revsuit.http.SetToken(form["Token"])
		log.Info("Update platform config [token] to %s", form["Token"])
	}

	if form["Database"] != revsuit.config.Database {
		revsuit.config.Database = form["Database"]
		initDatabase(form["Database"])
		log.Info("Update platform config [database] to %s", form["Database"])
	}

	if form["IpHeader"] != revsuit.config.IpHeader {
		revsuit.config.IpHeader = form["IpHeader"]
		revsuit.http.SetIpHeader(form["IpHeader"])
		log.Info("Update http config [ip_header] to %s", form["IpHeader"])
	}

	if form["Domain"] != revsuit.config.Domain {
		revsuit.config.Domain = form["Domain"]
		revsuit.dns.SetServerDomain(form["Domain"])
		log.Info("Update platform config [domain] to %s", form["Domain"])
		if revsuit.dns.Enable {
			revsuit.dns.Restart()
		}
	}

	if form["ExternalIP"] != revsuit.config.ExternalIP {
		revsuit.config.ExternalIP = form["ExternalIP"]
		revsuit.dns.SetServerIP(form["ExternalIP"])
		revsuit.ftp.SetPasvIP(form["ExternalIP"])
		log.Info("Update platform config [ExternalIP] to %s", form["ExternalIP"])
		if revsuit.dns.Enable {
			revsuit.dns.Restart()
		}
		if revsuit.ftp.Enable {
			revsuit.ftp.Restart()
		}
	}

	if form["Token"] != revsuit.config.Token {
		revsuit.config.Token = form["Token"]
		revsuit.http.SetToken(form["Token"])
		log.Info("Update platform config [token] to %s", form["Token"])
	}

	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"result": "update succeed",
	})
}

func (revsuit *Revsuit) getFtpConfig(c *gin.Context) {
	c.JSON(200, revsuit.ftp.Config)
}

func (revsuit *Revsuit) updateFtpConfig(c *gin.Context) {
	var form = ftp.Config{}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if form.Addr != revsuit.ftp.Addr {
		revsuit.ftp.Addr = form.Addr
		log.Info("Update ftp config [addr] to %s", form.Addr)
	}

	if form.PasvPort != revsuit.ftp.PasvPort {
		revsuit.ftp.PasvPort = form.PasvPort
		log.Info("Update ftp config [pasv_port] to %d", form.PasvPort)
	}

	if form.Enable != revsuit.ftp.Enable {
		log.Info("Update ftp config [enable] to %v", form.Enable)
		if form.Enable {
			go revsuit.ftp.Run()
		} else {
			revsuit.ftp.Stop()
		}
		return
	}

	if revsuit.ftp.Enable {
		revsuit.ftp.Restart()
	}
}

func (revsuit *Revsuit) getDnsConfig(c *gin.Context) {
	c.JSON(200, revsuit.dns.Config)
}

func (revsuit *Revsuit) updateDnsConfig(c *gin.Context) {
	var form = dns.Config{}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if form.Enable != revsuit.dns.Enable {
		log.Info("Update dns config [enable] to %v", form.Enable)
		if form.Enable {
			go revsuit.dns.Run()
		} else {
			revsuit.dns.Stop()
		}
		return
	}
}

func (revsuit *Revsuit) getMySQLConfig(c *gin.Context) {
	c.JSON(200, revsuit.mysql.Config)
}

func (revsuit *Revsuit) updateMySQLConfig(c *gin.Context) {
	var form = mysql.Config{}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if form.Addr != revsuit.mysql.Addr {
		revsuit.mysql.Addr = form.Addr
		log.Info("Update mysql config [addr] to %s", form.Addr)
	}

	if form.VersionString != revsuit.mysql.VersionString {
		revsuit.mysql.VersionString = form.VersionString
		log.Info("Update mysql config [version_string] to %s", form.VersionString)
	}

	if form.Enable != revsuit.mysql.Enable {
		log.Info("Update mysql config [enable] to %v", form.Enable)
		if form.Enable {
			go revsuit.mysql.Run()
		} else {
			revsuit.mysql.Stop()
		}
		return
	}

	if revsuit.mysql.Enable {
		revsuit.mysql.Restart()
	}
}

func (revsuit *Revsuit) getRmiConfig(c *gin.Context) {
	c.JSON(200, revsuit.rmi.Config)
}

func (revsuit *Revsuit) updateRmiConfig(c *gin.Context) {
	var form = rmi.Config{}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	if form.Addr != revsuit.rmi.Addr {
		revsuit.rmi.Addr = form.Addr
		log.Info("Update rmi config [addr] to %s", form.Addr)
	}

	if form.Enable != revsuit.rmi.Enable {
		log.Info("Update rmi config [enable] to %v", form.Enable)
		if form.Enable {
			go revsuit.rmi.Run()
		} else {
			revsuit.rmi.Stop()
		}
		return
	}

	if revsuit.rmi.Enable {
		revsuit.rmi.Restart()
	}
}

func (revsuit *Revsuit) getNoticeConfig(c *gin.Context) {
	c.JSON(200, revsuit.config.Notice)
}

func (revsuit *Revsuit) updateNoticeConfig(c *gin.Context) {
	var form = noticeConfig{}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  err.Error(),
			"data":   nil,
		})
		return
	}

	log.Info("Update notice config")
	revsuit.config.Notice = form
	initNotice(form)
}
