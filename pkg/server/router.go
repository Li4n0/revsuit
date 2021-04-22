package server

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/frontend"
	"github.com/li4n0/revsuit/pkg/dns"
	"github.com/li4n0/revsuit/pkg/mysql"
	"github.com/li4n0/revsuit/pkg/rhttp"
	log "unknwon.dev/clog/v2"
)

func (revsuit *Revsuit) registerRouter() {
	revsuit.http.Router = gin.Default()
	revsuit.registerPlatformRouter()
	revsuit.registerHttpRouter()
}

func (revsuit *Revsuit) registerPlatformRouter() {
	// /api need Authorization
	api := revsuit.http.Router.Group("/revsuit/api")
	api.Use(func(c *gin.Context) {
		cookieToken, err := c.Request.Cookie("token")
		if !(c.Request.Header.Get("Token") == revsuit.http.Token || err == nil && cookieToken.Value == revsuit.http.Token) {
			c.Abort()
			c.Status(403)
		}
	})
	revsuit.http.ApiGroup = api

	//platform routers
	api.GET("/events", events)
	api.GET("/ping", ping)
}

func (revsuit *Revsuit) registerHttpRouter() {
	revsuit.http.Router.NoRoute(revsuit.http.Receive)

	//register frontend
	fe, err := fs.Sub(frontend.FS, "dist")
	if err != nil {
		log.Fatal("Failed to sub path `dist`: %v", err)
	}
	revsuit.http.Router.StaticFS("/revsuit/admin", http.FS(fe))

	// init record router group
	recordGroup := revsuit.http.ApiGroup.Group("/record")

	httpGroup := recordGroup.Group("/http")
	httpGroup.GET("", rhttp.ListRecords)

	dnsGroup := recordGroup.Group("/dns")
	dnsGroup.GET("", dns.List)

	mysqlGroup := recordGroup.Group("/mysql")
	mysqlGroup.GET("", mysql.List)

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

	// init file router group
	fileGroup := revsuit.http.ApiGroup.Group("/file")
	fileGroup.GET("/mysql/:id", mysql.GetFile)

}
