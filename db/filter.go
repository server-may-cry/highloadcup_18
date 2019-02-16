package db

import (
	"sort"
	"strings"

	"github.com/server-may-cry/highloadcup_18/listhelper"
	"github.com/server-may-cry/highloadcup_18/structures"
)

type Filter struct {
	AllowedIDs []int
	SexEq      []byte

	EmailDomain []byte
	EmailLt     []byte
	EmailGt     []byte

	StatusEq  []byte
	StatusNeq []byte

	FnameEq   []byte
	FnameAny  []string
	FnameNull []byte

	SnameEq     []byte
	SnameStarts []byte
	SnameNull   []byte

	PhoneCode []byte
	PhoneNull []byte

	CountryEq   []byte
	CountryNull []byte

	CityEq   []byte
	CityAny  []string
	CityNull []byte

	BirthLt    int
	BirthGt    int
	BirthYear  int16
	JoinedYear int16

	InterestsContains []string
	InterestsAny      []string

	LikesContains    []int
	LikesContainsAny []int

	PremiumNow  []byte
	PremiumNull []byte
}

func (s *Storage) FilterNoDataFetch(f Filter) []int {
	var filtered []int
	if len(f.AllowedIDs) != 0 {
		filtered = f.AllowedIDs
	} else {
		filtered = s.allIDs
	}

	if len(f.SexEq) != 0 {
		filtered = listhelper.Intersection(filtered, s.sexIndex[string(f.SexEq)])
	}
	if len(f.EmailDomain) != 0 {
		filtered = listhelper.Intersection(filtered, s.emailDomainIndex[string(f.EmailDomain)])
	}
	if len(f.StatusEq) != 0 {
		filtered = listhelper.Intersection(filtered, s.statusIndex[string(f.StatusEq)])
	}
	if len(f.FnameEq) != 0 {
		filtered = listhelper.Intersection(filtered, s.fnameIndex[string(f.FnameEq)])
	}
	if len(f.SnameEq) != 0 {
		filtered = listhelper.Intersection(filtered, s.snameIndex[string(f.SnameEq)])
	}
	if len(f.PhoneCode) != 0 {
		filtered = listhelper.Intersection(filtered, s.phoneCodeIndex[string(f.PhoneCode)])
	}
	if len(f.CountryEq) != 0 {
		filtered = listhelper.Intersection(filtered, s.countryIndex[string(f.CountryEq)])
	}
	if len(f.CityEq) != 0 {
		filtered = listhelper.Intersection(filtered, s.cityIndex[string(f.CityEq)])
	}
	if f.BirthYear != 0 {
		filtered = listhelper.Intersection(filtered, s.bithYearIndex[f.BirthYear])
	}
	if f.JoinedYear != 0 {
		filtered = listhelper.Intersection(filtered, s.joinedYearIndex[f.JoinedYear])
	}

	if len(f.StatusNeq) != 0 {
		var statuses []int
		for status, ids := range s.statusIndex {
			if status != string(f.StatusNeq) {
				flag := false
				statuses = appendOrIntersect(&flag, statuses, ids)
			}
		}
		sort.Ints(statuses)
		filtered = listhelper.Intersection(filtered, statuses)
	}
	if len(f.SnameStarts) != 0 {
		var names []int
		for name, ids := range s.snameIndex {
			if strings.HasPrefix(name, string(f.SnameStarts)) {
				flag := false
				names = appendOrIntersect(&flag, names, ids)
			}
		}
		sort.Ints(names)
		filtered = listhelper.Intersection(filtered, names)
	}

	if len(f.FnameNull) != 0 {
		filtered = nullFilter(f.FnameNull, filtered, s.fnameIndex)
	}
	if len(f.SnameNull) != 0 {
		filtered = nullFilter(f.SnameNull, filtered, s.snameIndex)
	}
	if len(f.CountryNull) != 0 {
		filtered = nullFilter(f.CountryNull, filtered, s.countryIndex)
	}
	if len(f.CityNull) != 0 {
		filtered = nullFilter(f.CityNull, filtered, s.cityIndex)
	}

	if len(f.PhoneNull) != 0 {
		filterStrategy := string(f.PhoneNull) == "1"
		filtered = listhelper.Intersection(filtered, s.phoneNullIndex[filterStrategy])
	}
	if len(f.PremiumNull) != 0 {
		filterStrategy := string(f.PremiumNull[0]) == "1"
		filtered = listhelper.Intersection(filtered, s.premiumWasIndex[filterStrategy])
	}
	if len(f.PremiumNow) != 0 {
		filtered = listhelper.Intersection(filtered, s.premiumActiveIndex[false])
	}

	if len(f.InterestsContains) > 0 {
		contains := containsFilter(f.InterestsContains, s.interestsIndex)
		filtered = listhelper.Intersection(filtered, contains)
	}
	if len(f.LikesContains) > 0 {
		contains := containsLikesFilter(f.LikesContains, s.likesIndex)
		filtered = listhelper.Intersection(filtered, contains)
	}

	if len(f.FnameAny) > 0 {
		any := anyFilter(f.FnameAny, s.fnameIndex)
		filtered = listhelper.Intersection(filtered, any)
	}
	if len(f.CityAny) > 0 {
		any := anyFilter(f.CityAny, s.cityIndex)
		filtered = listhelper.Intersection(filtered, any)
	}
	if len(f.InterestsAny) > 0 {
		any := anyFilter(f.InterestsAny, s.interestsIndex)
		filtered = listhelper.Intersection(filtered, any)
	}
	if len(f.LikesContainsAny) > 0 {
		any := anyFilterLikes(f.LikesContainsAny, s.likesIndex)
		filtered = listhelper.Intersection(filtered, any)
	}

	if len(f.EmailLt) != 0 || len(f.EmailGt) != 0 {
		var emailScanResult []int
		emailLt := string(f.EmailLt)
		emailGt := string(f.EmailGt)
		s.emailUniqueIndex.Range(func(email, id interface{}) bool {
			em := email.(string)
			if len(f.EmailLt) != 0 && emailLt <= em {
				return true
			}
			if len(f.EmailGt) != 0 && emailGt >= em {
				return true
			}
			emailScanResult = append(emailScanResult, id.(int))
			return true
		})
		sort.Ints(emailScanResult)
		filtered = listhelper.Intersection(filtered, emailScanResult)
	}

	if f.BirthLt != 0 || f.BirthGt != 0 {
		var birthScanResult []int
		for birth, ids := range s.bithIndex {
			if f.BirthLt != 0 && f.BirthLt <= birth {
				continue
			}
			if f.BirthGt != 0 && f.BirthGt >= birth {
				continue
			}
			birthScanResult = append(birthScanResult, ids...)
		}
		sort.Ints(birthScanResult)
		filtered = listhelper.Intersection(filtered, birthScanResult)
	}

	return listhelper.RemoveDuplicates(filtered) // likes index contain duplicates. TODO resolve
}

