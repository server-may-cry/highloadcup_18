package db_test

import (
	"io/ioutil"
	"log"
	"math/rand"
	"testing"

	"github.com/server-may-cry/highloadcup_18/db"
	"github.com/server-may-cry/highloadcup_18/initloader"
)

var (
	storage      *db.Storage
	accountCount int
)

func initBench() {
	storage = db.New()
	log.SetOutput(ioutil.Discard)
	accountCount = initloader.Load(storage, "/tmp/data")
	storage.Reindex()
}

func BenchmarkReindex(b *testing.B) {
	initBench()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		storage.Reindex()
	}
}

func BenchmarkGet(b *testing.B) {
	initBench()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		id := rand.Intn(accountCount + accountCount/2)
		storage.Get(id)
	}
}

func BenchmarkGroup(b *testing.B) {
	initBench()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		storage.Group(db.Filter{
			SexEq: []byte("f"),
		}, []string{"sex", "status", "city"}, 40, true)
	}
}

func BenchmarkFilterNoDataFetch(b *testing.B) {
	initBench()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		storage.FilterNoDataFetch(db.Filter{
			SexEq: []byte("f"),
		})
	}
}

func BenchmarkFilter(b *testing.B) {
	initBench()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		storage.Filter(db.Filter{
			SexEq: []byte("f"),
		}, 40)
	}
}

func BenchmarkArrayIndex(b *testing.B) {
	a := make([]int, b.N)
	size := 0

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a[size] = i
		size++
	}
}

func BenchmarkArrayAppend(b *testing.B) {
	a := make([]int, 0, b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a = append(a, i)
	}
}
