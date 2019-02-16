package httphelper

import (
	"bytes"

	"github.com/valyala/fasthttp"
)

const stripPefixLength = len("/accounts/")

func GetAccountID(path []byte) (int, error) {
	cuttedPath := path[stripPefixLength:]
	slashPos := bytes.Index(cuttedPath, []byte("/"))
	accountID := cuttedPath[:slashPos]
	return fasthttp.ParseUint(accountID)
}
