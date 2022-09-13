package server

import (
	"io/fs"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/frontend"
	"github.com/li4n0/revsuit/internal/file"
	"github.com/li4n0/revsuit/pkg/dns"
	"github.com/li4n0/revsuit/pkg/ftp"
	"github.com/li4n0/revsuit/pkg/ldap"
	"github.com/li4n0/revsuit/pkg/mysql"
	"github.com/li4n0/revsuit/pkg/rhttp"
	"github.com/li4n0/revsuit/pkg/rmi"
	log "unknwon.dev/clog/v2"
)

func (revsuit *Revsuit) registerRouter() {
	revsuit.http.Router = gin.New()
	revsuit.http.Router.Use(recovery)
	if revsuit.logLevel == log.LevelTrace {
		revsuit.http.Router.Use(gin.Logger())
	}

	revsuit.registerPlatformRouter()
	revsuit.registerHttpRouter()
}

func (revsuit *Revsuit) registerPlatformRouter() {
	// /api need Authorization
	api := revsuit.http.Router.Group(revsuit.config.AdminPathPrefix + "/api")
	api.Use(func(c *gin.Context) {
		token, _ := c.Cookie("token")
		if token != revsuit.http.Token {
			if c.Request.Header.Get("Token") != revsuit.http.Token {
				c.Abort()
				c.Status(403)
			}
		}
	})
	revsuit.http.ApiGroup = api

	//platform routers
	api.GET("/auth", revsuit.auth)
	api.GET("/events", revsuit.events)
	api.GET("/ping", ping)
	api.GET("/version", version)
	api.GET("/getUpgrade", revsuit.getUpgrade)
}

func (revsuit *Revsuit) registerHttpRouter() {
	revsuit.http.Router.NoRoute(revsuit.http.Receive)

	//register frontend
	fe, err := fs.Sub(frontend.FS, "dist")
	if err != nil {
		log.Fatal("Failed to sub path `dist`: %v", err)
	}
	revsuit.http.Router.StaticFS(path.Clean(revsuit.config.AdminPathPrefix+"/admin"), http.FS(fe))

	// init settings router group
	settingsGroup := revsuit.http.ApiGroup.Group("setting")
	settingsGroup.GET("/exportRules", exportRules)
	settingsGroup.POST("/importRules", revsuit.importRules)
	settingsGroup.GET("/getPlatformConfig", revsuit.getPlatformConfig)
	settingsGroup.POST("/updatePlatformConfig", revsuit.updatePlatformConfig)
	settingsGroup.GET("/getDnsConfig", revsuit.getDnsConfig)
	settingsGroup.POST("/updateDnsConfig", revsuit.updateDnsConfig)
	settingsGroup.GET("/getRmiConfig", revsuit.getRmiConfig)
	settingsGroup.POST("/updateRmiConfig", revsuit.updateRmiConfig)
	settingsGroup.GET("/getLdapConfig", revsuit.getLdapConfig)
	settingsGroup.POST("/updateLdapConfig", revsuit.updateLdapConfig)
	settingsGroup.GET("/getMySQLConfig", revsuit.getMySQLConfig)
	settingsGroup.POST("/updateMySQLConfig", revsuit.updateMySQLConfig)
	settingsGroup.GET("/getFtpConfig", revsuit.getFtpConfig)
	settingsGroup.POST("/updateFtpConfig", revsuit.updateFtpConfig)
	settingsGroup.GET("/getNoticeConfig", revsuit.getNoticeConfig)
	settingsGroup.POST("/updateNoticeConfig", revsuit.updateNoticeConfig)

	// init record router group
	recordGroup := revsuit.http.ApiGroup.Group("/record")

	httpGroup := recordGroup.Group("/http")
	httpGroup.GET("", rhttp.Records)
	httpGroup.DELETE("", rhttp.Records)

	dnsGroup := recordGroup.Group("/dns")
	dnsGroup.GET("", dns.Records)
	dnsGroup.DELETE("", dns.Records)

	mysqlGroup := recordGroup.Group("/mysql")
	mysqlGroup.GET("", mysql.Records)
	mysqlGroup.DELETE("", mysql.Records)

	rmiGroup := recordGroup.Group("/rmi")
	rmiGroup.GET("", rmi.Records)
	rmiGroup.DELETE("", rmi.Records)

	ldapGroup := recordGroup.Group("/ldap")
	ldapGroup.GET("", ldap.Records)
	ldapGroup.DELETE("", ldap.Records)

	ftpGroup := recordGroup.Group("/ftp")
	ftpGroup.GET("", ftp.Records)
	ftpGroup.DELETE("", ftp.Records)

	// init rule router group
	ruleGroup := revsuit.http.ApiGroup.Group("/rule")

	httpGroup = ruleGroup.Group("/http")
	httpGroup.GET("", rhttp.ListRules)
	httpGroup.POST("", rhttp.UpsertRules)
	httpGroup.DELETE("", rhttp.DeleteRules)

	dnsGroup = ruleGroup.Group("/dns")
	dnsGroup.GET("", dns.ListRules)
	dnsGroup.POST("", dns.UpsertRules)
	dnsGroup.DELETE("", dns.DeleteRules)

	mysqlGroup = ruleGroup.Group("/mysql")
	mysqlGroup.GET("", mysql.ListRules)
	mysqlGroup.POST("", mysql.UpsertRules)
	mysqlGroup.DELETE("", mysql.DeleteRules)

	rmiGroup = ruleGroup.Group("/rmi")
	rmiGroup.GET("", rmi.ListRules)
	rmiGroup.POST("", rmi.UpsertRules)
	rmiGroup.DELETE("", rmi.DeleteRules)

	ldapGroup = ruleGroup.Group("/ldap")
	ldapGroup.GET("", ldap.ListRules)
	ldapGroup.POST("", ldap.UpsertRules)
	ldapGroup.DELETE("", ldap.DeleteRules)

	ftpGroup = ruleGroup.Group("/ftp")
	ftpGroup.GET("", ftp.ListRules)
	ftpGroup.POST("", ftp.UpsertRules)
	ftpGroup.DELETE("", ftp.DeleteRules)

	// init file router group
	fileGroup := revsuit.http.ApiGroup.Group("/file")
	fileGroup.GET("/:record_type/:id", file.GetFile)

}
