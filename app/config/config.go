// Package config defines application's configurations
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"gopkg.in/throttled/throttled.v1"

	"github.com/supinf/reinvent-sessions-api/app/misc"
)

func defaultConfig() Config {
	gopath := os.Getenv("GOPATH")
	return Config{
		Name:            "ReInvent-Sessions-API",
		Port:            80,
		Stage:           "",
		LogLevel:        4,
		AccessLog:       true,
		StaticFileHost:  "",
		StaticFilePath:  gopath + "/src/github.com/supinf/reinvent-sessions-api/app",
		ValidHost:       "",
		ValidUserAgent:  "",
		CorsMethods:     "",
		CorsOrigin:      "*",
		SecuredCookie:   false,
		Timeout:         60 * time.Second,
		LimitRatePerMin: 0,
		LimitBursts:     0,
		LimitVaryBy:     nil,
		LimitKeyCache:   0,
		AwsLog:          false,
		AwsRoleExpiry:   5 * time.Minute,
		DynamoDbLocal:   "",
		CognitoPoolID:   "us-east-1:12345",
		CognitoRoleArn:  "arn:aws:iam::12345",
		TwitterKey:      "",
		TwitterSecret:   "",
		TwitterCallback: "",
	}
}

// NewConfig returns a config struct created by
// merging environment variables and a config file.
func NewConfig() *Config {
	temp := environmentConfig()
	config := &temp

	if !config.complete() {
		config.merge(fileConfig())
	}
	defer func() {
		config.merge(defaultConfig())
		config.trimWhitespace()
	}()
	return config
}

func environmentConfig() Config {
	sTemp := os.Getenv("APP_LIMIT_VARYBY")
	var varyBy *throttled.VaryBy
	if sTemp == "Path" {
		varyBy = &throttled.VaryBy{Path: true}
	}
	if sTemp == "RemoteAddr" {
		varyBy = &throttled.VaryBy{RemoteAddr: true}
	}
	return Config{
		Name:            os.Getenv("APP_NAME"),
		Port:            misc.ParseUint16(os.Getenv("APP_PORT")),
		Stage:           os.Getenv("APP_STAGE"),
		LogLevel:        misc.Atoi(os.Getenv("APP_LOG_LEVEL")),
		AccessLog:       misc.ParseBool(os.Getenv("APP_ACCESS_LOG")),
		StaticFileHost:  os.Getenv("APP_STATIC_FILE_HOST"),
		StaticFilePath:  os.Getenv("APP_STATIC_FILE_PATH"),
		ValidHost:       os.Getenv("APP_VALID_HOST"),
		ValidUserAgent:  os.Getenv("APP_VALID_USER_AGENT"),
		CorsMethods:     os.Getenv("APP_ACCESS_CONTROL_ALLOW_METHODS"),
		CorsOrigin:      os.Getenv("APP_ACCESS_CONTROL_ALLOW_ORIGIN"),
		SecuredCookie:   misc.ParseBool(os.Getenv("APP_SECURED_COOKIE")),
		Timeout:         misc.ParseDuration(os.Getenv("APP_TIMEOUT")),
		LimitRatePerMin: misc.Atoi(os.Getenv("APP_LIMITRATE_PERMIN")),
		LimitBursts:     misc.Atoi(os.Getenv("APP_LIMIT_BURST")),
		LimitVaryBy:     varyBy,
		LimitKeyCache:   misc.Atoi(os.Getenv("APP_LIMIT_KEYCACHE")),
		AwsLog:          misc.ParseBool(os.Getenv("APP_AWS_LOG")),
		AwsRoleExpiry:   misc.ParseDuration(os.Getenv("APP_AWS_ROLE_EXPIRY")),
		DynamoDbLocal:   os.Getenv("DYNAMODB_PORT_8000_TCP_ADDR"),
		CognitoPoolID:   os.Getenv("APP_COGNITO_POOL_ID"),
		CognitoRoleArn:  os.Getenv("APP_COGNITO_ROLE_ARN"),
		TwitterKey:      os.Getenv("APP_TWITTER_CONSUMER_KEY"),
		TwitterSecret:   os.Getenv("APP_TWITTER_CONSUMER_SECRET"),
		TwitterCallback: os.Getenv("APP_TWITTER_CONSUMER_CALLBACK"),
	}
}

func fileConfig() Config {
	path := misc.NVL(os.Getenv("CONFIG_FILE_PATH"), "/etc/reinvent-sessions-api/config.json")
	file, err := os.Open(path)
	if err != nil {
		return Config{}
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Unable to read config file", "err:", err)
		return Config{}
	}
	if strings.TrimSpace(string(data)) == "" {
		return Config{}
	}
	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Error reading config json data. [message] ", err)
	}
	return config
}

func (config *Config) merge(arg Config) *Config {
	mine := reflect.ValueOf(config).Elem()
	theirs := reflect.ValueOf(&arg).Elem()

	for i := 0; i < mine.NumField(); i++ {
		myField := mine.Field(i)
		if misc.ZeroOrNil(myField.Interface()) {
			myField.Set(reflect.ValueOf(theirs.Field(i).Interface()))
		}
	}
	return config
}

func (config *Config) complete() bool {
	cfgs := reflect.ValueOf(config).Elem()

	for i := 0; i < cfgs.NumField(); i++ {
		if misc.ZeroOrNil(cfgs.Field(i).Interface()) {
			return false
		}
	}
	return true
}

func (config *Config) trimWhitespace() {
	cfgs := reflect.ValueOf(config).Elem()
	cfgAttrs := reflect.Indirect(reflect.ValueOf(config)).Type()

	for i := 0; i < cfgs.NumField(); i++ {
		field := cfgs.Field(i)
		if !field.CanInterface() {
			continue
		}
		attr := cfgAttrs.Field(i).Tag.Get("trim")
		if len(attr) == 0 {
			continue
		}
		if field.Kind() != reflect.String {
			continue
		}
		str := field.Interface().(string)
		field.SetString(strings.TrimSpace(str))
	}
}

// String returns a string representation of the config.
func (config *Config) String() string {
	return fmt.Sprintf(
		"Name: %v, Port: %v, Stage: %v, LogLevel: %v, AccessLog: %v, "+
			"StaticFileHost: %v, StaticFilePath: %v, "+
			"ValidHost: %v, ValidUserAgent: %v, CorsMethods: %v, CorsOrigin: %v, SecuredCookie: %v, "+
			"Timeout: %v, LimitRatePerMin: %v, LimitBursts: %v, LimitVaryBy: %v, LimitKeyCache: %v, "+
			"AwsRegion: %v, AwsLog: %v, AwsRoleExpiry: %v, DynamoDbLocal: %v, CognitoPoolID: %v, "+
			"CognitoRoleArn: %v, TwitterKey: %v, TwitterSecret: %v, TwitterCallback: %v",
		config.Name, config.Port, config.Stage, config.LogLevel, config.AccessLog,
		config.StaticFileHost, config.StaticFilePath,
		config.ValidHost, config.ValidUserAgent, config.CorsMethods, config.CorsOrigin, config.SecuredCookie,
		config.Timeout, config.LimitRatePerMin, config.LimitBursts, config.LimitVaryBy, config.LimitKeyCache,
		os.Getenv("AWS_REGION"), config.AwsLog, config.AwsRoleExpiry, config.DynamoDbLocal,
		config.CognitoPoolID, config.CognitoRoleArn,
		config.TwitterKey, config.TwitterSecret, config.TwitterCallback)
}
