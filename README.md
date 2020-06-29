# Lambda API

[Documentation](https://godoc.org/github.com/planningcenter/lambda-api)

## Usage

### http.Handler

```golang
import (
  "net/http"

  api "github.com/planningcenter/lambda-api"
 "github.com/planningcenter/lambda-api/middleware"
)

func main() {
  app := api.New()

  app.Draw(func(router api.Router) {
    router.Add(middleware.Recovery, middleware.RequestID, middleware.Logger)

    router.Handle("GET", "/", func(w http.ResponseWriter, req *http.Request) {
      w.WriteHeader(http.StatusNoContent)
    })
  })

  server := &http.Server{
    Addr:    ":8080",
    Handler: app,
  }

  if err := server.ListenAndServe(); err != nil {
    log.Fatal(err)
  }
}
```

### Lambda

```golang
import (
  "net/http"

  "github.com/aws/aws-lambda-go/lambda"
  api "github.com/planningcenter/lambda-api"
 "github.com/planningcenter/lambda-api/middleware"
)

var (
  app *api.API
)

func init() {
  app = api.New()

  app.Draw(func(router api.Router) {
    router.Add(middleware.Recovery, middleware.RequestID, middleware.Logger)

    router.Handle("GET", "/", func(w http.ResponseWriter, req *http.Request) {
      w.WriteHeader(http.StatusNoContent)
    })
  })
}

func main() {
  lambda.Start(app.LambdaHandler)
}
```
