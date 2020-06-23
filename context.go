package api

import "context"

var (
	contextValueKeyAPI = &struct{}{}
)

// Context returns a new context with the API
func Context(ctx context.Context, api *API) context.Context {
	return api.context(context.WithValue(ctx, contextValueKeyAPI, api))
}
