package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Employee struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Salary string `json:"salary"`
}

func main() {
	lambda.Start(GetEmployee)
}

//GetEmployee handler takes event request
func GetEmployee(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := r.PathParameters["id"]
	emp, err := fetchEmployee(id)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "",
		}, err
	}

	empData, err := json.Marshal(emp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(empData),
	}, nil
}

func fetchEmployee(id string) (Employee, error) {
	var emp Employee
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(session)

	input := &dynamodb.GetItemInput{
		TableName: aws.String("employees"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	result, err := svc.GetItem(input)
	if err != nil {
		return emp, err
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, &emp)
	if err != nil {
		return emp, err
	}
	return emp, nil
}
