package main

import (
	"testing"

	"github.com/kataras/iris/httptest"
)

func TestBasicAuth(t *testing.T) {

	e := httptest.New(t, newApp())

	e.GET("/hi").Expect().Status(httptest.StatusOK)

}