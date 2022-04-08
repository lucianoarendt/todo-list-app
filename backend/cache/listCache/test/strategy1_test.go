package listCache_test

import (
	"testing"

	"github.com/patrickmn/go-cache"
	"github.com/rafaelmf3/todo-list/cache/listCache"
)

var (
	stS1 = listCacheServiceTester{
		listCache.NewCacheStrategy1(cache.New(-1, -1)),
		"ListCacheStrategy1",
	}
)

func BenchmarkListCacheStrategy1Create(b *testing.B) {
	stS1.benchmarkListCacheServiceCreate(b)
}

func BenchmarkListCacheStrategy1ReadNoCache(b *testing.B) {
	stS1.benchmarkListCacheServiceReadNoCache(b)
}

func BenchmarkListCacheStrategy1ReadNoCacheAndWithCache(b *testing.B) {
	stS1.benchmarkListCacheServiceReadNoCacheAndWithCache(b)
}

func BenchmarkListCacheStrategy1ReadAllNoCache(b *testing.B) {
	stS1.benchmarkListCacheServiceReadAllNoCache(b)
}

func BenchmarkListCacheStrategy1ReadAllNoCacheAndWithCache(b *testing.B) {
	stS1.benchmarkListCacheServiceReadAllNoCacheAndWithCache(b)
}

func BenchmarkListCacheStrategy1UpdateNoCache(b *testing.B) {
	stS1.benchmarkListCacheServiceUpdateNoCache(b)
}

func BenchmarkListCacheStrategy1UpdateWithCache(b *testing.B) {
	stS1.benchmarkListCacheServiceUpdateWithCache(b)
}

func BenchmarkListCacheStrategy1Delete(b *testing.B) {
	stS1.benchmarkListCacheServiceDelete(b)
}

func BenchmarkListCacheStrategy1CreateSymbol(b *testing.B) {
	stS1.benchmarkListCacheServiceCreateSymbol(b)
}

func BenchmarkListCacheStrategy1DeleteSymbol(b *testing.B) {
	stS1.benchmarkListCacheServiceDeleteSymbol(b)
}
