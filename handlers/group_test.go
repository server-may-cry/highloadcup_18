package handlers_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/server-may-cry/highloadcup_18/db"
	"github.com/server-may-cry/highloadcup_18/handlers"
	"github.com/server-may-cry/highloadcup_18/initloader"
	"github.com/valyala/fasthttp"
)

func initBench() *db.Storage {
	storage := db.New()
	log.SetOutput(ioutil.Discard)
	initloader.Load(storage, "/tmp/data")
	storage.Reindex()
	return storage
}

func BenchmarkGroup(b *testing.B) {
	storage := initBench()
	h := handlers.GetGroupHandler(storage)

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.SetRequestURI("/accounts/group/?order=1&keys=country&birth=2002&sex=f&limit=45")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h(ctx)
		if ctx.Response.StatusCode() != http.StatusOK {
			b.Fatalf("Incorrect status code %d. Expected %d", ctx.Response.StatusCode(), http.StatusOK)
			b.FailNow()
		}
	}
}
