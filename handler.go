package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/guregu/dynamo"
	"github.com/yunspace/serverless-golang/aws/event/apigateway"
)

func DynamoDB() dynamo.Table {
	db := dynamo.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	return db.Table(os.Getenv("STAGE") + "-Todos")
}

func Ping(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	return apigateway.NewAPIGatewayResponseWithBody(200, "Pong"), nil
}

func Echo(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	var message struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(evt.Body), &message); err != nil {
		return apigateway.NewAPIGatewayResponseWithError(400, err), nil
	}
	return apigateway.NewAPIGatewayResponseWithBody(200, message), nil
}

func CreateTodo(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	var jsonParams struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	if err := json.Unmarshal([]byte(evt.Body), &jsonParams); err != nil {
		return apigateway.NewAPIGatewayResponseWithError(400, err), nil
	}
	todo := NewTodo(jsonParams.Title, jsonParams.Completed)
	if err := DynamoDB().Put(todo).Run(); err != nil {
		return apigateway.NewAPIGatewayResponseWithError(502, err), nil
	}
	return apigateway.NewAPIGatewayResponseWithBody(201, todo), nil
}

func DeleteTodo(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	var deletedTodo Todo
	ID := evt.PathParameters["id"]
	if ID == "" {
		return apigateway.NewAPIGatewayResponseWithBody(400, "Invalid query string"), nil
	}
	if err := DynamoDB().Delete("ID", ID).OldValue(&deletedTodo); err != nil {
		return apigateway.NewAPIGatewayResponseWithError(502, err), nil
	}
	return apigateway.NewAPIGatewayResponseWithBody(200, deletedTodo), nil
}

func FetchAllTodo(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	var todos []Todo
	if err := DynamoDB().Scan().All(&todos); err != nil {
		return apigateway.NewAPIGatewayResponseWithError(502, err), nil
	}
	return apigateway.NewAPIGatewayResponseWithBody(200, todos), nil
}

func FetchSingleTodo(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	ID := evt.PathParameters["id"]
	if ID == "" {
		return apigateway.NewAPIGatewayResponseWithBody(400, "Invalid query string"), nil
	}
	var todo Todo
	if err := DynamoDB().Get("ID", ID).One(&todo); err != nil {
		return apigateway.NewAPIGatewayResponseWithError(502, err), nil
	}
	return apigateway.NewAPIGatewayResponseWithBody(200, todo), nil
}

func UpdateTodo(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	var jsonParams struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	var newTodo Todo
	ID := evt.PathParameters["id"]
	if ID == "" {
		return apigateway.NewAPIGatewayResponseWithBody(400, "Invalid query string"), nil
	}
	if err := json.Unmarshal([]byte(evt.Body), &jsonParams); err != nil {
		return apigateway.NewAPIGatewayResponseWithError(400, err), nil
	}
	if err := DynamoDB().Update("ID", ID).Add("Title", jsonParams.Title).Add("Completed", jsonParams.Completed).Value(&newTodo); err != nil {
		return apigateway.NewAPIGatewayResponseWithError(502, err), nil
	}
	return apigateway.NewAPIGatewayResponseWithBody(200, newTodo), nil
}
