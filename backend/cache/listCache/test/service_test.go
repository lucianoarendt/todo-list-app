package listCache_test

import (
	"testing"

	"github.com/rafaelmf3/todo-list/cache/listCache"
	"github.com/rafaelmf3/todo-list/models"
	"gorm.io/gorm"
)

type listCacheServiceTester struct {
	cacheHandler listCache.ListCacheService
	logAs        string
}

//READALL TODO

var (
	id     = 1
	userID = 2

	simbles = models.Symbol{
		Model:  gorm.Model{ID: 2},
		Symbol: "simbles",
		ListID: uint(id),
	}

	testSymbol = models.Symbol{
		Model:  gorm.Model{ID: 1},
		Symbol: "simbolo O_o",
		ListID: uint(id),
	}

	testList = models.List{
		Model:     gorm.Model{ID: uint(id)},
		Name:      "nome",
		IsDefault: false,
		UserID:    uint(userID),
		Symbols: []models.Symbol{
			simbles,
		},
	}

	updatedList = models.List{
		Model:     gorm.Model{ID: uint(id)},
		Name:      "nome atualizado",
		IsDefault: false,
		UserID:    uint(userID),
		Symbols:   []models.Symbol{},
	}
)

func (st *listCacheServiceTester) benchmarkListCacheServiceCreate(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		st.cacheHandler.CreateOnCache(testList)
		fromCache, _ := st.cacheHandler.TryReadingFromCache(userID, id,
			emptyElseGetDataFrom,
		)
		if !testList.Equals(*fromCache) {
			b.Errorf("%sCreate want=%v, got=%v", st.logAs, testList, fromCache)
		}
	}
}

func (st *listCacheServiceTester) benchmarkListCacheServiceReadNoCache(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		fromElse, _ := st.cacheHandler.TryReadingFromCache(userID, id,
			func() (models.List, error) {
				return testList, nil
			})

		if !testList.Equals(*fromElse) {
			b.Errorf("%sReadNoCache want=%v, got=%v", st.logAs, testList, fromElse)
		}
	}
}

func (st *listCacheServiceTester) benchmarkListCacheServiceReadNoCacheAndWithCache(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		fromElse, _ := st.cacheHandler.TryReadingFromCache(userID, id,
			func() (models.List, error) {
				return testList, nil
			})

		fromCache, _ := st.cacheHandler.TryReadingFromCache(userID, id,
			emptyElseGetDataFrom,
		)

		if !compareLists(testList, *fromElse, *fromCache) {
			b.Errorf("%sReadNoCacheAndWithCache want=%v, got=%v,%v", st.logAs, testList, fromElse, fromCache)
		}
	}
}

func (st *listCacheServiceTester) benchmarkListCacheServiceReadAllNoCache(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fromElse, _ := st.cacheHandler.TryReadingAllFromCache(userID,
			func() ([]models.List, error) {
				return []models.List{testList}, nil
			})

		if len(fromElse) != 1 || !fromElse[0].Equals(testList) {
			b.Errorf("%sReadAllNoCache want=%v, got=%v", st.logAs, []models.List{testList}, fromElse)
		}
	}
}

func (st *listCacheServiceTester) benchmarkListCacheServiceReadAllNoCacheAndWithCache(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fromElse, _ := st.cacheHandler.TryReadingAllFromCache(userID,
			func() ([]models.List, error) {
				return []models.List{testList}, nil
			})

		fromCache, _ := st.cacheHandler.TryReadingAllFromCache(userID,
			func() ([]models.List, error) {
				return []models.List{}, nil
			})

		if len(fromElse) != 1 || !fromElse[0].Equals(testList) {
			b.Errorf("%sReadAllNoCache want=%v, got=%v", st.logAs, []models.List{testList}, fromElse)
		}

		if len(fromCache) != 1 || !fromCache[0].Equals(testList) {
			b.Errorf("%sReadAllNoCache want=%v, got=%v", st.logAs, []models.List{testList}, fromCache)
		}
	}
}

