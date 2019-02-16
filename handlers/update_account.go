package handlers

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/server-may-cry/highloadcup_18/validators"
	"github.com/valyala/fasthttp"

	"github.com/server-may-cry/highloadcup_18/db"
	"github.com/server-may-cry/highloadcup_18/httphelper"
	"github.com/server-may-cry/highloadcup_18/structures"
)

func GetUpdateAccountHandler(storage *db.Storage) fasthttp.RequestHandler {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return func(ctx *fasthttp.RequestCtx) {
		accountID, err := httphelper.GetAccountID(ctx.Path())
		if err != nil {
			ctx.SetStatusCode(http.StatusNotFound)
			return
		}
		account, exist := storage.Get(accountID)
		if !exist {
			ctx.SetStatusCode(http.StatusNotFound)
			return
		}

		var accountFromRequest structures.Account
		b := ctx.PostBody()
		err = json.Unmarshal(b, &accountFromRequest)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}

		if accountFromRequest.Fname != "" {
			account.Fname = accountFromRequest.Fname
		}
		if accountFromRequest.Sname != "" {
			account.Sname = accountFromRequest.Sname
		}
		if accountFromRequest.Email != "" {
			account.Email = accountFromRequest.Email
		}
		// interests
		if len(accountFromRequest.Interests) != 0 {
			account.Interests = accountFromRequest.Interests
		}
		if accountFromRequest.Status != "" {
			account.Status = accountFromRequest.Status
		}
		if accountFromRequest.Sex != "" {
			account.Sex = accountFromRequest.Sex
		}
		if accountFromRequest.Phone != "" {
			account.Phone = accountFromRequest.Phone
		}
		// likes
		if len(accountFromRequest.Likes) != 0 {
			account.Likes = accountFromRequest.Likes
		}
		if accountFromRequest.Birth != 0 {
			account.Birth = accountFromRequest.Birth
		}
		if accountFromRequest.City != "" {
			account.City = accountFromRequest.City
		}
		if accountFromRequest.Country != "" {
			account.Country = accountFromRequest.Country
		}
		if accountFromRequest.Joined != 0 {
			account.Joined = accountFromRequest.Joined
		}
		if valid := validators.IsAccountValid(account); !valid {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		err = storage.UpdateAccount(account)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}

		ctx.SetStatusCode(http.StatusAccepted)
		ctx.SetContentTypeBytes(contentType)
		ctx.SetBody(emptyBody)
	}
}
