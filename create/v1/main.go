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

type MyObject struct {
	Pk       string  `json:"pk"`
	Sk       string  `json:"sk"`
	Name     string  `json:"name"`
	Lastname *string `json:"lastname,omitempty"`
	Gender   *string `json:"gender,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	EventID  *string `json:"eventid,omitempty"`
	Capacity *string `json:"capacity,omitempty"`
	Date     string  `json:"date"`
	Hour     *string `json:"hour,omitempty"`
	Status   string  `json:"status"`
}

func CrearRegistro(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var object MyObject

	err := json.Unmarshal([]byte(request.Body), &object)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(object)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("eventstable"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(CrearRegistro)
}
