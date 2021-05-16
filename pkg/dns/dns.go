package dns

import (
	"strings"
	"sync"
	"time"

	"github.com/li4n0/revsuit/internal/database"
	"github.com/li4n0/revsuit/internal/newdns"
	"github.com/li4n0/revsuit/internal/qqwry"
	"github.com/li4n0/revsuit/internal/recycler"
	"github.com/li4n0/revsuit/internal/rule"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

type Server struct {
	Config
	rules      []*Rule
	rulesLock  sync.RWMutex
	livingLock sync.Mutex
}

var (
	server         *Server
	once           sync.Once
	rebindingCache = cache.New(5*time.Second, 10*time.Second)
)

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

func (s *Server) updateRules() error {
	db := database.DB.Model(new(Rule))
	defer s.rulesLock.Unlock()
	s.rulesLock.Lock()
	return errors.Wrap(db.Order("rank desc").Find(&s.rules).Error, "DNS update rules error")
}

func newSet(_rule *Rule, name, value, ip string, _type newdns.Type) []newdns.Set {
	set := []newdns.Set{
		{
			Name: name,
			Type: _type,
			Records: func() []newdns.Record {
				switch _rule.Type {
				case newdns.TXT:
					return []newdns.Record{{Data: []string{value}}}
				case newdns.CNAME, newdns.NS:
					return []newdns.Record{{Address: value + "."}}
				case newdns.REBINDING:

					// Get rebinding ip list
					values, ok := rebindingCache.Get(ip)
					if !ok {
						rebindingCache.Set(ip, strings.Split(value, ","), cache.DefaultExpiration)
						values = strings.Split(value, ",")
					}

					//Choose and delete first ip
					value := values.([]string)[0]
					if len(values.([]string)) > 1 {
						rebindingCache.Set(ip, values.([]string)[1:len(values.([]string))], cache.DefaultExpiration)
					} else {
						rebindingCache.Delete(ip)
					}

					log.Trace("DNS rebinding client[ip:%v] to %v", ip, value)
					return []newdns.Record{{Address: value}}
				default:
					return []newdns.Record{{Address: value}}
				}
			}(),
			TTL: _rule.TTL * time.Second,
		},
	}
	return set
}

// newZone creates new dns zone with root domain
func (s *Server) newZone(name string) *newdns.Zone {
	defer func() {
		if err := recover(); err != nil {
			recycler.Recycle(err)
		}
	}()

	domain := strings.TrimSuffix(name, ".")
	frags := strings.Split(domain, ".")
	zoneName := name
	if len(frags) >= 2 {
		zoneName = strings.Join(frags[len(frags)-2:], ".") + "."
	}
	zone := &newdns.Zone{
		Name:             zoneName,
		MasterNameServer: "ns1.hostmaster.com.",
		AllNameServers: []string{
			"ns1.hostmaster.com.",
			"ns2.hostmaster.com.",
			"ns3.hostmaster.com.",
		},
		Handler: func(lookedName, remoteAddr string) ([]newdns.Set, error) {
			ip := strings.Split(remoteAddr, ":")[0]

			for _, _rule := range s.getRules() {
				flag, flagGroup, vars := _rule.Match(domain)
				if flag == "" {
					continue
				}

				r, err := newRecord(_rule, flag, domain, ip, qqwry.Area(ip))
				if err != nil {
					log.Warn("DNS record(rule_id:%s) created failed :%s", _rule.Name, err)
					return nil, nil
				}
				log.Info("DNS record[id:%d rule:%s remote_ip:%s] has been created", r.ID, _rule.Name, ip)

				//only send to client or notify user when this connection recorded first time.
				var count int64
				if flagGroup != "" {
					database.DB.Where("rule_name=? and domain like ?", _rule.Name, "%"+flagGroup+"%").Model(&Record{}).Count(&count)
				}
				if count <= 1 {
					if _rule.PushToClient {
						r.PushToClient()
						if flagGroup != "" {
							log.Trace("DNS record[id:%d, flagGroup:%s] has been put to client message queue", r.ID, flagGroup)
						} else {
							log.Trace("DNS record[id:%d] has been put to client message queue", r.ID)
						}
					}
					//send notice
					if _rule.Notice {
						go func() {
							r.Notice()
							if flagGroup != "" {
								log.Trace("DNS record[id:%d, flagGroup:%s] notice has been sent", r.ID, flagGroup)
							} else {
								log.Trace("DNS record[id:%d] notice has been sent", r.ID)
							}
						}()
					}
				}
				if _rule.Value != "" {
					value := rule.CompileTpl(_rule.Value, vars)
					_type := _rule.Type
					if _rule.Type == newdns.REBINDING {
						_type = newdns.A
					}

					return newSet(_rule, name, value, ip, _type), nil
				}
			}

			return nil, nil
		},
	}
	return zone
}

func (s *Server) Stop() {
	log.Info("DNS server is stopping...")
	s.Enable = false
	s.livingLock.Unlock()
}

func (s *Server) Run() {
	s.Enable = true
	s.livingLock.Lock()

	defer func() {
		if s.Enable {
			log.Error("DNS Server exited unexpectedly")
		}
		s.Enable = false
		s.livingLock.Unlock()
	}()

	if err := s.updateRules(); err != nil {
		log.Error(err.Error())
		return
	}

	// create server
	server := newdns.NewServer(newdns.Config{
		Handler: func(name string) (*newdns.Zone, error) {
			return s.newZone(name), nil
		},
	})

	// run server
	log.Info("Starting DNS Server at :53")
	go func() {
		s.livingLock.Lock()
		if !s.Enable {
			server.Close()
		}
	}()

	err := server.Run(":53")
	defer server.Close()

	if err != nil {
		log.Error(err.Error())
		return
	}
}
