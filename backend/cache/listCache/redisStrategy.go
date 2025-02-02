package listCache

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/rafaelmf3/todo-list/database"
	"github.com/rafaelmf3/todo-list/models"
)

type cacheRedisStrategy struct {
	cache *redis.Client
}

func NewCacheRedisStrategy() ListCacheService {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", os.Getenv("REDIS_ADDRESS")),
		Password: os.Getenv("REDIS_PASS"),
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
	} else if err == nil {
		list.Unmarshal(json)
	} else {
		return nil, err
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
		} else if err == redis.Nil {
			list.PopulateWithSymbols(database.DB)
		} else {
			return nil, err
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
			pairs = append(pairs, fmt.Sprintf("%s_%d", userKey, lists[i].ID), lists[i])
		}
		pairs = append(pairs, userKey, true)
		c.cache.MSet(pairs...)
	} else if err == nil {
		lists, _ = c.getCacheLists(userKey)
	} else {
		return nil, err
	}

	return lists, nil
}

func (c *cacheRedisStrategy) ReadAllDefaultFromCache(elseGetDataFrom func() ([]models.List, error)) ([]models.List, error) {
	jsonList, err := c.cache.Get(defaultKey).Result()

	var lists []models.List
	if err == redis.Nil {
		var err error
		lists, err = elseGetDataFrom()
		if err != nil {
			return nil, err
		}

		marshaledLists, _ := json.Marshal(lists)

		err = c.cache.Set(defaultKey, marshaledLists, 0).Err()
		if err != nil {
			return nil, err
		}
	} else if err == nil {
		json.Unmarshal([]byte(jsonList), &lists)
	} else {
		return nil, err
	}

	return lists, nil
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
