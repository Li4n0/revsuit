package ftp

type Config struct {
	Enable   bool
	Addr     string
	PasvIP   string `yaml:"pasv_ip"`
	PasvPort int    `yaml:"pasv_port"`
}