func (st *listCacheServiceTester) benchmarkListCacheServiceUpdateNoCache(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		noCache, _ := st.cacheHandler.UpdateOnCache(testList)
		fromCache, _ := st.cacheHandler.TryReadingFromCache(userID, id,
			emptyElseGetDataFrom,
		)
		if !compareLists(testList, *noCache, *fromCache) {
			b.Errorf("%sUpdateNoCache want=%v, got=%v,%v", st.logAs, testList, noCache, fromCache)
		}
	}
}

func (st *listCacheServiceTester) benchmarkListCacheServiceUpdateWithCache(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		st.cacheHandler.CreateOnCache(testList)
		noCache, _ := st.cacheHandler.UpdateOnCache(updatedList)
		fromCache, _ := st.cacheHandler.TryReadingFromCache(userID, id,
			emptyElseGetDataFrom,
		)
		if !compareLists(updatedList, *noCache, *fromCache) || updatedList.Equals(testList) {
			b.Errorf("%sUpdateWithCache want=%v, got=%v,%v", st.logAs, testList, noCache, fromCache)
		}
	}
}

func (st *listCacheServiceTester) benchmarkListCacheServiceDelete(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		st.cacheHandler.CreateOnCache(testList)
		st.cacheHandler.DeleteOnCache(testList)
		cacheList, _ := st.cacheHandler.TryReadingFromCache(userID, id, emptyElseGetDataFrom)

		if compareLists(*cacheList, testList) {
			b.Errorf("%sDelete want=%v, got=%v", st.logAs, testList, cacheList)
		}
	}
}

func (st *listCacheServiceTester) benchmarkListCacheServiceCreateSymbol(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		st.cacheHandler.CreateOnCache(testList)
		st.cacheHandler.CreateSymbolOnCache(userID, testSymbol)
		cacheList, _ := st.cacheHandler.TryReadingFromCache(userID, id, emptyElseGetDataFrom)
		if testList.Equals(*cacheList) || !cacheList.Contains(testSymbol) {
			b.Errorf("%sCreateSymbol\n want=%v,\n got=%v", st.logAs, testList, cacheList)
		}
	}
}

func (st *listCacheServiceTester) benchmarkListCacheServiceDeleteSymbol(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		st.cacheHandler.CreateOnCache(testList)
		st.cacheHandler.DeleteSymbolOnCache(userID, simbles)
		cacheList, _ := st.cacheHandler.TryReadingFromCache(userID, id, emptyElseGetDataFrom)
		if testList.Equals(*cacheList) || cacheList.Contains(simbles) {
			b.Errorf("%sDeleteSymbol\n want=%v,\n got=%v", st.logAs, testList, cacheList)
		}
	}
}

func (st *listCacheServiceTester) benchMarkListCacheServiceReadAllDefaultNoCache(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		fromElse, _ := st.cacheHandler.ReadAllDefaultFromCache(
			func() ([]models.List, error) {
				return []models.List{testList}, nil
			})

		if len(fromElse) != 1 || !fromElse[0].Equals(testList) {
			b.Errorf("%sReadAllNoCache want=%v, got=%v", st.logAs, []models.List{testList}, fromElse)
		}
	}
}

func (st *listCacheServiceTester) benchMarkListCacheServiceReadAllDefaultNoCacheAndWithCache(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		fromElse, _ := st.cacheHandler.ReadAllDefaultFromCache(
			func() ([]models.List, error) {
				return []models.List{testList}, nil
			})
		fromCache, _ := st.cacheHandler.ReadAllDefaultFromCache(
			func() ([]models.List, error) {
				return []models.List{}, nil
			})

		if len(fromElse) != 1 || !fromElse[0].Equals(testList) ||
			len(fromCache) != 1 || !fromCache[0].Equals(testList) {
			b.Errorf("%sReadAllNoCache want=%v, got=%v,%v", st.logAs, []models.List{testList}, fromElse, fromCache)
		}
	}
}

func emptyElseGetDataFrom() (models.List, error) {
	return models.List{}, nil
}

func compareLists(lists ...models.List) bool {
	for i := 0; i < len(lists)-1; i++ {
		if !lists[i].Equals(lists[i+1]) {
			return false
		}
	}

	return true
}
