package main

import (
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
)

type APIGatewayResponse struct {
	StatusCode    int               `json:"statusCode"`
	Headers       map[string]string `json:"headers"`
	Body          interface{}       `json:"body"`
	Base64Encoded bool              `json:"isBase64Encoded"`
}

func Ping(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	response := APIGatewayResponse{
		StatusCode:    200,
		Headers:       nil,
		Body:          "pong",
		Base64Encoded: false,
	}
	return response, nil
}
