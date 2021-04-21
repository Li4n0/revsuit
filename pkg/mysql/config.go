package mysql

type Config struct {
	Enable bool
	Addr string
	VersionString string `yaml:"version_string"`
}
