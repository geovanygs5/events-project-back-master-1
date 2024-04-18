package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Event struct {
	Pk       string `json:"pk"`
	Sk       string `json:"sk"`
	Name     string `json:"name"`
	Capacity string `json:"capacity"`
	Date     string `json:"date"`
	Hour     string `json:"hour"`
	Status   string `json:"status"`
	Email    string `json:"email"`
	Eventid  string `json:"eventid"`
	Gender   string `json:"gender"`
	Lastname string `json:"lastname"`
	Phone    string `json:"phone"`
	Price    string `json:"price"`
	Bank     string `json:"bank"`
	Account  string `json:"account"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	paramValue := request.QueryStringParameters["param"]

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	svc := dynamodb.New(sess)
	input := &dynamodb.QueryInput{
		TableName: aws.String("events-testing-table"),
		IndexName: aws.String("sk-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"sk": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("METADATA#" + paramValue),
					},
				},
			},
		},
	}

	result, err := svc.Query(input)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	var eventsList []Event
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &eventsList)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	responseBody, err := json.Marshal(eventsList)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}

func main() {
	lambda.Start(handler)
}
