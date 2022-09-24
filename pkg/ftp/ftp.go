package ftp

import (
	"bufio"
	"bytes"
	"fmt"
	"go/types"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/recycler"
	"github.com/li4n0/revsuit/internal/rule"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

type Server struct {
	Config
	pasvIP      string
	rules       []*Rule
	rulesLock   sync.RWMutex
	livingLock  sync.Mutex
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

func (s *Server) SetPasvIP(ip string) *Server {
	s.pasvIP = ip
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
	return errors.Wrap(db.Order("base_rank desc").Find(&s.rules).Error, "FTP update rules error")
}

func getClientPasvConnAddress(ip, port string) string {
	dataPort, _ := strconv.Atoi(port)
	return fmt.Sprintf("%s:%d", ip, dataPort+1)
}

func (s *Server) authenticate(user, password string) (_rule *Rule, flag, flagGroup string, vars map[string]string) {
	for _, _rule := range s.getRules() {
		for _, s := range []string{user, password} {
			flag, flagGroup, vars = _rule.Match(s)
			if flag != "" {
				vars["user"] = user
				vars["password"] = password
				return _rule, flag, flagGroup, vars
			}
		}
	}
	return _rule, flag, flagGroup, vars
}

func (s *Server) getPasvAddressFromCache(ip, pasvAddressTpl string) (pasvAddress string) {
	pasvAddress = pasvAddressTpl
	if strings.Contains(pasvAddressTpl, ",") {
		values, ok := rebindingCache.Get(ip)
		if !ok {
			rebindingCache.Set(ip, strings.Split(pasvAddressTpl, ","), cache.DefaultExpiration)
			values = strings.Split(pasvAddressTpl, ",")
		}
		//Choose and delete first address
		pasvAddress = values.([]string)[0]
		if len(values.([]string)) > 1 {
			rebindingCache.Set(ip, values.([]string)[1:len(values.([]string))], cache.DefaultExpiration)
		} else {
			rebindingCache.Delete(ip)
		}
	}
	return pasvAddress
}

const (
	NeedAccount             = "332 Need account for login.\r\n"
	PasswordPlease          = "331 password please - version check\r\n"
	PasswordError           = "331 please specify the password\r\n"
	UserLogged              = "230 User logged in\r\n"
	NoSuchFile              = "550 %s: No such file or directory.\r\n"
	CommandNotFound         = "500 '%s': command not understood.\r\n"
	EnteringPassiveMode     = "227 Entering Passive Mode (%s,%v,%d)\r\n"
	OpeningBinaryMode       = "150 Opening BINARY mode data connection for '%s' (%d bytes).\r\n"
	OpeningBinaryModeUpload = "150 Opening BINARY mode data connection for '%s'.\r\n"
	TransferComplete        = "226 Transfer complete.\r\n"
	Goodbye                 = "221 Goodbye.\r\n"
	DirectoryChanged        = "250 Directory successfully changed.\r\n"
	CurrentDirectory        = "257 \"%s\" is the current directory\r\n"
)

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			recycler.Recycle(err)
		}
	}()
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
	clientPasvConnAddress := getClientPasvConnAddress(ip, port)
	buf := &bytes.Buffer{}
	connBuf := bufio.NewWriter(conn)

	var user, password, method, flag, flagGroup, pasvAddress, filename string
	status := CRASHED
	path := "/"
	uploadData := make([]byte, 0)
	var _rule *Rule
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
			frags := strings.SplitN(strings.TrimRight(buf.String(), "\r\n"), " ", 2)
			var cmd = frags[0]
			var args string
			if len(frags) > 1 {
				args = frags[1]
			}
			log.Trace("FTP connection[%s] exec command: %s", conn.RemoteAddr(), strings.TrimRight(buf.String(), "\r\n"))

			if _rule == nil && cmd != "USER" && cmd != "PASS" {
				_, _ = connBuf.WriteString(NeedAccount)
				break loop
			}

			switch cmd {
			case "USER":
				user = args
				_, _ = connBuf.WriteString(PasswordPlease)
			case "PASS":
				password = args
				if _rule, flag, flagGroup, vars = s.authenticate(user, password); _rule == nil {
					_, _ = connBuf.WriteString(PasswordError)
					break loop
				}

				log.Trace("FTP connection[%s] matched rule[rule_name: %s, flag: %s]", conn.RemoteAddr(), _rule.Name, flag)
				_, _ = connBuf.WriteString(UserLogged)

				if pasvAddress = s.getPasvAddressFromCache(ip, _rule.PasvAddress); pasvAddress == "" {
					pasvAddress = fmt.Sprintf("%s:%d", s.pasvIP, s.PasvPort)
				}
				pasvAddress = rule.CompileTpl(pasvAddress, vars)
				isRedirect = pasvAddress != fmt.Sprintf("%s:%d", s.pasvIP, s.PasvPort)
			case "SIZE":
				path += strings.TrimLeft(args, "/")
				if _rule == nil || isRedirect || len(_rule.Data) == 0 {
					_, _ = connBuf.WriteString(fmt.Sprintf(NoSuchFile, args))
					break
				}
				_, _ = connBuf.WriteString(fmt.Sprintf("213 %d\r\n", len(_rule.Data)))
			case "EPSV", "EPRT", "PORT":
				// refuse to use EPSV/EPRT/PORT in order to make the client to use PASV mode.
				_, _ = connBuf.WriteString(fmt.Sprintf(CommandNotFound, cmd))
			case "PASV":
				//Just so that ide does not prompt that there may be a nil value
				if _rule != nil {
					pasvIP, pasvPort, err := net.SplitHostPort(pasvAddress)
					if err != nil {
						log.Warn("FTP failed to split rule[id:%d] pasv_address(%s) :%s", _rule.ID, pasvAddress, err)
						break
					}

					port, err := strconv.Atoi(pasvPort)
					if err != nil {
						log.Warn("FTP failed to convert rule[id:%d] pasv_port(%s) :%s", _rule.ID, pasvPort, err)
						break
					}
					ret := fmt.Sprintf(EnteringPassiveMode, strings.ReplaceAll(pasvIP, ".", ","), float64(port/256), port%256)
					_, _ = connBuf.WriteString(ret)
					if isRedirect {
						log.Trace("FTP connection[%s] will be redirect[pasv_address: %s]", conn.RemoteAddr(), pasvAddress)
					}
				}
			case "RETR":
				//Just so that ide does not prompt that there may be a nil value
				if _rule != nil {
					filename = args
					method = DOWNLOAD

					//send data to client
					_, _ = connBuf.WriteString(fmt.Sprintf(OpeningBinaryMode, filename, len(_rule.Data)))
					_ = connBuf.Flush()
					s.dataChannel <- map[string]interface{}{clientPasvConnAddress: []byte(rule.CompileTpl(_rule.Data, vars))}
					_, _ = connBuf.WriteString(TransferComplete)
				}

			case "STOR":
				filename = args
				method = UPLOAD

				_, _ = connBuf.WriteString(fmt.Sprintf(OpeningBinaryModeUpload, filename))
				_ = connBuf.Flush()
				//only could read data send to local pasv server.
				if !isRedirect {
					dataChannel := make(chan []byte)
					s.dataChannel <- map[string]interface{}{clientPasvConnAddress: dataChannel}
					uploadData = <-dataChannel
					log.Trace("FTP connection[%s] uploaded %d bytes", conn.RemoteAddr(), len(uploadData))
				}
				_, _ = connBuf.WriteString(TransferComplete)
			case "QUIT":
				_, _ = connBuf.WriteString(Goodbye)
				status = FINISHED
				break loop
			case "CWD":
				_, _ = connBuf.WriteString(DirectoryChanged)
				path += strings.TrimRight(args, "\r\n") + "/"
			case "PWD":
				_, _ = connBuf.WriteString(fmt.Sprintf(CurrentDirectory, path))
			default:
				_, _ = conn.Write([]byte("230 more data please!\r\n"))
			}
			_ = connBuf.Flush()
		}
		buf = &bytes.Buffer{}
	}

	if _rule != nil {
		createRecord(_rule, flag, flagGroup, user, password, method, path, filename, ip, uploadData, status)
	}
}

