package service

import (
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

// GeoService looks up the country ISO code for an IP address.
// If no database path was configured it is a no-op and always returns "".
type GeoService struct {
	db *geoip2.Reader // nil when GeoDBPath is empty
}

func NewGeoService(dbPath string) *GeoService {
	if dbPath == "" {
		return &GeoService{}
	}
	db, err := geoip2.Open(dbPath)
	if err != nil {
		log.Printf("geo service: failed to open GeoIP database at %q: %v (country lookup disabled)", dbPath, err)
		return &GeoService{}
	}
	log.Printf("geo service: loaded GeoIP database from %q", dbPath)
	return &GeoService{db: db}
}

// LookupCountry returns the ISO 3166-1 alpha-2 country code (e.g. "US", "DE")
// for the given IP string. Returns "" if lookup is unavailable or fails.
func (g *GeoService) LookupCountry(ipStr string) string {
	if g.db == nil {
		return ""
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}
	record, err := g.db.Country(ip)
	if err != nil {
		return ""
	}
	return record.Country.IsoCode
}

func (g *GeoService) Close() {
	if g.db != nil {
		g.db.Close()
	}
}
