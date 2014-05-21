package web

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteString(t *testing.T) {
	str := "Hello World!"
	m := martini.Classic()
	m.Use(ContextWithCookieSecret("secret"))
	m.Get("/", func(ctx *Context) {
		ctx.WriteString(str)
	})
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	m.ServeHTTP(res, req)
	if res.Body.String() != str {
		t.Errorf("WriteString Error")
	}
}

func TestAbort(t *testing.T) {
	str := "Hello World!"
	m := martini.Classic()
	m.Use(ContextWithCookieSecret("secret"))
	m.Get("/", func(ctx *Context) {
		ctx.Abort(401, str)
	})
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	m.ServeHTTP(res, req)
	if res.Code != 401 {
		t.Error("Response Code Error")
	}
	if res.Body.String() != str {
		t.Error("Abort Content Error")
	}
}
