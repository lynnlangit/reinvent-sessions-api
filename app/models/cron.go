package models

import (
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/supinf/reinvent-sessions-api/app/aws"
	"github.com/supinf/reinvent-sessions-api/app/config"
	"github.com/supinf/reinvent-sessions-api/app/logs"
	"github.com/supinf/reinvent-sessions-api/app/misc"
)

var cronTable string

// CronResult represents result of crons
type CronResult struct {
	Key           string    `json:"key"`
	Memo          string    `json:"memo"`
	LastStartDate time.Time `json:"lastStartDate"`
	LastEndDate   time.Time `json:"lastEndDate"`
}

var cronOnce sync.Once

func init() {
	cronOnce.Do(func() {
		r, _ := regexp.Compile("[^a-zA-Z0-9_\\.]")
		cronTable = r.ReplaceAllString(strings.ToLower(config.NewConfig().Name), "-") + "-crons"
		comfirmCronTableExists()
	})
}

// GetCronResults lists all cron results from DynamoDB
//  @return results models.CronResult
func GetCronResults() (result []CronResult, count int64, err error) {
	records, count, err := aws.DynamoRecords(cronTable)
	if err != nil {
		return result, 0, nil
	}
	return toCronResults(records), count, nil
}

// GetCronResult retrives a specified cron result from DynamoDB
//  @return session models.Session
func GetCronResult(key string) (result CronResult, err error) {
	record, err := aws.DynamoRecord(cronTable, map[string]*dynamodb.AttributeValue{
		"Key": aws.DynamoAttributeS(key),
	})
	if err != nil {
		return result, nil
	}
	return toCronResult(record), nil
}

// cast DynamoDB records to CronResult
func toCronResults(records []map[string]*dynamodb.AttributeValue) (results []CronResult) {
	for _, record := range records {
		result := toCronResult(record)
		results = append(results, result)
	}
	if len(results) == 0 {
		results = make([]CronResult, 0)
	}
	return results
}

func toCronResult(record map[string]*dynamodb.AttributeValue) CronResult {
	result := CronResult{}
	result.Key = aws.DynamoS(record, "Key")
	result.Memo = aws.DynamoS(record, "Memo")
	result.LastStartDate = misc.StringToTime(aws.DynamoS(record, "LastStartDate"))
	result.LastEndDate = misc.StringToTime(aws.DynamoS(record, "LastEndDate"))
	return result
}

// PersistCronResult saves itself to DynamoDB
func PersistCronResult(key string, memo string, s, e time.Time) error {
	items := map[string]*dynamodb.AttributeValue{}
	items["Key"] = aws.DynamoAttributeS(key)
	items["Memo"] = aws.DynamoAttributeS(memo)
	items["LastStartDate"] = aws.DynamoAttributeS(misc.TimeToString(s))
	items["LastEndDate"] = aws.DynamoAttributeS(misc.TimeToString(e))
	_, err := aws.DynamoPutItem(cronTable, items)
	return err
}

func comfirmCronTableExists() error {
	if _, err := aws.DynamoTable(cronTable); err == nil {
		return nil
	}
	logs.Debug.Print("[model] Cron table was not found. Try to make it. @aws.DynamoCreateTable")
	attributes := map[string]string{
		"Key": "S",
	}
	keys := map[string]string{
		"Key": "HASH",
	}
	_, err := aws.DynamoCreateTable(cronTable, attributes, keys, 1, 1)
	return err
}
