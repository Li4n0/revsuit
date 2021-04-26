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
	log "unknwon.dev/clog/v2"
)

type Server struct {
	rules     []*Rule
	rulesLock sync.RWMutex
}

var (
	server         *Server
	once           sync.Once
	rebindingCache = cache.New(5*time.Second, 10*time.Second)
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

func (s *Server) Run() {

	if err := s.updateRules(); err != nil {
		log.Fatal(err.Error())
	}

	//create new dns zone with root domain
	newZone := func(name string) *newdns.Zone {
		defer func() {
			if err := recover(); err != nil {
				recycler.Recycle(err)
			}
		}()

		domain := strings.TrimSuffix(name, ".")
		frags := strings.Split(domain, ".")
		zoneName := ""
		if len(frags) >= 2 {
			zoneName = strings.Join(frags[len(frags)-2:], ".") + "."
		} else {
			zoneName = name
		}
		return &newdns.Zone{
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

					//only send to client when this connection recorded first time.
					if _rule.PushToClient {
						if flagGroup != "" {
							var count int64
							database.DB.Where("rule_name=? and domain like ?", _rule.Name, "%"+flagGroup+"%").Model(&Record{}).Count(&count)
							if count <= 1 {
								r.PushToClient()
								log.Trace("DNS record[id%d] has been put to client message queue", r.ID)
							}
						} else {
							r.PushToClient()
							log.Trace("DNS record[id%d] has been put to client message queue", r.ID)
						}
					}

					//send notice
					if _rule.Notice {
						go func() {
							r.Notice()
							log.Trace("DNS record[id%d] notice has been sent", r.ID)
						}()
					}

					if _rule.Value != "" {
						value := rule.CompileTpl(_rule.Value, vars)
						_type := _rule.Type
						if _rule.Type == newdns.REBINDING {
							_type = newdns.A
						}

						return []newdns.Set{
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

										log.Trace("DNS rebinding client(ip:%v) to %v", ip, value)
										return []newdns.Record{{Address: value}}
									default:
										return []newdns.Record{{Address: value}}
									}
								}(),
								TTL: _rule.TTL * time.Second,
							},
						}, nil
					}
				}

				return nil, nil
			},
		}
	}

	// create server
	server := newdns.NewServer(newdns.Config{
		Handler: func(name string) (*newdns.Zone, error) {
			return newZone(name), nil
		},
	})

	// run server
	log.Info("Starting DNS Server at :53")
	err := server.Run(":53")
	if err != nil {
		log.Fatal(err.Error())
	}

}
