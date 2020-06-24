package api

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	adapter "github.com/awslabs/aws-lambda-go-api-proxy/core"
)

var (
	accessor = &adapter.RequestAccessor{}
)

// LambdaHandler is the func that should be passed to the `lambda.Start` call
func (a *API) LambdaHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req, err := accessor.EventToRequestWithContext(ctx, event)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	writer := adapter.NewProxyResponseWriter()

	a.ServeHTTP(writer, req)

	response, err := writer.GetProxyResponse()
	if err != nil {
		return adapter.GatewayTimeout(), errors.New("gateway timeout while processing request")
	}

	return response, nil
}
