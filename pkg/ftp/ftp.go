package ftp

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/qqwry"
	"github.com/li4n0/revsuit/internal/rule"
	log "unknwon.dev/clog/v2"
)

type Server struct {
	Config
	rules     []*Rule
	rulesLock sync.RWMutex
}

type Status string

const (
	CRASHED  Status = "CRASHED"
	FINISHED Status = "FINISHED"
)

var (
	server *Server
	once   sync.Once
)

func GetServer() *Server {
	once.Do(func() {
		server = &Server{rulesLock: sync.RWMutex{}}
	})
	return server
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

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(time.Second * 30)); err != nil {
		log.Error("FTP set connection deadline error:%v", err.Error())
	}

	if _, err := conn.Write([]byte("220 (vsFTPd 3.0.2)\r\n")); err != nil {
		log.Error("FTP write connection error:%v", err.Error())
	}

	ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
	buf := &bytes.Buffer{}

	var user, password, path, flag, flagGroup string
	status := CRASHED
	var matchedRule *Rule
	var vars map[string]string

loop:
	for {
		data := make([]byte, 2048)
		n, err := conn.Read(data)
		if err != nil {
			break
		}
		buf.Write(data[:n])

		if buf.Len() > 4 {
			cmd := string(buf.Bytes()[:4])
			switch cmd {
			case "USER":
				user = strings.TrimRight(string(buf.Bytes()[5:]), "\r\n")
				_, _ = conn.Write([]byte("331 password please - version check\r\n"))
			case "PASS":
				password = strings.TrimRight(string(buf.Bytes()[5:]), "\r\n")
				_, _ = conn.Write([]byte("230 User logged in\r\n"))

				for _, _rule := range s.getRules() {
					for _, s := range []string{user, password} {
						flag, flagGroup, vars = _rule.Match(s)
						if flag != "" {
							vars["user"] = user
							vars["password"] = password
							break
						}
					}
					if flag == "" {
						continue
					}
					matchedRule = _rule
				}

			case "QUIT":
				_, _ = conn.Write([]byte("221 Goodbye.\r\n"))
			case "RETR":
				path += "/" + strings.TrimRight(string(buf.Bytes()[5:]), "\r\n")
				_, _ = conn.Write([]byte("451 Nope\r\n"))
				_, _ = conn.Write([]byte("221 Goodbye.\r\n"))
				status = FINISHED
				break loop
			case "EPSV", "EPRT", "PORT":
				// refuse to use EPSV/EPRT/PORT in order to make the client to use PASV mode.
				_, _ = conn.Write([]byte(fmt.Sprintf("500 '%s': command not understood.\r\n", cmd)))
			case "PASV":
				// return rule's pasv_address or default pasv address
				ret := fmt.Sprintf("227 Entering Passive Mode (%s,%v,%d)\r\n", strings.ReplaceAll(s.PasvIP, ".", ","), float64(s.PasvPort/256), s.PasvPort%256)

				if matchedRule != nil && matchedRule.PasvAddress != "" {
					pasvAddress := rule.CompileTpl(matchedRule.PasvAddress, vars)
					pasvIP, pasvPort, err := net.SplitHostPort(pasvAddress)
					if err != nil {
						log.Warn("FTP failed to split rule[id%d] pasv_address(%s) :%s", matchedRule.ID, pasvAddress, err.Error())
						break
					}
					port, err := strconv.Atoi(pasvPort)
					if err != nil {
						log.Warn("FTP failed to convert rule[id%d] pasv_port(%s) :%s", matchedRule.ID, pasvPort, err.Error())
						break
					}
					ret = fmt.Sprintf("227 Entering Passive Mode (%s,%v,%d)\r\n", strings.ReplaceAll(pasvIP, ".", ","), float64(port/256), port%256)
				}
				_, _ = conn.Write([]byte(ret))
			default:
				cmd = string(buf.Bytes()[:3])
				if cmd == "CWD" {
					_, _ = conn.Write([]byte("250 Directory successfully changed.\r\n"))
					path += "/" + strings.TrimRight(string(buf.Bytes()[4:]), "\r\n")
				} else if cmd == "PWD" {
					_, _ = conn.Write([]byte("257 \"/\" is the current directory\r\n"))
				} else {
					_, _ = conn.Write([]byte("230 more data please!\r\n"))
				}
			}
		}
		buf = &bytes.Buffer{}
	}

	if matchedRule != nil {
		_rule := matchedRule
		area := qqwry.Area(ip)

		// create new record
		r, err := NewRecord(_rule, flag, user, password, path, ip, area, status)
		if err != nil {
			log.Error("FTP record[rule_id:%d] created failed :%s", _rule.ID, err.Error())
			return
		}
		log.Info("FTP record[id:%d rule:%s remote_ip:%s] has been created", r.ID, _rule.Name, ip)

		//only send to client when this connection recorded first time.
		if _rule.PushToClient {
			if flagGroup != "" {
				var count int64
				database.DB.Where("rule_name=? and raw like ?", _rule.Name, "%"+flagGroup+"%").Model(&Record{}).Count(&count)
				if count <= 1 {
					r.PushToClient()
					log.Trace("FTP record[id%d] has been put to client message queue", r.ID)
				}
			}
			r.PushToClient()
			log.Trace("FTP record[id%d] has been put to client message queue", r.ID)
		}

		//send notice
		if _rule.Notice {
			go func() {
				r.Notice()
				log.Trace("FTP record[id%d] notice has been sent", r.ID)
			}()
		}
	}
}

func (s *Server) Run() {
	if err := s.updateRules(); err != nil {
		log.Fatal(err.Error())
	}

	// run server
	log.Info("Starting FTP Server at %v", s.Addr)

	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		pasvAddress := fmt.Sprintf("%s:%d", strings.Split(s.Addr, ":")[0], s.PasvPort)
		log.Info("Start to listen FTP PASV port at %v", pasvAddress)
		listener, err := net.Listen("tcp", pasvAddress)
		if err != nil {
			log.Fatal("FTP failed to listen on pasv port : %v", err)
		}
		for {
			tcpConn, err := listener.Accept()
			if err != nil {
				log.Error("FTP accept connection error: %v", err)
				continue
			}
			_ = tcpConn.Close()
		}
	}()

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Error("FTP accept connection error: %v", err)
			continue
		}
		go s.handleConnection(tcpConn)
	}

}
