package ftp

type Config struct {
	Enable   bool
	Addr     string
	PasvPort int `yaml:"pasv_port"`
}
