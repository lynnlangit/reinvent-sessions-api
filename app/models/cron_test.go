package models

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/supinf/reinvent-sessions-api/app/config"
	"github.com/supinf/reinvent-sessions-api/app/misc"
)

func TestCronTable(t *testing.T) {
	actual := cronTable
	r, _ := regexp.Compile("[^a-zA-Z0-9_\\.]")
	expected := r.ReplaceAllString(strings.ToLower(config.NewConfig().Name), "-") + "-crons"
	if actual != expected {
		t.Errorf("CronTableName was not fixed in init func. Expected %v, but got %v", expected, actual)
		return
	}
}

func TestToCronResult(t *testing.T) {
	actual := toCronResult(map[string]*dynamodb.AttributeValue{
		"Key": &dynamodb.AttributeValue{
			S: awssdk.String("key-name"),
		},
		"LastStartDate": &dynamodb.AttributeValue{
			S: awssdk.String("2015-08-01T12:03:04Z"),
		},
		"LastEndDate": &dynamodb.AttributeValue{
			S: awssdk.String("2015-08-02T23:45:06+0900"),
		},
	})
	expected := CronResult{}
	expected.Key = "key-name"
	expected.LastStartDate = misc.StringToTime("2015-08-01T12:03:04Z")
	expected.LastEndDate = misc.StringToTime("2015-08-02T23:45:06+0900")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}
