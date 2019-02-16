package httphelper_test

import (
	"testing"

	"github.com/server-may-cry/highloadcup_18/httphelper"
)

func TestGetAccountID(t *testing.T) {
	var pathToIntTestData = []struct {
		in     []byte
		result int
	}{
		{
			[]byte("/accounts/1706/suggest/"),
			1706,
		},
		{
			[]byte("/accounts/314941/recommend/"),
			314941,
		},
		{
			[]byte("/accounts/123/"),
			123,
		},
		{
			[]byte("/accounts/asdadasd/"),
			-1,
		},
		{
			[]byte("/accounts/curlulilevtile/cyulgedeer/"),
			-1,
		},
		{
			[]byte("/accounts/hasdeceralitdos/"),
			-1,
		},
	}
	for _, data := range pathToIntTestData {
		result, _ := httphelper.GetAccountID(data.in)
		if result != data.result {
			t.Errorf("got %d, want %d", result, data.result)
		}
	}

}
