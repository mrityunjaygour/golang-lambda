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
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type Employee struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Salary string `json:"salary"`
}

// type UpdateInfo struct {
// 	Id     string `json:"id"`
// 	Name   string `json:"newname,omitempty"`
// 	Email  string `json:"newemail,omitempty"`
// 	Salary string `json:"newsalary,omitempty"`
// }

func main() {
	lambda.Start(UpdateHandler)
}

func UpdateHandler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var emp Employee
	err := json.Unmarshal([]byte(r.Body), &emp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "",
		}, err
	}
	updatedEmp, err := updateEmployee(emp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "",
		}, err
	}
	data, _ := json.Marshal(updatedEmp)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       string(data),
	}, nil
}

func updateEmployee(emp Employee) (Employee, error) {

	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(session)

	upd := expression.
		Set(expression.Name("name"), expression.Value(emp.Name)).
		Set(expression.Name("email"), expression.Value(emp.Email)).
		Set(expression.Name("salary"), expression.Value(emp.Salary))

	expr, err := expression.NewBuilder().WithUpdate(upd).Build()
	if err != nil {
		return emp, err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("employees"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(emp.Id),
			},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}

	result, err := svc.UpdateItem(input)
	if err != nil {
		return emp, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Attributes, &emp)
	if err != nil {
		return emp, err
	}

	return emp, nil
}
