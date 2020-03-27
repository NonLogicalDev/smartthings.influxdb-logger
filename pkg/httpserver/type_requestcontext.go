package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RequestContext struct {
	context.Context

	Id   int
	ResW http.ResponseWriter
	Req  *http.Request

	cacheBody    []byte
	cacheBodyErr error
}

func (rCtx *RequestContext) WriteJSON(data interface{}) {
	_ = json.NewEncoder(rCtx.ResW).Encode(data)
}

func (rCtx *RequestContext) WriteStatus(code int) {
	rCtx.ResW.WriteHeader(code)
}

func (rCtx *RequestContext) WriteHeader(name, value string) {
	rCtx.ResW.Header().Add(name, value)
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

func (rCtx *RequestContext) Log(header string, msg string, args ...interface{}) {
	fmt.Printf("[%s]:(%d) %s: %s\n",
		time.Now().Format(time.RFC3339),
		rCtx.Id,
		header,
		fmt.Sprintf(msg, args...),
	)
}

