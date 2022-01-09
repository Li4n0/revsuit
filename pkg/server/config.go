package server

import (
	"github.com/li4n0/revsuit/internal/ipinfo"
	"github.com/li4n0/revsuit/pkg/dns"
	"github.com/li4n0/revsuit/pkg/ftp"
	"github.com/li4n0/revsuit/pkg/ldap"
	"github.com/li4n0/revsuit/pkg/mysql"
	"github.com/li4n0/revsuit/pkg/rhttp"
	"github.com/li4n0/revsuit/pkg/rmi"
)

type noticeConfig struct {
	DingTalk string
	Lark     string
	WeiXin   string
	Slack    string
}

type Config struct {
	Version            float64
	Addr               string
	Token              string
	Domains            []string
	ExternalIP         string `yaml:"external_ip"`
	AdminPathPrefix    string `yaml:"admin_path_prefix"`
	Database           string
	LogLevel           string        `yaml:"log_level"`
	CheckUpgrade       bool          `yaml:"check_upgrade"`
	IpLocationDatabase ipinfo.Config `yaml:"ip_location_database"`
	Notice             noticeConfig
	HTTP               rhttp.Config
	DNS                dns.Config
	MySQL              mysql.Config
	RMI                rmi.Config
	LDAP               ldap.Config
	FTP                ftp.Config
}
