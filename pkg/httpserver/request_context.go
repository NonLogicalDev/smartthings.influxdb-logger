package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

type RequestContext struct {
	context.Context

	Id  int

	Res http.ResponseWriter
	Req *http.Request
	Log *zap.Logger

	cacheBody    []byte
	cacheBodyErr error
}

func (rCtx *RequestContext) WriteJSON(data interface{}) {
	_ = json.NewEncoder(rCtx.Res).Encode(data)
}

func (rCtx *RequestContext) WriteStatus(code int) {
	rCtx.Res.WriteHeader(code)
}

func (rCtx *RequestContext) WriteHeader(name, value string) {
	rCtx.Res.Header().Add(name, value)
}

func (rCtx *RequestContext) ReadJSON(out interface{}) error {
	if rCtx.cacheBody == nil {
		rCtx.cacheBody, rCtx.cacheBodyErr = ioutil.ReadAll(rCtx.Req.Body)
	}
	if rCtx.cacheBodyErr != nil {
		return rCtx.cacheBodyErr
	}
	return json.Unmarshal(rCtx.cacheBody, out)
}

func (rCtx *RequestContext) OnError(code int, err error) {
	if err != nil {
		if errors.Is(err, &json.UnmarshalTypeError{}) {
			rCtx.Log = rCtx.Log.With(zap.String("request-body", string(rCtx.cacheBody)))
		}
		rCtx.WriteStatus(code)
		rCtx.WriteError(err)
		panic(err)
	}
}

func (rCtx *RequestContext) WriteError(err error) {
	if err != nil {
		rCtx.WriteJSON(map[string]interface{}{
			"status": "error",
			"error":  err,
		})
	}
}
