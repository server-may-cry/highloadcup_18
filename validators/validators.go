package validators

import (
	"strings"

	"github.com/server-may-cry/highloadcup_18/structures"
)

func IsAccountValid(a structures.Account) bool {
	if len(a.Email) > 100 {
		return false
	}

	if i := strings.Index(a.Email, "@"); i == -1 {
		return false
	}

	if len(a.Fname) > 50 {
		return false
	}

	if len(a.Sname) > 50 {
		return false
	}

	if len(a.Phone) > 16 {
		return false
	}

	if a.Sex != structures.SexF && a.Sex != structures.SexM {
		return false
	}

	if len(a.City) > 50 {
		return false
	}

	if a.Birth == 0 {
		return false
	}

	if a.Joined == 0 {
		return false
	}

	if a.Status != structures.StatusFree &&
		a.Status != structures.StatusHold &&
		a.Status != structures.StatusAllComplex {
		return false
	}

	for _, interest := range a.Interests {
		if len(interest) > 100 {
			return false
		}
	}

	return true
}
