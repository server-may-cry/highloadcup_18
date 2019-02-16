package handlers

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/server-may-cry/highloadcup_18/db"
	"github.com/server-may-cry/highloadcup_18/structures"
	"github.com/server-may-cry/highloadcup_18/validators"
	"github.com/valyala/fasthttp"
)

func GetNewAccountHandler(storage *db.Storage) fasthttp.RequestHandler {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return func(ctx *fasthttp.RequestCtx) {
		var request structures.Account
		b := ctx.PostBody()
		err := json.Unmarshal(b, &request)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		if !validators.IsAccountValid(request) {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		err = storage.AddAccount(request)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		ctx.SetStatusCode(http.StatusCreated)
		ctx.SetContentTypeBytes(contentType)
		ctx.SetBody(emptyBody)
	}
}
