package main

import (
	"github.com/kataras/iris/httptest"
	"testing"
)

func TestTerminal(t *testing.T) {
	e := httptest.New(t, newApp(), httptest.URL("http://127.0.0.1:8080"))

	e.GET("/api/v1.0/terminal").Expect().Body().Contains("OK")
}