func (s *Storage) Filter(f Filter, limit int) []structures.Account {
	filtered := s.FilterNoDataFetch(f)

	sort.Sort(sort.Reverse(sort.IntSlice(filtered)))

	l := len(filtered)
	if limit > l {
		limit = l
	}

	return s.GetByList(filtered[:limit])
}

func containsFilter(containsFilter []string, index map[string][]int) []int {
	var contains []int
	initFlag := false
	for _, v := range containsFilter {
		contains = appendOrIntersect(&initFlag, contains, index[v])
	}
	sort.Ints(contains)
	return contains
}
func containsLikesFilter(containsFilter []int, index [][]int) []int {
	var contains []int
	initFlag := false
	for _, v := range containsFilter {
		contains = appendOrIntersect(&initFlag, contains, index[v])
	}
	sort.Ints(contains)
	return contains
}

func anyFilter(anyValue []string, index map[string][]int) []int {
	var any []int
	for _, value := range anyValue {
		any = append(any, index[value]...)
	}
	any = listhelper.RemoveDuplicates(any)
	sort.Ints(any)
	return any
}
func anyFilterLikes(anyValue []int, index [][]int) []int {
	var any []int
	for _, value := range anyValue {
		any = append(any, index[value]...)
	}
	any = listhelper.RemoveDuplicates(any)
	sort.Ints(any)
	return any
}

func nullFilter(value []byte, filtered []int, index map[string][]int) []int {
	if string(value) == "1" {
		return listhelper.Intersection(filtered, index[""])
	}
	var notNull []int
	for name, ids := range index {
		if name != "" {
			notNull = append(notNull, ids...)
		}
	}
	sort.Ints(notNull)
	return listhelper.Intersection(filtered, notNull)
}

func appendOrIntersect(intersect *bool, a, b []int) []int {
	if *intersect {
		return listhelper.Intersection(a, b)
	}
	*intersect = true
	return append(a, b...)
}
