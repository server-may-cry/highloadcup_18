package handlers

import (
	"math"
	"net/http"
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/server-may-cry/highloadcup_18/db"
	"github.com/server-may-cry/highloadcup_18/functions"
	"github.com/server-may-cry/highloadcup_18/httphelper"
	"github.com/server-may-cry/highloadcup_18/listhelper"
	"github.com/server-may-cry/highloadcup_18/structures"
	"github.com/valyala/fasthttp"
)

type AccountResponse struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	Fname  string `json:"fname,omitempty"`
	Sname  string `json:"sname,omitempty"`
	Status string `json:"status"`
}

type SugggestResponse struct {
	Accounts []AccountResponse `json:"accounts"`
}

type accountWithSimilarityIndex struct {
	similarityIndex float64
	account         structures.Account
}

func GetSuggestHandler(storage *db.Storage) fasthttp.RequestHandler {
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
		if limitInt < 1 {
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
		if len(a.Likes) == 0 {
			ctx.SetBody(emptySuggestResponse)
			return
		}
		userLikes := make([]int, 0, len(a.Likes))
		for _, v := range a.Likes {
			userLikes = append(userLikes, v.ID)
		}
		filtered := storage.Filter(db.Filter{
			SexEq:            []byte(a.Sex),
			LikesContainsAny: userLikes,
			CountryEq:        queryParameters.Peek("country"),
			CityEq:           queryParameters.Peek("city"),
		}, math.MaxInt32)
		filteredLen := len(filtered)
		if filteredLen == 0 {
			ctx.SetBody(emptySuggestResponse)
			return
		}
		compatibleList := make([]accountWithSimilarityIndex, 0, filteredLen)
		for _, v := range filtered {
			compatibleList = append(compatibleList, accountWithSimilarityIndex{
				account:         v,
				similarityIndex: functions.CalculateSimularity(a.Likes, v.Likes),
			})
		}
		sort.SliceStable(compatibleList, func(i, j int) bool {
			return compatibleList[i].similarityIndex > compatibleList[j].similarityIndex
		})
		var listForSuggest []int
		for _, v := range compatibleList {
			sortedPartOfListToSuggest := make([]int, 0, len(v.account.Likes))
			for _, like := range v.account.Likes {
				sortedPartOfListToSuggest = append(sortedPartOfListToSuggest, like.ID)
			}
			sort.Sort(sort.Reverse(sort.IntSlice(sortedPartOfListToSuggest)))
			listForSuggest = append(listForSuggest, sortedPartOfListToSuggest...)
		}
		listForSuggest = listhelper.RemoveDuplicates(listForSuggest)
		finalListToSuggest := listhelper.Diff(listForSuggest, userLikes)
		finalListLen := len(finalListToSuggest)
		if limitInt > finalListLen {
			limitInt = finalListLen
		}
		responseAccounts := make([]AccountResponse, 0, limitInt)
		for _, v := range finalListToSuggest[:limitInt] {
			a, _ := storage.Get(v)
			responseAccounts = append(responseAccounts, AccountResponse{
				ID:     a.ID,
				Email:  a.Email,
				Fname:  a.Fname,
				Sname:  a.Sname,
				Status: a.Status,
			})
		}

		ctx.SetContentTypeBytes(contentType)
		response := SugggestResponse{
			Accounts: responseAccounts,
		}
		b, _ := json.Marshal(response)
		ctx.SetBody(b)
	}
}
