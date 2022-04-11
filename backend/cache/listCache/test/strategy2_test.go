package listCache_test

import (
	"testing"

	"github.com/patrickmn/go-cache"
	"github.com/rafaelmf3/todo-list/cache/listCache"
)

var (
	stS2 = listCacheServiceTester{
		listCache.NewCacheStrategy2(cache.New(-1, -1)),
		"ListCacheStrategy2",
	}
)

func BenchmarkListCacheStrategy2Create(b *testing.B) {
	stS2.benchmarkListCacheServiceCreate(b)
}

func BenchmarkListCacheStrategy2ReadNoCache(b *testing.B) {
	stS2.benchmarkListCacheServiceReadNoCache(b)
}

func BenchmarkListCacheStrategy2ReadNoCacheAndWithCache(b *testing.B) {
	stS2.benchmarkListCacheServiceReadNoCacheAndWithCache(b)
}

func BenchmarkListCacheStrategy2ReadAllNoCache(b *testing.B) {
	stS2.benchmarkListCacheServiceReadAllNoCache(b)
}

func BenchmarkListCacheStrategy2ReadAllNoCacheAndWithCache(b *testing.B) {
	stS2.benchmarkListCacheServiceReadAllNoCacheAndWithCache(b)
}

func BenchmarkListCacheStrategy2UpdateNoCache(b *testing.B) {
	stS2.benchmarkListCacheServiceUpdateNoCache(b)
}

func BenchmarkListCacheStrategy2UpdateWithCache(b *testing.B) {
	stS2.benchmarkListCacheServiceUpdateWithCache(b)
}

func BenchmarkListCacheStrategy2Delete(b *testing.B) {
	stS2.benchmarkListCacheServiceDelete(b)
}

func BenchmarkListCacheStrategy2CreateSymbol(b *testing.B) {
	stS2.benchmarkListCacheServiceCreateSymbol(b)
}

func BenchmarkListCacheStrategy2DeleteSymbol(b *testing.B) {
	stS2.benchmarkListCacheServiceDeleteSymbol(b)
}

func BenchmarkListCacheStrategy2ReadAllDefaultNoCache(b *testing.B) {
	stS2.benchMarkListCacheServiceReadAllDefaultNoCache(b)
}

func BenchmarkListCacheStrategy2ReadAllDefaultNoCacheAndWithCache(b *testing.B) {
	stS2.benchMarkListCacheServiceReadAllDefaultNoCacheAndWithCache(b)
}
