package handlers

import (
	"bytes"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/server-may-cry/highloadcup_18/db"
	"github.com/valyala/fasthttp"
)

type GroupResponse struct {
	Groups []map[string]interface{} `json:"groups"`
}

func GetGroupHandler(storage *db.Storage) fasthttp.RequestHandler {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var allowedFilters = map[string]struct{}{
		"birth":     {},
		"keys":      {},
		"sex":       {},
		"interests": {},
		"country":   {},
		"joined":    {},
		"city":      {},
		"order":     {},
		"limit":     {},
		"query_id":  {},
		"likes":     {},
		"status":    {},
	}
	var allowedKeys = map[string]struct{}{
		"sex":       {},
		"status":    {},
		"interests": {},
		"country":   {},
		"city":      {},
	}
	return func(ctx *fasthttp.RequestCtx) {
		queryParameters := ctx.URI().QueryArgs()
		var stopRequest bool
		queryParameters.VisitAll(func(key, val []byte) {
			if _, ok := allowedFilters[string(key)]; !ok || len(val) == 0 {
				ctx.SetStatusCode(http.StatusBadRequest)
				stopRequest = true
			}
		})
		if stopRequest {
			return
		}
		limitInt, err := queryParameters.GetUint("limit")
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		if limitInt < 0 {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		order := string(queryParameters.Peek("order"))
		if order != "" && order != "-1" && order != "1" {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		keys := bytes.Split(queryParameters.Peek("keys"), []byte(","))
		keysStrings := make([]string, 0, len(keys))
		for _, key := range keys {
			skey := string(key)
			if _, ok := allowedKeys[skey]; !ok {
				ctx.SetStatusCode(http.StatusBadRequest)
				return
			}
			keysStrings = append(keysStrings, skey)
		}
		ctx.SetContentTypeBytes(contentType)
		grouped := storage.Group(db.Filter{
			StatusEq:          queryParameters.Peek("status"),
			LikesContains:     stringWithComasToInts(queryParameters.Peek("likes")),
			SexEq:             queryParameters.Peek("sex"),
			CountryEq:         queryParameters.Peek("country"),
			BirthYear:         int16(queryParameters.GetUintOrZero("birth")),
			InterestsContains: stringWithComasToStrings(queryParameters.Peek("interests")),
			JoinedYear:        int16(queryParameters.GetUintOrZero("joined")),
			CityEq:            queryParameters.Peek("city"),
		}, keysStrings, limitInt, order == "-1")
		if len(grouped) == 0 {
			ctx.SetBody(emptyGroupBody)
			return
		}
		groups := make([]map[string]interface{}, 0, len(grouped))
		for _, v := range grouped {
			m := make(map[string]interface{}, len(v.Keys)+1)
			m["count"] = v.Count
			for i, vv := range v.Keys {
				m[i] = vv
			}
			groups = append(groups, m)
		}
		response := GroupResponse{
			Groups: groups,
		}
		b, _ := json.Marshal(response)
		ctx.SetBody(b)
	}
}
