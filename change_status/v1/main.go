package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Request struct {
	Pk     string `json:"pk"`
	Sk     string `json:"sk"`
	Status string `json:"status"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req Request
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	err = updateItemStatus(req.Pk, req.Sk, req.Status)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func updateItemStatus(pk string, sk string, newStatus string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		return err
	}

	svc := dynamodb.New(sess)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#S": aws.String("status"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String(newStatus),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(pk),
			},
			"sk": {
				S: aws.String(sk),
			},
		},
		TableName:        aws.String("eventstable"),
		UpdateExpression: aws.String("SET #S = :s"),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
