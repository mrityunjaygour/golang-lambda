package main

import (
	"encoding/json"
	"fmt"
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
	lambda.Start(PostHandler)
}

func PostHandler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var emp Employee
	err := json.Unmarshal([]byte(r.Body), &emp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       emp.Email,
		}, err
	}
	emp, err = createEmployee(emp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       emp.Email,
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       emp.Email + "created successfully",
	}, nil
}

func createEmployee(emp Employee) (Employee, error) {

	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(session)

	result, err := dynamodbattribute.MarshalMap(emp)
	if err != nil {
		fmt.Println("Failed to marshall request")
		return Employee{}, err
	}

	input := &dynamodb.PutItemInput{
		Item:      result,
		TableName: aws.String("employees"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Failed to write to db")
		return Employee{}, err
	}

	return Employee{Id: emp.Id, Name: emp.Name, Email: emp.Email, Salary: emp.Salary}, nil
}
