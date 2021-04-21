package server

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/record"
	log "unknwon.dev/clog/v2"
)

func ping(c *gin.Context) {
	c.SetCookie("token", c.Request.Header["Token"][0], 0, "/revsuit/api/", c.Request.Host, true, true)
	c.String(200, "pong")
}

func events(c *gin.Context) {
	log.Info("Receive connection from ", c.Request.RemoteAddr)
	c.Stream(func(w io.Writer) bool {
		c.SSEvent("message", "connect succeed")
		select {
		case <-c.Writer.CloseNotify():
			return false
		case r := <-record.Channel():
			c.SSEvent("message", r.GetFlag())
		}
		return true
	})
	log.Info(c.Request.RemoteAddr, "disconnect")
}
