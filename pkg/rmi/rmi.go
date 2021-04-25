package rmi

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/qqwry"
	log "unknwon.dev/clog/v2"
)

type Server struct {
	Config
	rules     []*Rule
	rulesLock sync.RWMutex
}

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
	defer s.rulesLock.Unlock()
	s.rulesLock.Lock()
	return db.Order("rank desc").Find(&s.rules).Error
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(time.Second * 30)); err != nil {
		log.Error("RMI set connection deadline error:%v", err.Error())
	}

	ip, port, _ := net.SplitHostPort(conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Error("RMI read connection error:%v", err.Error())
	}

	if !bytes.Contains(buf, []byte{0x4a, 0x52, 0x4d, 0x49}) {
		return
	}

	send := []byte{0x4e}
	bs := make([]byte, 8)
	binary.BigEndian.PutUint16(bs, uint16(len(ip)))
	send = append(send, bs...)
	send = append(send, []byte(ip)...)
	send = append(send, []byte{0x00, 0x00}...)
	uintPort, _ := strconv.Atoi(port)
	bs = make([]byte, 8)
	binary.BigEndian.PutUint16(bs, uint16(uintPort))
	send = append(send, bs...)

	_, err = conn.Write(send)
	if err != nil {
		log.Error("RMI write connection error: %v", err.Error())
	}

	data := make([]byte, 512)

	for length := 0; length < 50; {
		n, err := conn.Read(data)
		if err != nil {
			log.Error("RMI read connection error: %v", err.Error())
		}
		length += n
	}

	frags := bytes.Split(data, []byte{0xdf, 0x74})
	path := strings.TrimRight(string(frags[len(frags)-1][2:]), "\x00")

	for _, _rule := range s.getRules() {
		flag, flagGroup, _ := _rule.Match(path)
		if flag == "" {
			continue
		}

		area := qqwry.Area(ip)

		// create new record
		r, err := NewRecord(_rule, flag, path, ip, area)
		if err != nil {
			log.Error("RMI record[rule_id:%d] created failed :%s", _rule.ID, err.Error())
			return
		}
		log.Info("RMI record[id:%d rule:%s remote_ip:%s] has been created", r.ID, _rule.Name, ip)

		//only send to client when this connection recorded first time.
		if _rule.PushToClient {
			if flagGroup != "" {
				var count int64
				database.DB.Where("rule_name=? and raw like ?", _rule.Name, "%"+flagGroup+"%").Model(&Record{}).Count(&count)
				if count <= 1 {
					r.PushToClient()
					log.Trace("RMI record[id%d] has been put to client message queue", r.ID)
				}
			}
			r.PushToClient()
			log.Trace("RMI record[id%d] has been put to client message queue", r.ID)
		}

		//send notice
		if _rule.Notice {
			go func() {
				r.Notice()
				log.Trace("RMI record[id%d] notice has been sent", r.ID)
			}()
		}
	}
}

func (s *Server) Run() {
	if err := s.updateRules(); err != nil {
		log.Fatal(err.Error())
	}

	// run server
	log.Info("Starting RMI Server at %v", s.Addr)

	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Error("RMI accept connection error: %v", err.Error())
			continue
		}
		go s.handleConnection(tcpConn)
	}

}
