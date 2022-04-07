package listCache_test

import (
	"testing"

	"github.com/patrickmn/go-cache"
	"github.com/rafaelmf3/todo-list/cache/listCache"
	"github.com/rafaelmf3/todo-list/models"
)

var (
	cacheHandler listCache.ListCacheService = listCache.NewCacheStrategy1(cache.New(-1, -1))
)

func BenchmarkStrategy1Create(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		list := models.List{
			Name: "nome",
		}
		list.ID = 1
		list.UserID = 2
		cacheHandler.CreateOnCache(list)
	}
}
func BenchmarkStrategy1ReadNoCache(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		list := models.List{
			Name: "nome",
		}
		list.ID = 1
		list.UserID = 2
		fromCache, _ := cacheHandler.ReadFromCache(int(list.ID), int(list.UserID), func() (models.List, error) {
			return list, nil
		})

		if list.ID != fromCache.ID || list.UserID != fromCache.UserID || list.Name != fromCache.Name {
			b.Errorf("Strategy1ReadNoCache want=%v, got=%v", list, fromCache)
		}
	}
}

func BenchmarkStrategy1(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		list := models.List{
			Name: "nome",
		}
		list.ID = 1
		list.UserID = 2
		cacheHandler.CreateOnCache(list)
	}
}
