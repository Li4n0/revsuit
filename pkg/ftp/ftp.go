package ftp

import (
	"bytes"
	"fmt"
	"go/types"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/file"
	"github.com/li4n0/revsuit/internal/qqwry"
	"github.com/li4n0/revsuit/internal/recycler"
	"github.com/li4n0/revsuit/internal/rule"
	"github.com/patrickmn/go-cache"
	log "unknwon.dev/clog/v2"
)

type Server struct {
	Config
	rules       []*Rule
	rulesLock   sync.RWMutex
	dataChannel chan map[string]interface{}
}

type Status string

const (
	CRASHED  Status = "CRASHED"
	FINISHED Status = "FINISHED"
)

type Method = string

const (
	DOWNLOAD Method = "DOWNLOAD"
	UPLOAD   Method = "UPLOAD"
)

var (
	server         *Server
	once           sync.Once
	rebindingCache = cache.New(5*time.Second, 10*time.Second)
)

func GetServer() *Server {
	once.Do(func() {
		server = &Server{rulesLock: sync.RWMutex{}, dataChannel: make(chan map[string]interface{}, 10)}
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
	defer s.rulesLock.Unlock()
	s.rulesLock.Lock()
	return db.Order("rank desc").Find(&s.rules).Error
}

func (s *Server) handleConnection(conn net.Conn) {
	log.Trace("New FTP connection from addr [%s]", conn.RemoteAddr())
	defer func() {
		_ = conn.Close()
		if err := recover(); err != nil {
			recycler.Recycle(err)
		}
	}()

	if err := conn.SetDeadline(time.Now().Add(time.Second * 30)); err != nil {
		log.Warn("FTP set connection deadline error:%v", err)
	}

	if _, err := conn.Write([]byte("220 (vsFTPd 3.0.2)\r\n")); err != nil {
		log.Warn("FTP write connection error:%v", err)
	}

	ip, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
	buf := &bytes.Buffer{}

	var user, password, method, flag, flagGroup, pasvAddress, filename string
	status := CRASHED
	path := "/"
	uploadData := make([]byte, 0)
	var matchedRule *Rule
	var vars map[string]string
	var isRedirect bool

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
			log.Trace("FTP connection[%s] exec command: %s", conn.RemoteAddr(), strings.TrimRight(buf.String(), "\r\n"))
			switch cmd {
			case "USER":
				user = strings.TrimRight(string(buf.Bytes()[5:]), "\r\n")
				_, _ = conn.Write([]byte("331 password please - version check\r\n"))
			case "PASS":
				password = strings.TrimRight(string(buf.Bytes()[5:]), "\r\n")
				for _, _rule := range s.getRules() {
					for _, s := range []string{user, password} {
						flag, flagGroup, vars = _rule.Match(s)
						if flag != "" {
							vars["user"] = user
							vars["password"] = password
							log.Trace(
								"FTP connection[%s] matched rule[rule_name: %s, flag: %s]",
								conn.RemoteAddr(), _rule.Name, flag,
							)
							break
						}
					}
					if flag == "" {
						continue
					}
					matchedRule = _rule
					pasvAddress = _rule.PasvAddress
				}

				if flag == "" {
					_, _ = conn.Write([]byte("331 please specify the password\r\n"))
					break loop
				}

				_, _ = conn.Write([]byte("230 User logged in\r\n"))
				if strings.Contains(pasvAddress, ",") {
					values, ok := rebindingCache.Get(ip)
					if !ok {
						rebindingCache.Set(ip, strings.Split(pasvAddress, ","), cache.DefaultExpiration)
						values = strings.Split(pasvAddress, ",")
					}
					//Choose and delete first address
					pasvAddress = values.([]string)[0]
					if len(values.([]string)) > 1 {
						rebindingCache.Set(ip, values.([]string)[1:len(values.([]string))], cache.DefaultExpiration)
					} else {
						rebindingCache.Delete(ip)
					}
				}
				if pasvAddress == "" {
					pasvAddress = fmt.Sprintf("%s:%d", s.PasvIP, s.PasvPort)
				}
				isRedirect = rule.CompileTpl(pasvAddress, vars) != fmt.Sprintf("%s:%d", s.PasvIP, s.PasvPort)
				if isRedirect {
					log.Trace("FTP connection[%s] will be redirect[pasv_address: %s]", conn.RemoteAddr(), pasvAddress)
				}
			case "SIZE":
				path += strings.TrimLeft(strings.TrimRight(string(buf.Bytes()[5:]), "\r\n"), "/")
				if matchedRule == nil || isRedirect || len(matchedRule.Data) == 0 {
					_, _ = conn.Write([]byte(fmt.Sprintf("550 %s: No such file or directory.\r\n", strings.TrimRight(string(buf.Bytes()[5:]), "\r\n"))))
					break
				}
				_, _ = conn.Write([]byte(fmt.Sprintf("213 %d\r\n", len(matchedRule.Data))))
			case "EPSV", "EPRT", "PORT":
				// refuse to use EPSV/EPRT/PORT in order to make the client to use PASV mode.
				_, _ = conn.Write([]byte(fmt.Sprintf("500 '%s': command not understood.\r\n", cmd)))
			case "PASV":
				// return rule's pasv_address or default pasv address
				if matchedRule != nil {
					pasvAddress := rule.CompileTpl(pasvAddress, vars)
					pasvIP, pasvPort, err := net.SplitHostPort(pasvAddress)
					if err != nil {
						log.Warn("FTP failed to split rule[id%d] pasv_address(%s) :%s", matchedRule.ID, pasvAddress, err)
						break
					}
					port, err := strconv.Atoi(pasvPort)
					if err != nil {
						log.Warn("FTP failed to convert rule[id%d] pasv_port(%s) :%s", matchedRule.ID, pasvPort, err)
						break
					}
					ret := fmt.Sprintf("227 Entering Passive Mode (%s,%v,%d)\r\n", strings.ReplaceAll(pasvIP, ".", ","), float64(port/256), port%256)
					_, _ = conn.Write([]byte(ret))
				}
			case "RETR":
				method = DOWNLOAD
				filename := strings.TrimRight(string(buf.Bytes()[5:]), "\r\n")
				if matchedRule == nil {
					_, _ = conn.Write([]byte("451 Nope\r\n"))
					_, _ = conn.Write([]byte("221 Goodbye.\r\n"))
					break
				}

				_, _ = conn.Write([]byte(
					fmt.Sprintf("150 Opening BINARY mode data connection for '%s' (%d bytes).\r\n", filename, len(matchedRule.Data))))
				dataPort, _ := strconv.Atoi(port)
				s.dataChannel <- map[string]interface{}{fmt.Sprintf("%s:%d", ip, dataPort+1): []byte(rule.CompileTpl(matchedRule.Data, vars))}
				_, _ = conn.Write([]byte("226 Transfer complete.\r\n"))
			case "STOR":
				method = UPLOAD
				filename = strings.TrimRight(string(buf.Bytes()[5:]), "\r\n")
				_, _ = conn.Write([]byte(fmt.Sprintf("150 Opening BINARY mode data connection for '%s'.\r\n", filename)))
				if !isRedirect {
					dataChannel := make(chan []byte)
					dataPort, _ := strconv.Atoi(port)
					s.dataChannel <- map[string]interface{}{fmt.Sprintf("%s:%d", ip, dataPort+1): dataChannel}
					uploadData = <-dataChannel
					log.Trace(
						"FTP connection[%s] uploaded %d bytes",
						conn.RemoteAddr(), len(uploadData),
					)
				}
				_, _ = conn.Write([]byte("226 Transfer complete.\r\n"))
			case "QUIT":
				_, _ = conn.Write([]byte("221 Goodbye.\r\n"))
				status = FINISHED
				break loop
			default:
				cmd = string(buf.Bytes()[:3])
				if cmd == "CWD" {
					_, _ = conn.Write([]byte("250 Directory successfully changed.\r\n"))
					path += strings.TrimRight(string(buf.Bytes()[4:]), "\r\n") + "/"
				} else if cmd == "PWD" {
					_, _ = conn.Write([]byte(fmt.Sprintf("257 \"%s\" is the current directory\r\n", path)))
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
		var r *Record
		var err error
		// create new record
		if len(uploadData) != 0 {
			r, err = NewRecord(_rule, flag, user, password, method, path, ip, area, file.FTPFile{Name: filename, Content: uploadData}, status)
		} else {
			r, err = NewRecord(_rule, flag, user, password, method, path, ip, area, file.FTPFile{}, status)
		}
		if err != nil {
			log.Warn("FTP record[rule_id:%d] created failed :%s", _rule.ID, err)
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

func (s *Server) handlePasvConnection(conn net.Conn, data map[string]interface{}) {
	remoteAddress := conn.RemoteAddr().String()
	switch v := data[remoteAddress].(type) {
	case types.Nil:
		s.dataChannel <- data
		return
	case []byte:
		_, err := conn.Write(v)
		if err != nil {
			log.Warn("FTP PASV server sent data to connection[%s] failed with error: %s", remoteAddress, err)
		}
		log.Trace("FTP PASV server has sent data to connection[%s]", remoteAddress)
	case chan []byte:
		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Warn("FTP PASV server received data from connection[%s] failed with error: %s", remoteAddress, err)
		}
		v <- buf
		log.Trace("FTP PASV server has received data from connection[%s]", remoteAddress)
	}

	_ = conn.Close()
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
		for data := range s.dataChannel {
			tcpConn, err := listener.Accept()
			if err != nil {
				log.Warn("FTP accept connection error: %v", err)
				continue
			}
			s.handlePasvConnection(tcpConn, data)
		}
	}()

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Warn("FTP accept connection error: %v", err)
			continue
		}
		go s.handleConnection(tcpConn)
	}

}
