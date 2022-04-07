package listCache

import (
	"fmt"
	"strconv"

	"github.com/rafaelmf3/todo-list/models"
)

type ListCacheService interface {
	CreateOnCache(list models.List) error
	ReadFromCache(userID int, listID int, getDataFrom func() (models.List, error)) (*models.List, error) //todo
	UpdateOnCache(list models.List) (*models.List, error)
	DeleteOnCache(list models.List) error
	ReadAllFromCache(userID int, getDataFrom func() ([]models.List, error)) ([]models.List, error) //todo
	ReadAllDefaultFromCache() error
	CreateSymbolOnCache(userID int, symbol models.Symbol) error
	DeleteSymbolOnCache(userID int, symbol models.Symbol) error
}

func mountKeys(userID int, listID int) (string, string) {
	userKey := strconv.Itoa(userID)
	listKey := fmt.Sprintf("%d_%d", userID, listID)

	return userKey, listKey
}
