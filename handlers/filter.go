package handlers

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/server-may-cry/highloadcup_18/db"
	"github.com/valyala/fasthttp"
)

type FilterResponse struct {
	Accounts []map[string]interface{} `json:"accounts"`
}

func GetFilterHandler(storage *db.Storage) fasthttp.RequestHandler {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var allowedFilters = map[string]struct{}{
		"sex_eq": {},

		"email_domain": {},
		"email_lt":     {},
		"email_gt":     {},

		"status_eq":  {},
		"status_neq": {},

		"fname_eq":   {},
		"fname_any":  {},
		"fname_null": {},

		"sname_eq":     {},
		"sname_starts": {},
		"sname_null":   {},

		"phone_code": {},
		"phone_null": {},

		"country_eq":   {},
		"country_null": {},

		"city_eq":   {},
		"city_any":  {},
		"city_null": {},

		"birth_lt":   {},
		"birth_gt":   {},
		"birth_year": {},

		"interests_contains": {},
		"interests_any":      {},

		"likes_contains": {},

		"premium_now":  {},
		"premium_null": {},

		"limit": {},

		"query_id": {},
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

		ctx.SetContentTypeBytes(contentType)
		filtered := storage.Filter(db.Filter{
			SexEq: queryParameters.Peek("sex_eq"),

			EmailDomain: queryParameters.Peek("email_domain"),
			EmailLt:     queryParameters.Peek("email_lt"),
			EmailGt:     queryParameters.Peek("email_gt"),

			StatusEq:  queryParameters.Peek("status_eq"),
			StatusNeq: queryParameters.Peek("status_neq"),

			FnameEq:   queryParameters.Peek("fname_eq"),
			FnameAny:  stringWithComasToStrings(queryParameters.Peek("fname_any")),
			FnameNull: queryParameters.Peek("fname_null"),

			SnameEq:     queryParameters.Peek("sname_eq"),
			SnameStarts: queryParameters.Peek("sname_starts"),
			SnameNull:   queryParameters.Peek("sname_null"),

			PhoneCode: queryParameters.Peek("phone_code"),
			PhoneNull: queryParameters.Peek("phone_null"),

			CountryEq:   queryParameters.Peek("country_eq"),
			CountryNull: queryParameters.Peek("country_null"),

			CityEq:   queryParameters.Peek("city_eq"),
			CityAny:  stringWithComasToStrings(queryParameters.Peek("city_any")),
			CityNull: queryParameters.Peek("city_null"),

			BirthLt:   queryParameters.GetUintOrZero("birth_lt"),
			BirthGt:   queryParameters.GetUintOrZero("birth_gt"),
			BirthYear: int16(queryParameters.GetUintOrZero("birth_year")),

			InterestsContains: stringWithComasToStrings(queryParameters.Peek("interests_contains")),
			InterestsAny:      stringWithComasToStrings(queryParameters.Peek("interests_any")),

			LikesContains: stringWithComasToInts(queryParameters.Peek("likes_contains")),

			PremiumNow:  queryParameters.Peek("premium_now"),
			PremiumNull: queryParameters.Peek("premium_null"),
		}, limitInt)
		if len(filtered) == 0 {
			ctx.SetBody(emptyFilterResponse)
			return
		}
		response := FilterResponse{
			Accounts: make([]map[string]interface{}, len(filtered)),
		}
		for i, v := range filtered {
			response.Accounts[i] = map[string]interface{}{
				"id":    v.ID,
				"email": v.Email,
			}

			if queryParameters.Has("sex_eq") {
				response.Accounts[i]["sex"] = v.Sex
			}
			if queryParameters.Has("status_eq") || queryParameters.Has("status_neq") {
				response.Accounts[i]["status"] = v.Status
			}
			if queryParameters.Has("fname_eq") || queryParameters.Has("fname_any") || queryParameters.Has("fname_null") {
				if v.Fname != "" {
					response.Accounts[i]["fname"] = v.Fname
				}
			}
			if queryParameters.Has("sname_eq") || queryParameters.Has("sname_starts") || queryParameters.Has("sname_null") {
				if v.Sname != "" {
					response.Accounts[i]["sname"] = v.Sname
				}
			}

			if queryParameters.Has("phone_code") || queryParameters.Has("phone_null") {
				if v.Phone != "" {
					response.Accounts[i]["phone"] = v.Phone
				}
			}
			if queryParameters.Has("country_eq") || queryParameters.Has("country_null") {
				if v.Country != "" {
					response.Accounts[i]["country"] = v.Country
				}
			}
			if queryParameters.Has("city_eq") || queryParameters.Has("city_any") || queryParameters.Has("city_null") {
				if v.City != "" {
					response.Accounts[i]["city"] = v.City
				}
			}
			if queryParameters.Has("birth_lt") || queryParameters.Has("birth_gt") || queryParameters.Has("birth_year") {
				response.Accounts[i]["birth"] = v.Birth
			}
			if queryParameters.Has("premium_now") || queryParameters.Has("premium_null") {
				if v.Premium.Start != 0 {
					response.Accounts[i]["premium"] = v.Premium
				}
			}
		}
		b, _ := json.Marshal(response)
		ctx.SetBody(b)
	}
}
