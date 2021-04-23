package rhttp

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/qqwry"
	log "unknwon.dev/clog/v2"
)

type Server struct {
	Addr      string
	Token     string
	IpHeader  string
	Router    *gin.Engine
	ApiGroup  *gin.RouterGroup
	rules     []*Rule
	rulesLock sync.RWMutex
}

const (
	queryVar  = `\$\{query\.(.+?)\}`
	bodyVar   = `\$\{body\.(.+?)\}`
	headerVar = `\$\{header\.(.+?)\}`
)

var (
	server           *Server
	once             sync.Once
	queryVarMatcher  = regexp.MustCompile(queryVar)
	bodyVarMatcher   = regexp.MustCompile(bodyVar)
	headerVarMatcher = regexp.MustCompile(headerVar)
)

func GetServer() *Server {
	once.Do(func() {
		letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		randString := func(n int) string {
			b := make([]byte, n)
			for i := range b {
				b[i] = letterBytes[rand.Intn(len(letterBytes))]
			}
			return string(b)
		}
		server = &Server{
			Addr:  ":80",
			Token: randString(9),
			rules: make([]*Rule, 0),
		}
	})
	return server
}

func (s *Server) SetAddr(addr string) *Server {
	s.Addr = addr
	return s
}

func (s *Server) SetToken(token string) *Server {
	s.Token = token
	return s
}

func (s *Server) SetIpHeader(header string) *Server {
	s.IpHeader = header
	return s
}

func (s *Server) getRules() []*Rule {
	defer s.rulesLock.RUnlock()
	s.rulesLock.RLock()
	return s.rules
}

func (s *Server) updateRules() error {
	db := database.DB.Model(new(Rule))
	s.rulesLock.Lock()
	db.Order("rank desc").Find(&s.rules)
	s.rulesLock.Unlock()
	return nil
}

func (s *Server) Run() {
	if err := s.updateRules(); err != nil {
		log.Error(err.Error())
	}
	log.Info("Starting HTTP Server at %s, token:%s", s.Addr, s.Token)
	err := s.Router.Run(s.Addr)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getRawRequest(r *http.Request) ([]byte, error) {
	// to resolve httputil.DumpRequestOut error with "" protocol
	r.URL.Scheme = "http"
	r.URL.Host = "revsuit"

	return httputil.DumpRequestOut(r, true)
}

func compileTpl(c *gin.Context, tpl string) (compiled string) {
	compiled = tpl
	if queryVarMatcher.FindString(tpl) != "" {
		compiled = queryVarMatcher.ReplaceAllString(compiled, c.Query(queryVarMatcher.FindStringSubmatch(tpl)[1]))
	}

	if bodyVarMatcher.FindString(tpl) != "" {
		compiled = bodyVarMatcher.ReplaceAllString(compiled, c.PostForm(bodyVarMatcher.FindStringSubmatch(tpl)[1]))
	}

	if headerVarMatcher.FindString(tpl) != "" {
		compiled = headerVarMatcher.ReplaceAllString(compiled, c.GetHeader(headerVarMatcher.FindStringSubmatch(tpl)[1]))
	}
	return compiled
}

func (s *Server) Receive(c *gin.Context) {
	u := c.Request.URL.String()
	for _, _rule := range s.getRules() {
		flag, flagGroup := _rule.Match(u)
		if flag == "" {
			continue
		}

		var (
			ip   = strings.Split(c.Request.RemoteAddr, ":")[0]
			area = qqwry.Area(ip)
		)

		if ip1 := c.Request.Header.Get(s.IpHeader); s.IpHeader != "" && ip1 != "" {
			ip = ip1
			delete(c.Request.Header, s.IpHeader)
		}

		raw, err := getRawRequest(c.Request)
		if err != nil {
			log.Warn(err.Error())
		}

		// create new record
		r, err := NewRecord(_rule, flag, c.Request.Method, u, ip, area, string(raw))
		if err != nil {
			log.Error("HTTP record(rule_id:%d) created failed :%s", _rule.ID, err.Error())
			code, err := strconv.Atoi(compileTpl(c, _rule.ResponseStatusCode))
			if err != nil || code < 100 || code > 600 {
				code = 400
			}

			c.String(code, compileTpl(c, _rule.ResponseBody))
			return
		}
		log.Info("HTTP record(id:%d,rule:%s,remote_ip:%s) has been created", r.ID, _rule.Name, ip)

		//only send to client when this connection recorded first time.
		if _rule.PushToClient {
			if flagGroup != "" {
				var count int64
				database.DB.Where("rule_name=? and raw like ?", _rule.Name, "%"+flagGroup+"%").Model(&Record{}).Count(&count)
				if count <= 1 {
					r.PushToClient()
					log.Trace("HTTP record(id:%d) has been put to client message queue", r.ID)
				}
			}
			r.PushToClient()
			log.Trace("HTTP record(id:%d) has been put to client message queue", r.ID)
		}

		//send notice
		if _rule.Notice {
			go func() {
				r.Notice()
				log.Trace("HTTP record(id:%d) notice has been sent", r.ID)
			}()
		}

		for header, value := range _rule.ResponseHeaders {
			c.Header(compileTpl(c, header), compileTpl(c, value))
		}

		code, err := strconv.Atoi(compileTpl(c, _rule.ResponseStatusCode))
		if err != nil || code < 100 || code > 600 {
			code = 400
		}

		c.String(code, compileTpl(c, _rule.ResponseBody))
		return
	}

}
