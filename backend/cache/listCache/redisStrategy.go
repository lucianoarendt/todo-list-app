package listCache

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/rafaelmf3/todo-list/database"
	"github.com/rafaelmf3/todo-list/models"
)

type cacheRedisStrategy struct {
	cache *redis.Client
}

func NewCacheRedisStrategy() ListCacheService {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		DB:       0,
	})

	return &cacheRedisStrategy{
		cache: client,
	}
}

func (c *cacheRedisStrategy) CreateOnCache(list models.List) error {
	_, listKey := mountKeys(int(list.UserID), int(list.ID))

	c.cache.Set(listKey, list, 0)

	return nil
}

func (c *cacheRedisStrategy) TryReadingFromCache(userID int, listID int, elseGetDataFrom func() (models.List, error)) (*models.List, error) {
	_, cacheKey := mountKeys(userID, listID)

	json, err := c.cache.Get(cacheKey).Result()

	var list models.List
	if err == redis.Nil {
		var err error
		list, err = elseGetDataFrom()
		if err != nil {
			return nil, err
		}

		c.cache.Set(cacheKey, list, 0)
	} else {
		list.Unmarshal(json)
	}
	return &list, nil
}

func (c *cacheRedisStrategy) UpdateOnCache(list models.List) (*models.List, error) {
	_, cacheKey := mountKeys(int(list.UserID), int(list.ID))
	if list.Symbols == nil {
		jsonList, err := c.cache.Get(cacheKey).Result()
		if err == nil {
			var cacheAsList models.List
			cacheAsList.Unmarshal(jsonList)
			list.Symbols = cacheAsList.Symbols
		} else {
			list.PopulateWithSymbols(database.DB)
		}
	}
	//-----------------

	c.cache.Set(cacheKey, list, 0)

	return &list, nil
}

func (c *cacheRedisStrategy) DeleteOnCache(list models.List) error {
	_, cacheKey := mountKeys(int(list.UserID), int(list.ID))

	return c.cache.Del(cacheKey).Err()
}

func (c *cacheRedisStrategy) TryReadingAllFromCache(userID int, elseGetDataFrom func() ([]models.List, error)) ([]models.List, error) {
	userKey, _ := mountKeys(userID, -1)
	_, err := c.cache.Get(userKey).Result()

	var pairs []interface{}

	var lists []models.List
	if err == redis.Nil {
		lists, err = elseGetDataFrom()
		if err != nil {
			return nil, err
		}

		for i := range lists {
			lists[i].PopulateWithSymbols(database.DB)
			pairs = append(pairs, fmt.Sprintf("%s_%d", userKey, lists[i].ID), lists[i])
		}
		pairs = append(pairs, userKey, true)
		c.cache.MSet(pairs...)
	} else {
		lists, _ = c.getCacheLists(userKey)
	}

	return lists, nil
}

func (c *cacheRedisStrategy) ReadAllDefaultFromCache() ([]models.List, error) {
	return nil, nil
}

func (c *cacheRedisStrategy) CreateSymbolOnCache(userID int, symbol models.Symbol) error {
	//Cache Handling
	_, cacheKey := mountKeys(userID, int(symbol.ListID))

	listJson, err := c.cache.Get(cacheKey).Result()

	if err == nil {
		var cacheAsList models.List
		cacheAsList.Unmarshal(listJson)

		SymbolsCopy := make([]models.Symbol, len(cacheAsList.Symbols))
		copy(SymbolsCopy, cacheAsList.Symbols)

		cacheAsList.Symbols = append(SymbolsCopy, symbol)
		c.cache.Set(cacheKey, cacheAsList, 0)
	}

	//--------------------
	return nil
}

func (c *cacheRedisStrategy) DeleteSymbolOnCache(userID int, symbol models.Symbol) error {
	//Cache handling
	_, cacheKey := mountKeys(userID, int(symbol.ListID))

	listJson, err := c.cache.Get(cacheKey).Result()

	if err == nil {
		var cacheAsList models.List
		cacheAsList.Unmarshal(listJson)

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

		c.cache.Set(cacheKey, cacheAsList, 0)
	}
	//-----------------
	return nil
}

func (c *cacheRedisStrategy) getCacheLists(userID string) ([]models.List, bool) {
	pattern := fmt.Sprintf("%s_*", userID)
	keys, _ := c.cache.Keys(pattern).Result()

	listsJson, _ := c.cache.MGet(keys...).Result()
	lists := make([]models.List, len(listsJson))
	for i := range listsJson {
		lists[i].Unmarshal(listsJson[i].(string))
	}

	return lists, len(lists) > 0
}