func (s *Server) handlePasvConnection(conn net.Conn, data map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			recycler.Recycle(err)
		}
	}()

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
		buf, err := io.ReadAll(conn)
		if err != nil {
			log.Warn("FTP PASV server received data from connection[%s] failed with error: %s", remoteAddress, err)
		}
		v <- buf
		log.Trace("FTP PASV server has received data from connection[%s]", remoteAddress)
	}

	_ = conn.Close()
}

// run pasv server
func (s *Server) runPasvServer() (net.Listener, error) {
	pasvAddress := fmt.Sprintf("%s:%d", strings.Split(s.Addr, ":")[0], s.PasvPort)
	log.Info("Start to listen FTP PASV port at %v, PasvIP is %v", pasvAddress, s.pasvIP)
	listener, err := net.Listen("tcp", pasvAddress)

	if err != nil {
		return nil, errors.Wrap(err, "FTP failed to listen on pasv port")
	}

	go func() {
		for data := range s.dataChannel {
			tcpConn, err := listener.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), net.ErrClosed.Error()) {
					log.Warn("FTP accept connection error: %v", err)
				} else {
					break
				}
				continue
			}
			s.handlePasvConnection(tcpConn, data)
		}
	}()
	return listener, nil
}

func (s *Server) Stop() {
	log.Info("FTP Server is stopping...")
	s.Enable = false
	s.livingLock.Unlock()
}

func (s *Server) Restart() {
	s.Stop()
	time.Sleep(time.Second * 2)
	go s.Run()
}

func (s *Server) Run() {
	s.Enable = true
	s.livingLock.Lock()
	defer func() {
		if s.Enable {
			log.Error("FTP Server exited unexpectedly")
			s.Enable = false
			s.livingLock.Unlock()
		}
	}()

	if err := s.UpdateRules(); err != nil {
		log.Error(err.Error())
		return
	}

	pasvListener, err := s.runPasvServer()
	if err != nil {
		log.Error(err.Error())
	}
	defer func() {
		if pasvListener != nil {
			_ = pasvListener.Close()
		}
	}()

	// run ftp server
	log.Info("Starting FTP Server at %v", s.Addr)

	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Error(errors.Wrap(err, "FTP failed to start").Error())
		return
	}

	go func() {
		s.livingLock.Lock()
		if !s.Enable {
			_ = listener.Close()
		}
		s.livingLock.Unlock()
	}()

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Warn("FTP accept connection error: %v", err)
			} else {
				break
			}
			continue
		}
		go s.handleConnection(tcpConn)
	}
}
