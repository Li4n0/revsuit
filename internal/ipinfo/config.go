package ipinfo

type Config struct {
	Database      string
	GeoLicenseKey string `yaml:"geo_license_key"`
}
