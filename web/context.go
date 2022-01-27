package web

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func (ctx *Context) Json(data interface{}, err error) error {
	status := http.StatusOK
	if err != nil {
		status = http.StatusInternalServerError
	}
	ctx.w.Header().Set("Content-Type", "application/json")
	ctx.w.WriteHeader(status)
	encoder := json.NewEncoder(ctx.w)
	err = encoder.Encode(data)
	if err != nil {
		http.Error(ctx.w, err.Error(), 500)
	}
	return err
}

func (ctx *Context) Success(data interface{}) error {
	return ctx.Json(data, nil)
}

func (ctx *Context) Fail(data interface{}, err error) error {
	return ctx.Json(data, err)
}
