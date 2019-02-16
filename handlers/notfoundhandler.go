package handlers

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

func NotFoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusNotFound)
}
