package functions

import (
	"math"

	"github.com/server-may-cry/highloadcup_18/accounthelper"
	"github.com/server-may-cry/highloadcup_18/structures"
)

// max int64 9,223,372,036,854,775,808
// max int32             2,147,483,647
// max age diff - int32

var (
	statusGroupWeight   = int64(math.Pow(2, 48)) / 4 // 70,368,744,177,664
	interestMatchWeight = int64(math.Pow(2, 32))     //      4,294,967,296
)

func CalculateCompatibility(me, to structures.Account) int64 {
	var total int64

	if accounthelper.IsPremium(to.Premium) {
		total += statusGroupWeight * 3
	}
	if to.Status == structures.StatusFree {
		total += statusGroupWeight * 2
	} else if to.Status == structures.StatusAllComplex {
		total += statusGroupWeight
	}
	toInterestsMap := make(map[string]struct{}, len(to.Interests))
	for _, item := range to.Interests {
		toInterestsMap[item] = struct{}{}
	}

	for _, meInterest := range me.Interests {
		if _, ok := toInterestsMap[meInterest]; ok {
			total += interestMatchWeight
		}
	}

	var ageDiff = me.Birth - to.Birth
	if ageDiff < 0 {
		ageDiff = -ageDiff
	}
	total -= int64(ageDiff)

	return total
}
