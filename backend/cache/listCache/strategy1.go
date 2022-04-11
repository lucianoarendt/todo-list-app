package listCache

import (
	"fmt"
	"strings"

	"github.com/patrickmn/go-cache"
	"github.com/rafaelmf3/todo-list/database"
	"github.com/rafaelmf3/todo-list/models"
)

type cacheStrategy1 struct {
	cache *cache.Cache
}

func NewCacheStrategy1(cache *cache.Cache) ListCacheService {
	return &cacheStrategy1{
		cache: cache,
	}
}

func (c *cacheStrategy1) CreateOnCache(list models.List) error {

	_, listKey := mountKeys(int(list.UserID), int(list.ID))

	c.cache.SetDefault(listKey, list)

	return nil
}

func (c *cacheStrategy1) TryReadingFromCache(userID int, listID int, elseGetDataFrom func() (models.List, error)) (*models.List, error) {
	_, cacheKey := mountKeys(userID, listID)
	listCache, existsOnCache := c.cache.Get(cacheKey)

	var list models.List
	if !existsOnCache {
		var err error
		list, err = elseGetDataFrom()
		if err != nil {
			return nil, err
		}

		c.cache.SetDefault(cacheKey, list)
	} else {
		list = listCache.(models.List)
	}

	return &list, nil
}

func (c *cacheStrategy1) UpdateOnCache(list models.List) (*models.List, error) {

	_, cacheKey := mountKeys(int(list.UserID), int(list.ID))
	if list.Symbols == nil {
		cacheList, existsOnCache := c.cache.Get(cacheKey)
		if existsOnCache {
			cacheAsList := cacheList.(models.List)
			list.Symbols = cacheAsList.Symbols
		} else {
			list.PopulateWithSymbols(database.DB)
		}
	}
	//-----------------

	c.cache.SetDefault(cacheKey, list)

	return &list, nil
}

func (c *cacheStrategy1) DeleteOnCache(list models.List) error {
	_, cacheKey := mountKeys(int(list.UserID), int(list.ID))
	c.cache.Delete(cacheKey)
	return nil
}

func (c *cacheStrategy1) TryReadingAllFromCache(userID int, elseGetDataFrom func() ([]models.List, error)) ([]models.List, error) {
	userKey, _ := mountKeys(userID, -1)
	_, existsOnCache := c.cache.Get(userKey)

	var lists []models.List
	if !existsOnCache {
		var err error
		lists, err = elseGetDataFrom()
		if err != nil {
			return nil, err
		}

		for i := range lists {
			c.cache.SetDefault(fmt.Sprintf("%s_%d", userKey, lists[i].ID), lists[i])
		}
		c.cache.SetDefault(userKey, true)
	} else {
		lists, _ = c.getCacheLists(userKey)
	}
	return lists, nil
}

func (c *cacheStrategy1) ReadAllDefaultFromCache(elseGetDataFrom func() ([]models.List, error)) ([]models.List, error) {
	var lists []models.List

	listsCache, existsOnCache := c.cache.Get("default")

	if !existsOnCache {
		lists, _ = elseGetDataFrom()
		c.cache.SetDefault("default", lists)

	} else {
		lists = listsCache.([]models.List)
	}

	return lists, nil
}

func (c *cacheStrategy1) CreateSymbolOnCache(userID int, symbol models.Symbol) error {

	//Cache Handling
	_, cacheKey := mountKeys(userID, int(symbol.ListID))

	listCache, existsOnCache := c.cache.Get(cacheKey)

	if existsOnCache {
		cacheAsList := listCache.(models.List)
		SymbolsCopy := make([]models.Symbol, len(cacheAsList.Symbols))
		copy(SymbolsCopy, cacheAsList.Symbols)

		cacheAsList.Symbols = append(SymbolsCopy, symbol)
		c.cache.SetDefault(cacheKey, cacheAsList)
	}

	//--------------------
	return nil
}

func (c *cacheStrategy1) DeleteSymbolOnCache(userID int, symbol models.Symbol) error {
	//Cache handling
	_, cacheKey := mountKeys(userID, int(symbol.ListID))

	listCache, existsOnCache := c.cache.Get(cacheKey)

	if existsOnCache {
		cacheAsList := listCache.(models.List)
		SymbolsCopy := make([]models.Symbol, len(cacheAsList.Symbols))
		copy(SymbolsCopy, cacheAsList.Symbols)

		//Delete a symbol
		for i, e := range SymbolsCopy {
			if e.ID == uint(symbol.ID) {
				SymbolsCopy[i] = SymbolsCopy[len(SymbolsCopy)-1]
				SymbolsCopy = SymbolsCopy[:len(SymbolsCopy)-1]
				break
			}
		}
		//-----------------
		cacheAsList.Symbols = SymbolsCopy

		c.cache.SetDefault(cacheKey, cacheAsList)
	}
	//-----------------

	return nil
}

func (c *cacheStrategy1) getCacheLists(userID string) ([]models.List, bool) {
	itemsMap := c.cache.Items()
	var lists []models.List

	for k := range itemsMap {
		if strings.Contains(k, fmt.Sprintf("%s_", userID)) {
			lists = append(lists, itemsMap[k].Object.(models.List))
		}
	}

	return lists, len(lists) > 0
}
