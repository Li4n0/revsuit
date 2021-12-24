package rhttp

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/ipinfo"
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

	httpServer http.Server
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

func (s *Server) UpdateRules() error {
	db := database.DB.Model(new(Rule))
	defer s.rulesLock.Unlock()
	s.rulesLock.Lock()
	return db.Order("base_rank desc").Find(&s.rules).Error
}

func (s *Server) startHttpServer() {
	log.Info("Starting HTTP Server at %s, token:%s", s.Addr, s.Token)
	s.httpServer = http.Server{
		Addr:    s.Addr,
		Handler: s.Router,
	}
	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err.Error())
	}
}

func (s *Server) stopHttpServer() {
	log.Info("HTTP Server is stopping...")
	err := s.httpServer.Shutdown(context.TODO())
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (s *Server) Restart() {
	//only need to stop http server, because it will start in a loop
	s.stopHttpServer()
}

func (s *Server) Run() {
	if err := s.UpdateRules(); err != nil {
		log.Warn(err.Error())
	}
	for {
		s.startHttpServer()
		time.Sleep(2 * time.Second)
	}
}

func getRawRequest(r *http.Request) ([]byte, error) {
	// to resolve httputil.DumpRequestOut error with "" protocol
	r.URL.Scheme = "http"
	r.URL.Host = "revsuit"

	return httputil.DumpRequestOut(r, true)
}

func compileTpl(c *gin.Context, tpl string, vars map[string]string) (compiled string) {
	compiled = tpl
	for _, submatch := range queryVarMatcher.FindAllStringSubmatch(tpl, -1) {
		compiled = strings.ReplaceAll(compiled, submatch[0], c.Query(submatch[1]))
	}

	for _, submatch := range bodyVarMatcher.FindAllStringSubmatch(tpl, -1) {
		compiled = strings.ReplaceAll(compiled, submatch[0], c.PostForm(submatch[1]))
	}

	for _, submatch := range headerVarMatcher.FindAllStringSubmatch(tpl, -1) {
		compiled = strings.ReplaceAll(compiled, submatch[0], c.GetHeader(submatch[1]))
	}

	for n, v := range vars {
		compiled = strings.ReplaceAll(compiled, "${"+n+"}", v)
	}

	return compiled
}

func (s *Server) Receive(c *gin.Context) {
	u := c.Request.URL.String()
	raw, err := getRawRequest(c.Request)
	if err != nil {
		log.Warn(err.Error())
	}
	for _, _rule := range s.getRules() {
		flag, flagGroup, vars := _rule.Match(string(raw))
		if flag == "" {
			continue
		}

		ip := strings.Split(c.Request.RemoteAddr, ":")[0]

		if ip1 := c.Request.Header.Get(s.IpHeader); s.IpHeader != "" && ip1 != "" {
			ip = ip1
			delete(c.Request.Header, s.IpHeader)
		}

		// create new record
		r, err := NewRecord(_rule, flag, c.Request.Method, u, ip, ipinfo.Area(ip), string(raw))
		if err != nil {
			log.Warn("HTTP record[rule_id:%d] created failed :%s", _rule.ID, err)
			code, err := strconv.Atoi(compileTpl(c, _rule.ResponseStatusCode, vars))
			if err != nil || code < 100 || code > 600 {
				code = 400
			}

			c.String(code, compileTpl(c, _rule.ResponseBody, vars))
			return
		}
		log.Info("HTTP record[id:%d rule:%s remote_ip:%s] has been created", r.ID, _rule.Name, ip)

		//only send to client or notify user when this connection recorded first time.
		var count int64
		if flagGroup != "" {
			database.DB.Where("rule_name=? and raw_request like ?", _rule.Name, "%"+flagGroup+"%").Model(&Record{}).Count(&count)
		}
		if count <= 1 {
			if _rule.PushToClient {
				r.PushToClient()
				if flagGroup != "" {
					log.Trace("HTTP record[id:%d, flagGroup:%s] has been put to client message queue", r.ID, flagGroup)
				} else {
					log.Trace("HTTP record[id:%d] has been put to client message queue", r.ID)
				}
			}
			//send notice
			if _rule.Notice {
				go func() {
					r.Notice()
					if flagGroup != "" {
						log.Trace("HTTP record[id:%d, flagGroup:%s]  notice has been sent", r.ID, flagGroup)
					} else {
						log.Trace("HTTP record[id:%d] notice has been sent", r.ID)
					}
				}()
			}
		}

		for header, value := range _rule.ResponseHeaders {
			c.Header(compileTpl(c, header, vars), compileTpl(c, value, vars))
		}

		code, err := strconv.Atoi(compileTpl(c, _rule.ResponseStatusCode, vars))
		if err != nil || code < 100 || code > 600 {
			code = 400
		}

		c.String(code, compileTpl(c, _rule.ResponseBody, vars))
		return
	}

}
