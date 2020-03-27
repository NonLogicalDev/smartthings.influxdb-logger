package httpserver

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
)

type (
	Handler = func(w http.ResponseWriter, r *http.Request)
	RequestHandler interface {
		HandleRequest (ctx *RequestContext)
	}
)

func RegisterHandlers(root string, influxUrl string, mux *http.ServeMux)  {
	wrap := func(handler RequestHandler) Handler  {
		return func(w http.ResponseWriter, r *http.Request) {
			rCtx := &RequestContext{
				Context: context.Background(),
				Id: rand.Intn(9999999),
				ResW: w, Req: r,
			}

			defer func() {
				r := recover()
				if r == nil {
					return
				}

				if err, ok := r.(error); ok {
					rCtx.Log("PANIC", "%v", err)
				}
			}()

			rCtx.Log("|START",">>>>>>>>>>>")
			handler.HandleRequest(rCtx)
			rCtx.Log("|END  ","<<<<<<<<<<<")
		}
	}

	mux.HandleFunc(
		fmt.Sprintf("%s/smtdata", root),
		wrap(NewSMTHandler(influxUrl)),
	)
}