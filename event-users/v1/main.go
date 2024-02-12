package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type User struct {
	PK       string `json:"pk"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Phone    string `json:"phone"`
	Date     string `json:"date"`
	Status   string `json:"status"`
	EventID  string `json:"eventid"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	svc := dynamodb.New(sess)

	eventID := request.QueryStringParameters["param"]

	queryInput := &dynamodb.QueryInput{
		TableName: aws.String("eventstable"),
		IndexName: aws.String("eventid-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"eventid": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("Event#" + eventID),
					},
				},
			},
		},
	}

	result, err := svc.Query(queryInput)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	var users []User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	usersJson, err := json.Marshal(users)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(usersJson),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
