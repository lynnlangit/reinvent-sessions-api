package models

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/supinf/reinvent-sessions-api/app/config"
)

func TestSessionTable(t *testing.T) {
	actual := sessionTable
	r, _ := regexp.Compile("[^a-zA-Z0-9_\\.]")
	expected := r.ReplaceAllString(strings.ToLower(config.NewConfig().Name), "-") + "-sessions"
	if actual != expected {
		t.Errorf("CronTableName was not fixed in init func. Expected %v, but got %v", expected, actual)
		return
	}
}

func TestToSessionResult(t *testing.T) {
	actual := toSession(map[string]*dynamodb.AttributeValue{
		"ID": &dynamodb.AttributeValue{
			S: awssdk.String("123"),
		},
		"TypeId": &dynamodb.AttributeValue{
			N: awssdk.String(strconv.Itoa(123)),
		},
	})
	expected := Session{}
	expected.ID = "123"
	expected.TypeID = 123

	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
		return
	}
}
