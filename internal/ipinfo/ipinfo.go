package ipinfo

import (
	"github.com/li4n0/revsuit/internal/ipinfo/geoip"
	"github.com/li4n0/revsuit/internal/ipinfo/qqwry"
	log "unknwon.dev/clog/v2"
)

var db Database

func Area(ip string) string {
	if db != nil {
		return db.Area(ip)
	}
	return ""
}

func Init(config Config) {
	switch config.Database {
	case "qqwry":
		db = qqwry.New()
	case "geoip":
		db = geoip.New(config.GeoLicenseKey)
	default:
		log.Fatal("wrong ip location database type")
	}
}
