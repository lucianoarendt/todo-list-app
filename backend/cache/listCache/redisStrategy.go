package listCache

import (
	"github.com/go-redis/redis"
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

	return nil
}

func (c *cacheRedisStrategy) ReadFromCache(userID int, listID int, getDataFrom func() (models.List, error)) (*models.List, error) {

	return nil, nil
}

func (c *cacheRedisStrategy) UpdateOnCache(list models.List) (*models.List, error) {

	return nil, nil
}

func (c *cacheRedisStrategy) DeleteOnCache(list models.List) error {

	return nil
}

func (c *cacheRedisStrategy) ReadAllFromCache(userID int, getDataFrom func() ([]models.List, error)) ([]models.List, error) {

	return nil, nil
}

func (c *cacheRedisStrategy) ReadAllDefaultFromCache() error {
	return nil
}

func (c *cacheRedisStrategy) CreateSymbolOnCache(userID int, symbol models.Symbol) error {

	return nil
}

func (c *cacheRedisStrategy) DeleteSymbolOnCache(userID int, symbol models.Symbol) error {

	return nil
}
