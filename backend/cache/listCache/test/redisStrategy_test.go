package listCache_test

import (
	"testing"

	"github.com/rafaelmf3/todo-list/cache/listCache"
)

var (
	stRedis = listCacheServiceTester{
		listCache.NewCacheRedisStrategy(),
		"ListCacheRedisStrategy",
	}
)

func BenchmarkListCacheRedisStrategyCreate(b *testing.B) {
	stRedis.benchmarkListCacheServiceCreate(b)
}

func BenchmarkListCacheRedisStrategyReadNoCache(b *testing.B) {
	stRedis.benchmarkListCacheServiceReadNoCache(b)
}

func BenchmarkListCacheRedisStrategyReadNoCacheAndWithCache(b *testing.B) {
	stRedis.benchmarkListCacheServiceReadNoCacheAndWithCache(b)
}

func BenchmarkListCacheRedisStrategyReadAllNoCache(b *testing.B) {
	stRedis.benchmarkListCacheServiceReadAllNoCache(b)
}

func BenchmarkListCacheRedisStrategyReadAllNoCacheAndWithCache(b *testing.B) {
	stRedis.benchmarkListCacheServiceReadAllNoCacheAndWithCache(b)
}

func BenchmarkListCacheRedisStrategyUpdateNoCache(b *testing.B) {
	stRedis.benchmarkListCacheServiceUpdateNoCache(b)
}

func BenchmarkListCacheRedisStrategyUpdateWithCache(b *testing.B) {
	stRedis.benchmarkListCacheServiceUpdateWithCache(b)
}

func BenchmarkListCacheRedisStrategyDelete(b *testing.B) {
	stRedis.benchmarkListCacheServiceDelete(b)
}

func BenchmarkListCacheRedisStrategyCreateSymbol(b *testing.B) {
	stRedis.benchmarkListCacheServiceCreateSymbol(b)
}

func BenchmarkListCacheRedisStrategyDeleteSymbol(b *testing.B) {
	stRedis.benchmarkListCacheServiceDeleteSymbol(b)
}
