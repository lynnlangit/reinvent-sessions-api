package config

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/throttled/throttled"
)

func TestDefaultConfig(t *testing.T) {
	var expected interface{}
	actual := defaultConfig()

	expected = "ReInvent-Sessions-API"
	if actual.Name != expected {
		t.Errorf("Expected %v, but got %v", expected, actual.Name)
		return
	}
	expected = uint16(80)
	if actual.Port != expected {
		t.Errorf("Expected %v, but got %v", expected, actual.Port)
		return
	}
	expected = true
	if actual.AccessLog != expected {
		t.Errorf("Expected %v, but got %v", expected, actual.AccessLog)
		return
	}
	var varyBy *throttled.VaryBy
	if actual.LimitVaryBy != varyBy {
		t.Errorf("Expected %v, but got %v", expected, actual.LimitVaryBy)
		return
	}
}

func TestMerge(t *testing.T) {
	cfg := Config{
		Name:        "Test",
		Port:        8080,
		LogLevel:    6,
		LimitVaryBy: &throttled.VaryBy{RemoteAddr: true},
	}
	actual := cfg.merge(defaultConfig())
	gopath := os.Getenv("GOPATH")
	expected := &Config{
		Name:            "Test",
		Port:            8080,
		LogLevel:        6,
		AccessLog:       true,
		StaticFileHost:  "",
		StaticFilePath:  gopath + "/src/github.com/supinf/reinvent-sessions-api/app",
		Timeout:         60 * time.Second,
		LimitRatePerMin: 0,
		LimitBursts:     0,
		LimitVaryBy:     &throttled.VaryBy{RemoteAddr: true},
		LimitKeyCache:   0,
		DynamoDbLocal:   "",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}

func TestComplete(t *testing.T) {
	actual := defaultConfig()
	if actual.complete() {
		t.Errorf("Unexpected result. %v", actual)
		return
	}
	actual = *actual.merge(Config{
		AccessLog:       true,
		StaticFileHost:  "some-host",
		LimitRatePerMin: 1,
		LimitBursts:     1,
		LimitVaryBy:     &throttled.VaryBy{RemoteAddr: true},
		LimitKeyCache:   1,
		DynamoDbLocal:   "dynamo",
	})
	if !actual.complete() {
		t.Errorf("Unexpected result. %v", actual)
		return
	}
}

func TestTrimWhitespace(t *testing.T) {
	actual := defaultConfig()
	actual.Name = " 　a b 　　c "
	actual.trimWhitespace()

	expected := "a b 　　c"
	if actual.Name != expected {
		t.Errorf("Expected %v, but got %v", expected, actual.Name)
		return
	}
}
