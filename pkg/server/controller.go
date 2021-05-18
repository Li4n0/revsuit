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
	"github.com/li4n0/revsuit/internal/record"
	"github.com/li4n0/revsuit/internal/recycler"
	log "unknwon.dev/clog/v2"
)

func auth(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", c.Request.Header["Token"][0], 0, "/revsuit/api/", "", false, true)
	c.String(200, "pong")
}

func ping(c *gin.Context) {
	c.String(200, "pong")
}

func events(c *gin.Context) {
	log.Info("Receive client connection from %v", c.Request.RemoteAddr)
	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Writer.CloseNotify():
			return false
		case r := <-record.Channel():
			c.SSEvent("message", r.GetFlag())
		}
		return true
	})
	log.Info("Client %s disconnect", c.Request.RemoteAddr)
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
