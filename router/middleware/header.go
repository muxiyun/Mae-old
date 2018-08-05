package middleware

import (
	"github.com/kataras/iris"
	"time"
	"net/http"
)

func NoCache(ctx iris.Context) {
	ctx.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	ctx.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	ctx.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	ctx.Next()
}

// Options is a middleware function that appends headers
// for options requests and aborts then exits the middleware
// chain and ends the request.
func Options(ctx iris.Context) {
	if ctx.Method() != "OPTIONS" {
		ctx.Next()
	} else {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		ctx.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		ctx.Header("Content-Type", "application/json")
		ctx.StatusCode(http.StatusOK)
	}
}

// Secure is a middleware function that appends security
// and resource access headers.
func Secure(ctx iris.Context) {

	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("X-Frame-Options", "DENY")
	ctx.Header("X-Content-Type-Options", "nosniff")
	ctx.Header("X-XSS-Protection", "1; mode=block")
	if ctx.Request().TLS != nil {
		ctx.Header("Strict-Transport-Security", "max-age=31536000")
	}
	ctx.Next()

	// Also consider adding Content-Security-Policy headers
	// c.Header("Content-Security-Policy", "script-src 'self' https://cdnjs.cloudflare.com")
}