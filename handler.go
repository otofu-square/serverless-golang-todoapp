package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/guregu/dynamo"
	"github.com/satori/go.uuid"
	"github.com/yunspace/serverless-golang/aws/event/apigateway"
)

type Todo struct {
	ID        string `dynamo:"ID"`
	Title     string
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTodo(title string, completed bool) *Todo {
	id := uuid.NewV4().String()
	now := time.Now()
	return &Todo{
		ID:        id,
		Title:     title,
		Completed: completed,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func DynamoDB() dynamo.Table {
	db := dynamo.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	return db.Table("Todos")
}

func Ping(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	response := apigateway.NewAPIGatewayResponseWithBody(200, "Pong")
	return response, nil
}

func CreateTodo(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	todo := NewTodo("Test todo", false)
	err := DynamoDB().Put(todo).Run()
	if err == nil {
		return apigateway.NewAPIGatewayResponseWithBody(201, todo), nil
	} else {
		return apigateway.NewAPIGatewayResponseWithError(502, err), nil
	}
}
