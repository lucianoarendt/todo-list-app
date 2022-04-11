package lists

import (
	"strconv"
	"strings"

	"github.com/patrickmn/go-cache"
	"github.com/rafaelmf3/todo-list/cache/listCache"
	"github.com/rafaelmf3/todo-list/database"
	"github.com/rafaelmf3/todo-list/middleware"
	"github.com/rafaelmf3/todo-list/models"

	"github.com/gofiber/fiber/v2"
)

type ListService interface {
	Create(c *fiber.Ctx) error
	Read(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	ReadAll(c *fiber.Ctx) error
	CreateSymbol(c *fiber.Ctx) error
	DeleteSymbol(c *fiber.Ctx) error
}

type listServiceImpl struct {
	secret       string
	cache        cache.Cache
	cacheHandler listCache.ListCacheService
}

func NewListService(secret string) ListService {
	l := &listServiceImpl{
		secret: secret,
		cache:  *cache.New(-1, -1),
	}
	l.cacheHandler = listCache.NewCacheStrategy1(&l.cache)
	return l
}

const MaxListsAmount = 10
const MaxSimbolsAmount = 50

func (l *listServiceImpl) Create(c *fiber.Ctx) error {
	//Authenticates
	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	userID, _ := strconv.Atoi(claims.Issuer)
	//---------------

	//mounts body
	var bodyData models.List
	if err := c.BodyParser(&bodyData); err != nil {
		return err
	}
	bodyData.UserID = uint(userID)
	//--------------

	if err := bodyData.CreateList(database.DB); err != nil {
		return l.handleListError(err, c)
	}

	//Handles cache
	l.cacheHandler.CreateOnCache(bodyData)

	return c.JSON(bodyData)
}

func (l *listServiceImpl) Read(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error converting query parameter to an int",
		})
	}

	//Authenticates
	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	userID, _ := strconv.Atoi(claims.Issuer)
	//---------------

	list, err := l.cacheHandler.TryReadingFromCache(userID, id,
		func() (models.List, error) {
			list := models.List{}
			err := list.ReadListById(database.DB, userID, id)
			return list, err
		})
	if err != nil {
		return l.handleListError(err, c)
	}

	return c.JSON(list)
}

func (l *listServiceImpl) ReadAll(c *fiber.Ctx) error { //totest
	//Authenticates
	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	userID, _ := strconv.Atoi(claims.Issuer)
	//---------------

	lists, err := l.cacheHandler.TryReadingAllFromCache(userID,
		func() ([]models.List, error) {
			list := models.List{}
			return list.ReadAllLists(database.DB, userID)
		})
	if err != nil {
		l.handleListError(err, c)
	}

	defaultLists, _ := l.cacheHandler.ReadAllDefaultFromCache(
		func() ([]models.List, error) {
			var list models.List

			return list.ReadAllDefault(database.DB)
		},
	)

	lists = append(lists, defaultLists...)
	return c.JSON(lists)
}

func (l *listServiceImpl) Update(c *fiber.Ctx) error {
	//validates query param
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error converting query parameter to an int",
		})
	}
	//-----------------

	//Authenticates
	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	userID, _ := strconv.Atoi(claims.Issuer)
	//---------------

	//Mounts Body
	var bodyData models.List
	if err := c.BodyParser(&bodyData); err != nil {
		return err
	}
	//------------

	list, err := bodyData.UpdateList(database.DB, userID, id)
	if err != nil {
		return l.handleListError(err, c)
	}

	list, _ = l.cacheHandler.UpdateOnCache(*list)

	return c.JSON(list)
}

func (l *listServiceImpl) Delete(c *fiber.Ctx) error {
	//validates query param
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error converting query parameter to an int",
		})
	}
	//-----------------

	//Authenticates
	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	userID, _ := strconv.Atoi(claims.Issuer)
	//----------------

	var list models.List
	if err := list.DeleteListByID(database.DB, userID, id); err != nil {
		return l.handleListError(err, c)
	}

	l.cacheHandler.DeleteOnCache(list)

	return c.JSON(fiber.Map{
		"message": "List successfully deleted",
		"data":    list,
	})
}

func (l *listServiceImpl) DeleteSymbol(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error converting query parameter to an int",
		})
	}

	//Authenticates
	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	userID, _ := strconv.Atoi(claims.Issuer)
	//-------------

	var symbol models.Symbol
	if err := symbol.DeleteSymbol(database.DB, userID, id); err != nil {
		return l.handleListError(err, c)
	}

	l.cacheHandler.DeleteSymbolOnCache(userID, symbol)

	return c.JSON(fiber.Map{
		"message": "Symbol successfully deleted",
	})
}

func (l *listServiceImpl) CreateSymbol(c *fiber.Ctx) error {
	listID, err := strconv.Atoi(c.Query("list_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error converting list_id to an int",
		})
	}

	//Authenticates
	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	userID, _ := strconv.Atoi(claims.Issuer)
	//-------------

	//Mounts body
	var bodyData models.Symbol
	if err := c.BodyParser(&bodyData); err != nil {
		return err
	}
	//-------------

	if err := bodyData.CreateSymbol(database.DB, listID); err != nil {
		return l.handleListError(err, c)
	}

	l.cacheHandler.CreateSymbolOnCache(userID, bodyData)

	return c.JSON(bodyData)
}

func (l *listServiceImpl) handleListError(err error, c *fiber.Ctx) error {
	if strings.Contains(err.Error(), "Not Found") {
		c.Status(fiber.StatusNotFound)
	} else {
		c.Status(fiber.StatusBadRequest)
	}
	return c.JSON(fiber.Map{
		"message": err.Error(),
	})
}
