package api_test

import (
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	api "github.com/planningcenter/lambda-api"
	"github.com/planningcenter/lambda-api/middleware"
)

func ExampleAPI_LambdaHandler() {
	app := api.New()

	app.Draw(func(router api.Router) {
		router.Add(middleware.Recovery, middleware.RequestID, middleware.Logger)

		router.Handle("GET", "/", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	})

	lambda.Start(app.LambdaHandler)
}
