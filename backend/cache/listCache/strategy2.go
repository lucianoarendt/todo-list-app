package listCache

import (
	"fmt"
	"strconv"

	"github.com/patrickmn/go-cache"
	"github.com/rafaelmf3/todo-list/database"
	"github.com/rafaelmf3/todo-list/models"
)

type cacheStrategy2 struct {
	cache *cache.Cache
}

func NewCacheStrategy2(cache *cache.Cache) ListCacheService {
	return &cacheStrategy2{
		cache: cache,
	}
}

func (c *cacheStrategy2) CreateOnCache(list models.List) error {
	userKey, listKey := mountKeys(int(list.UserID), int(list.ID))
	//Handles cache
	ids, userFound := c.cache.Get(userKey)
	if userFound {
		//At user level
		idsSlice := ids.([]string)
		idsSlice = append(idsSlice, strconv.Itoa(int(list.ID)))
		c.cache.SetDefault(userKey, idsSlice)
	}
	//At list level
	c.cache.SetDefault(listKey, list)
	//---------------

	return nil
}

func (c *cacheStrategy2) TryReadingFromCache(userID int, listID int, elseGetDataFrom func() (models.List, error)) (*models.List, error) {
	_, listKey := mountKeys(userID, listID)
	listCache, existsOnCache := c.cache.Get(listKey)

	var err error
	var list models.List
	if !existsOnCache {
		list, err = elseGetDataFrom()
		if err != nil {
			return nil, err
		}

		c.cache.SetDefault(listKey, list)
	} else {
		list = listCache.(models.List)
	}

	return &list, err
}

func (c *cacheStrategy2) UpdateOnCache(list models.List) (*models.List, error) {
	//Handle Cache
	_, listKey := mountKeys(int(list.UserID), int(list.ID))
	if list.Symbols == nil {
		cacheList, existsOnCache := c.cache.Get(listKey)
		if existsOnCache {
			cacheAsList := cacheList.(models.List)
			list.Symbols = cacheAsList.Symbols
		} else {
			list.PopulateWithSymbols(database.DB)
		}
	}

	c.cache.SetDefault(listKey, list)
	//-----------------------

	return &list, nil
}

func (c *cacheStrategy2) DeleteOnCache(list models.List) error {
	userKey, listKey := mountKeys(int(list.UserID), int(list.ID))
	//Handles Cache
	ids, userFound := c.cache.Get(userKey)
	if userFound {
		//Delete id from user cache
		idsSlice := ids.([]string)
		for i, e := range idsSlice {
			if e == strconv.Itoa(int(list.UserID)) {
				idsSlice[i] = idsSlice[len(idsSlice)-1]
				idsSlice = idsSlice[:len(idsSlice)-1]
				break
			}
		}
		c.cache.SetDefault(userKey, idsSlice)
		//-----------------------
	}
	c.cache.Delete(listKey)
	//-----------------
	return nil
}

func (c *cacheStrategy2) TryReadingAllFromCache(userID int, elseGetDataFrom func() ([]models.List, error)) ([]models.List, error) {

	userKey, _ := mountKeys(userID, -1)

	lists, existsOnCache := c.getCacheLists(userKey)

	if !existsOnCache {
		var err error
		lists, err = elseGetDataFrom()
		if err != nil {
			return nil, err
		}

		ids := make([]string, len(lists))
		for i := range lists {
			ids[i] = strconv.Itoa(int(lists[i].ID))
			c.cache.SetDefault(fmt.Sprintf("%s_%d", userKey, lists[i].ID), lists[i])
		}
		c.cache.SetDefault(userKey, ids)
	}

	return lists, nil
}

func (c *cacheStrategy2) ReadAllDefaultFromCache() ([]models.List, error) {
	return nil, nil
}

func (c *cacheStrategy2) CreateSymbolOnCache(userID int, symbol models.Symbol) error {
	//Cache Handling
	_, listKey := mountKeys(userID, int(symbol.ListID))

	listCache, existsOnCache := c.cache.Get(listKey)

	if existsOnCache {
		cacheAsList := listCache.(models.List)
		SymbolsCopy := make([]models.Symbol, len(cacheAsList.Symbols))
		copy(SymbolsCopy, cacheAsList.Symbols)

		cacheAsList.Symbols = append(SymbolsCopy, symbol)
		c.cache.SetDefault(listKey, cacheAsList)
	}
	//--------------------

	return nil
}

func (c *cacheStrategy2) DeleteSymbolOnCache(userID int, symbol models.Symbol) error {
	_, listKey := mountKeys(userID, int(symbol.ListID))

	//Cache handling

	listCache, existsOnCache := c.cache.Get(listKey)

	if existsOnCache {
		cacheAsList := listCache.(models.List)
		SymbolsCopy := make([]models.Symbol, len(cacheAsList.Symbols))
		copy(SymbolsCopy, cacheAsList.Symbols)

		//Delete a symbol
		for i, e := range SymbolsCopy {
			if e.ID == symbol.ID {
				SymbolsCopy[i] = SymbolsCopy[len(SymbolsCopy)-1]
				SymbolsCopy = SymbolsCopy[:len(SymbolsCopy)-1]
			}
		}
		//-----------------
		cacheAsList.Symbols = SymbolsCopy

		c.cache.SetDefault(listKey, cacheAsList)
	}
	//-----------------
	return nil
}

func (c *cacheStrategy2) getCacheLists(userID string) ([]models.List, bool) {
	var lists []models.List

	ids, foundUser := c.cache.Get(userID)
	if foundUser {
		idsAsSlice := ids.([]string)

		for _, e := range idsAsSlice {
			listCache, _ := c.cache.Get(userID + "_" + e)
			listCacheAsList := listCache.(models.List)

			lists = append(lists, listCacheAsList)
		}
	}

	return lists, foundUser
}
