package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"time"

	"github.com/server-may-cry/highloadcup_18/db"
	"github.com/server-may-cry/highloadcup_18/handlers"
	"github.com/server-may-cry/highloadcup_18/initloader"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/pprofhandler"
)

type NullLogger struct{}

func (NullLogger) Printf(format string, args ...interface{}) {}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	log.Println("start")

	path := flag.String("path", "/tmp/data", "")
	memprofile := flag.String("memprofile", "", "write memory profile to this file")
	port := flag.String("port", ":80", "HTTP port to listen")
	flag.Parse()

	storage := db.New()
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	log.Printf("storage initialized %d", stat.Alloc/1024/1024)
	n := time.Now()
	initloader.Load(storage, *path)
	runtime.GC()
	runtime.ReadMemStats(&stat)
	log.Printf("all json readed. Momery used %d", stat.Alloc/1024/1024)
	storage.Debug()
	d := time.Since(n)
	log.Println("initialized")
	log.Println(d)
	storage.Reindex()
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	filterHandler := handlers.GetFilterHandler(storage)
	groupHandler := handlers.GetGroupHandler(storage)
	recommendHandler := handlers.GetRecomendHandler(storage)
	suggestHandler := handlers.GetSuggestHandler(storage)
	updateHandler := handlers.GetUpdateAccountHandler(storage)
	newHandler := handlers.GetNewAccountHandler(storage)
	likesHandler := handlers.GetAddLikesHandler(storage)

	var pprofPrefix = []byte("/debug/pprof/")
	var (
		filter    = []byte("lter") // /accounts/filter/
		group     = []byte("roup") // /accounts/group/
		recommend = []byte("mend") // /accounts/{id}/recommend/
		suggest   = []byte("gest") // /accounts/{id}/suggest/
		new       = []byte("/new") // /accounts/new/
		likes     = []byte("ikes") // /accounts/likes/
		// /accounts/{id}/
	)
	m := func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered '%s':\n%s", r, string(debug.Stack()))
			}
		}()
		p := ctx.Path()
		if bytes.HasPrefix(p, pprofPrefix) {
			pprofhandler.PprofHandler(ctx)
			return
		}

		l := len(p)
		sliceForCompare := p[l-5 : l-1]
		if ctx.IsGet() {
			if isEqualBytes(sliceForCompare, filter) {
				filterHandler(ctx)
			} else if isEqualBytes(sliceForCompare, group) {
				groupHandler(ctx)
			} else if isEqualBytes(sliceForCompare, recommend) {
				recommendHandler(ctx)
			} else if isEqualBytes(sliceForCompare, suggest) {
				suggestHandler(ctx)
			} else {
				handlers.NotFoundHandler(ctx)
			}
		} else {
			if isEqualBytes(sliceForCompare, likes) {
				likesHandler(ctx)
			} else if isEqualBytes(sliceForCompare, new) {
				newHandler(ctx)
			} else {
				updateHandler(ctx)
			}
		}
	}

	//go stillAlive()
	log.Println("ready")
	s := &fasthttp.Server{
		Handler:                       m,
		DisableHeaderNamesNormalizing: true,
		ReduceMemoryUsage:             true,
		// NoDefaultServerHeader:         true,
		// NoDefaultContentType:          true,
		// Logger: NullLogger{},
	}
	err := s.ListenAndServe(*port)
	if err != nil {
		log.Fatal(err)
	}
}

func stillAlive() {
	timer := time.NewTicker(5 * time.Second)
	for {
		<-timer.C
		log.Println("Still alive")
	}
}

func isEqualBytes(a, b []byte) bool {
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
