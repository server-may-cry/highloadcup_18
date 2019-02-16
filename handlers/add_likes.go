package handlers

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/server-may-cry/highloadcup_18/db"
	"github.com/valyala/fasthttp"
)

func GetAddLikesHandler(storage *db.Storage) fasthttp.RequestHandler {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return func(ctx *fasthttp.RequestCtx) {
		var request struct {
			Likes []db.Like `json:"likes"`
		}
		b := ctx.PostBody()
		err := json.Unmarshal(b, &request)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}

		for _, like := range request.Likes {
			if _, exists := storage.Get(like.Liker); !exists {
				ctx.SetStatusCode(http.StatusBadRequest)
				return
			}
			if _, exists := storage.Get(like.Likee); !exists {
				ctx.SetStatusCode(http.StatusBadRequest)
				return
			}
		}

		storage.AddLikes(request.Likes)
		ctx.SetStatusCode(http.StatusAccepted)
		ctx.SetContentTypeBytes(contentType)
		ctx.SetBody(emptyBody)
	}
}
