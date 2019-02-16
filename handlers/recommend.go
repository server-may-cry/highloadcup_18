package handlers

import (
	"math"
	"net/http"
	"sort"

	"github.com/server-may-cry/highloadcup_18/functions"

	jsoniter "github.com/json-iterator/go"
	"github.com/server-may-cry/highloadcup_18/db"
	"github.com/server-may-cry/highloadcup_18/httphelper"
	"github.com/server-may-cry/highloadcup_18/structures"
	"github.com/valyala/fasthttp"
)

type Account struct {
	ID      int                 `json:"id"`
	Email   string              `json:"email"`
	Status  string              `json:"status"`
	Fname   string              `json:"fname,omitempty"`
	Sname   string              `json:"sname,omitempty"`
	Birth   int                 `json:"birth"`
	Premium *structures.Premium `json:"premium,omitempty"`
}

type RecomendResponse struct {
	Accounts []Account `json:"accounts"`
}

type accountWithCompatibleIndex struct {
	compatibleIndex int64
	account         structures.Account
}

func GetRecomendHandler(storage *db.Storage) fasthttp.RequestHandler {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var allowedFilters = map[string]struct{}{
		"country":  {},
		"city":     {},
		"limit":    {},
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
		accountID, err := httphelper.GetAccountID(ctx.Path())
		if err != nil {
			ctx.SetStatusCode(http.StatusNotFound)
			return
		}
		a, exist := storage.Get(accountID)
		if !exist {
			ctx.SetStatusCode(http.StatusNotFound)
			return
		}
		ctx.SetContentTypeBytes(contentType)
		searchSex := "f"
		if a.Sex == "f" {
			searchSex = "m"
		}
		if len(a.Interests) == 0 {
			ctx.SetBody(emptyRecommendBody)
			return
		}
		filtered := storage.Filter(db.Filter{
			SexEq:        []byte(searchSex),
			CountryEq:    queryParameters.Peek("country"),
			CityEq:       queryParameters.Peek("city"),
			InterestsAny: a.Interests,
		}, math.MaxInt32)
		if len(filtered) == 0 {
			ctx.SetBody(emptyRecommendBody)
			return
		}
		compatibleList := make([]accountWithCompatibleIndex, 0, len(filtered))
		for _, v := range filtered {
			compatibleList = append(compatibleList, accountWithCompatibleIndex{
				account:         v,
				compatibleIndex: functions.CalculateCompatibility(a, v),
			})
		}
		lenCompatibleList := len(compatibleList)
		sort.SliceStable(compatibleList, func(i, j int) bool {
			return compatibleList[i].compatibleIndex > compatibleList[j].compatibleIndex
		})
		if lenCompatibleList < limitInt {
			limitInt = lenCompatibleList
		}
		response := RecomendResponse{
			Accounts: make([]Account, limitInt),
		}
		for i, v := range compatibleList[:limitInt] {
			response.Accounts[i].ID = v.account.ID
			response.Accounts[i].Email = v.account.Email
			response.Accounts[i].Status = v.account.Status
			response.Accounts[i].Fname = v.account.Fname
			response.Accounts[i].Sname = v.account.Sname
			response.Accounts[i].Birth = v.account.Birth
			if v.account.Premium.Start != 0 {
				cp := v.account.Premium
				response.Accounts[i].Premium = &cp
			}
		}
		b, _ := json.Marshal(response)
		ctx.SetBody(b)
	}
}
