package main

import (
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
)

/// Create
func Ping(evt *apigatewayproxyevt.Event, ctx *runtime.Context) (interface{}, error) {
	return "Pong", nil
}
