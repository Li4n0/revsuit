package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/recycler"
	"github.com/li4n0/revsuit/internal/update"
	log "unknwon.dev/clog/v2"
)

func (revsuit *Revsuit) auth(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", c.Request.Header.Get("Token"), 0, revsuit.config.AdminPathPrefix, "", false, true)
	c.String(200, "pong")
}

func ping(c *gin.Context) {
	c.String(200, "pong")
}

func (revsuit *Revsuit) events(c *gin.Context) {
	id := revsuit.addClient(c)
	log.Info("Receive client[id:%d] connection from %v", id, c.Request.RemoteAddr)
	c.Stream(func(w io.Writer) bool {
		<-c.Writer.CloseNotify()
		return false
	})
	revsuit.removeClient(id)
	log.Info("Client[id:%d, remote_addr:%s] disconnect", id, c.Request.RemoteAddr)
}

func recovery(c *gin.Context) {
	timeFormat := func(t time.Time) string {
		var timeString = t.Format("2006/01/02 - 15:04:05")
		return timeString
	}
	defer func() {
		if err := recover(); err != nil {
			// Check for a broken connection, as it is not really a
			// condition that warrants a panic stack trace.
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}
			httpRequest, _ := httputil.DumpRequest(c.Request, false)
			headers := strings.Split(string(httpRequest), "\r\n")
			for idx, header := range headers {
				current := strings.Split(header, ":")
				if current[0] == "Authorization" {
					headers[idx] = current[0] + ": *"
				}
			}
			if brokenPipe {
				recycler.Recycle(fmt.Sprintf("%s\n%s", err, string(httpRequest)))
			} else if gin.IsDebugging() {
				recycler.Recycle(fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s\n", timeFormat(time.Now()), strings.Join(headers, "\r\n"), err))
			} else {
				recycler.Recycle(fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n", timeFormat(time.Now()), err))
			}

			// If the connection is dead, we can't write a status to it.
			if brokenPipe {
				_ = c.Error(err.(error)) // nolint: errcheck
				c.Abort()
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}
	}()
	c.Next()
}

func version(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "succeed",
		"error":  nil,
		"result": VERSION,
	})
}

func (revsuit *Revsuit) getUpgrade(c *gin.Context) {
	if !revsuit.config.CheckUpgrade {
		c.JSON(200, gin.H{
			"status": "succeed",
			"error":  nil,
			"result": gin.H{
				"upgradeable": false,
				"message":     "config of check upgrade is false",
			},
		})
		return
	}
	if upgradeable, release, err := update.CheckUpgrade(VERSION); err == nil && upgradeable {
		c.JSON(200, gin.H{
			"status": "succeed",
			"error":  nil,
			"result": gin.H{
				"upgradeable": true,
				"version":     release.Version,
				"release":     release.URL,
			},
		})
	} else {
		message := VERSION + " is the latest"
		if err != nil {
			message = err.Error()
		}

		c.JSON(200, gin.H{
			"status": "succeed",
			"error":  nil,
			"result": gin.H{
				"upgradeable": false,
				"message":     message,
			},
		})
	}
}
