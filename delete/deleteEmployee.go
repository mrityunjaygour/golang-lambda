package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	lambda.Start(DeleteEmployeeHandler)
}

func DeleteEmployeeHandler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := r.PathParameters["id"]
	err := deleteEmployee(id)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func deleteEmployee(id string) error {

	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(session)

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("employees"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	_, err := svc.DeleteItem(input)
	if err != nil {
		return err
	}
	return nil
}
