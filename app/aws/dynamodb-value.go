package aws

/**
 * @see https://github.com/aws/aws-sdk-go/blob/master/service/dynamodb/api.go
 */

import (
	"strconv"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoS gets string from AttributeValue
func DynamoS(data map[string]*dynamodb.AttributeValue, key string) string {
	if value, ok := data[key]; ok {
		return *value.S
	}
	return ""
}

// DynamoN gets int from AttributeValue
func DynamoN(data map[string]*dynamodb.AttributeValue, key string) int {
	if value, ok := data[key]; ok {
		i, _ := strconv.Atoi(*value.N)
		return i
	}
	return 0
}

// DynamoAttributeS makes string to DynamoDB String Attribute
func DynamoAttributeS(value string) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		S: awssdk.String(value),
	}
}

// DynamoAttributeN makes int to DynamoDB Number Attribute
func DynamoAttributeN(value int) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		N: awssdk.String(strconv.Itoa(value)),
	}
}
