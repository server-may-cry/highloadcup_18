package validators_test

import (
	"testing"

	"github.com/server-may-cry/highloadcup_18/structures"
	"github.com/server-may-cry/highloadcup_18/validators"
)

func TestIsAccountValid(t *testing.T) {
	var accountsValidatorTestData = []struct {
		in  structures.Account
		out bool
	}{
		{
			structures.Account{
				Email:  "asd@asd.asd",
				Sex:    "f",
				Status: structures.StatusFree,
				Birth:  1,
				Joined: 1,
			},
			true,
		},
		{
			structures.Account{
				Email:  "asd@asd.asd",
				Sex:    "f",
				Status: "asdasdas",
				Birth:  1,
				Joined: 1,
			},
			false,
		},
		{
			structures.Account{
				Email:  "ittacaxet",
				Sex:    "f",
				Status: structures.StatusHold,
				Birth:  1,
				Joined: 1,
			},
			false,
		},
	}
	for _, testData := range accountsValidatorTestData {
		out := validators.IsAccountValid(testData.in)
		if out != testData.out {
			t.Errorf("got %v, want %v", out, testData.out)
		}
	}
}
