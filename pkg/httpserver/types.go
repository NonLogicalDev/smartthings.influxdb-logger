package httpserver

import "net/http"

type (
	Handler        = func(w http.ResponseWriter, r *http.Request)
	RequestHandler interface {
		HandleRequest(ctx *RequestContext)
	}
)
