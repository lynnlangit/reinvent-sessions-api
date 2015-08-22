package config

import (
	"time"

	"github.com/throttled/throttled"
)

// Config defines the application configurations
type Config struct {
	Name            string `trim:"true"`
	Port            uint16
	Stage           string
	LogLevel        int
	AccessLog       bool
	StaticFileHost  string `trim:"true"`
	StaticFilePath  string `trim:"true"`
	ValidHost       string `trim:"true"`
	ValidUserAgent  string `trim:"true"`
	CorsMethods     string `trim:"true"`
	CorsOrigin      string `trim:"true"`
	Timeout         time.Duration
	LimitRatePerMin int
	LimitBursts     int
	LimitVaryBy     *throttled.VaryBy
	LimitKeyCache   int
	AwsLog          bool
	AwsRoleExpiry   time.Duration
	DynamoDbLocal   string `trim:"true"`
}
