package httpserver

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"go.uber.org/zap"
)

func RegisterHandlers(log *zap.Logger, root string, influxUrl string, mux *http.ServeMux) {
	handler, err := NewSMTHandler(influxUrl)
	if err != nil {
		panic(err)
	}
	mux.HandleFunc(
		fmt.Sprintf("%s", root), wrapHandler(log, handler),
	)
}

func wrapHandler(log *zap.Logger, handler RequestHandler) Handler {
	return func(response http.ResponseWriter, request *http.Request) {
		reqID := rand.Intn(9999999)
		rCtx := &RequestContext{
			Context: context.Background(),
			Id:      reqID,
			Req:     request,
			Res:     response,
			Log:     log.Named("RequestHandler").With(
				zap.String("uri-path", request.URL.String()),
				zap.Int("trace-id", reqID),
			),
		}

		defer func() {
			r := recover()
			if r == nil {
				return
			}
			if err, ok := r.(error); ok {
				rCtx.Log.Error("request panicked", zap.Error(err))
			} else {
				rCtx.Log.Error("request panicked", zap.Any("panic", r))
			}
		}()


		rCtx.Log.Debug("request started")
		handler.HandleRequest(rCtx)
		rCtx.Log.Debug("request finished")
	}
}
