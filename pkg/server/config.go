package server

import (
	"github.com/li4n0/revsuit/pkg/dns"
	"github.com/li4n0/revsuit/pkg/ftp"
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
	Version    float64
	Addr       string
	Token      string
	Domain     string
	ExternalIP string `yaml:"external_ip"`
	Database   string
	LogLevel   string `yaml:"log_level"`
	Notice     noticeConfig
	rhttp.Config
	DNS   dns.Config
	MySQL mysql.Config
	RMI   rmi.Config
	FTP   ftp.Config
}
