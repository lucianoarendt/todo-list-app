package other

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/patrickmn/go-cache"
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
	secret string
	cache  cache.Cache
}

func NewListService(secret string) ListService {
	return &listServiceImpl{
		secret: secret,
		cache:  *cache.New(-1, -1),
	}
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
	//---------------

	//mounts body
	var bodyData models.List
	if err := c.BodyParser(&bodyData); err != nil {
		return err
	}

	userID, _ := strconv.Atoi(claims.Issuer)
	bodyData.UserID = uint(userID)
	//--------------

	if err := bodyData.CreateList(database.DB); err != nil {
		l.handleListError(err, c)
	}

	//Handles cache
	ids, userFound := l.cache.Get(claims.Issuer)
	if userFound {
		//At user level
		idsSlice := ids.([]string)
		idsSlice = append(idsSlice, strconv.Itoa(int(bodyData.ID)))
		l.cache.SetDefault(claims.Issuer, idsSlice)
	}
	//At list level
	l.cache.SetDefault(fmt.Sprintf("%s_%d", claims.Issuer, bodyData.ID), bodyData)
	//---------------

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

	cacheKey := fmt.Sprintf("%d_%d", userID, id)
	listCache, existsOnCache := l.cache.Get(cacheKey)

	var list models.List
	if !existsOnCache {
		if err := list.ReadListById(database.DB, userID, id); err != nil {
			l.handleListError(err, c)
		}

		l.cache.SetDefault(cacheKey, list)
	} else {
		list = listCache.(models.List)
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

	lists, existsOnCache := l.getCacheLists(claims.Issuer)

	if !existsOnCache {
		list := models.List{}
		lists, err := list.ReadAllLists(database.DB, userID)
		if err != nil {
			l.handleListError(err, c)
		}

		ids := make([]string, len(lists))
		for i := range lists {
			lists[i].PopulateWithSymbols(database.DB)

			ids[i] = strconv.Itoa(int(lists[i].ID))
			l.cache.SetDefault(fmt.Sprintf("%s_%d", claims.Issuer, lists[i].ID), lists[i])
		}
		l.cache.SetDefault(claims.Issuer, ids)
	}
	lists = append(lists, l.readAllDefaultLists()...)
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
	//---------------

	//Mounts Body
	var bodyData models.List
	if err := c.BodyParser(&bodyData); err != nil {
		return err
	}

	userID, _ := strconv.Atoi(claims.Issuer)
	bodyData.UserID = uint(userID)
	//------------

	list, err := bodyData.UpdateList(database.DB, id)
	if err != nil {
		l.handleListError(err, c)
	}

	//Handle Cache
	cacheKey := fmt.Sprintf("%s_%d", claims.Issuer, id)
	if list.Symbols == nil {
		cacheList, existsOnCache := l.cache.Get(cacheKey)
		if existsOnCache {
			cacheAsList := cacheList.(models.List)
			list.Symbols = cacheAsList.Symbols
		} else {
			list.PopulateWithSymbols(database.DB)
		}
	}

	l.cache.SetDefault(cacheKey, list)
	//-----------------------

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
	//----------------

	var list models.List
	if err := list.DeleteListByID(database.DB, id); err != nil {
		l.handleListError(err, c)
	}

	//Handles Cache
	ids, userFound := l.cache.Get(claims.Issuer)
	if userFound {
		//Delete id from user cache
		idsSlice := ids.([]string)
		for i, e := range idsSlice {
			if e == strconv.Itoa(id) {
				idsSlice[i] = idsSlice[len(idsSlice)-1]
				idsSlice = idsSlice[:len(idsSlice)-1]
				break
			}
		}
		l.cache.SetDefault(claims.Issuer, idsSlice)
		//-----------------------
	}
	l.cache.Delete(fmt.Sprintf("%s_%d", claims.Issuer, id))
	//-----------------

	return c.JSON(fiber.Map{
		"message": "List successfully deleted",
	})
}

func (l *listServiceImpl) DeleteSymbol(c *fiber.Ctx) error {
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
	//-------------

	//Gets Symbol and validates existence
	var symbol models.Symbol
	database.DB.Where("id=?", id).Find(&symbol)

	if symbol.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Could not found symbol by the given id",
			"id":      id,
		})
	}
	//-------------------------

	if err := database.DB.Delete(&symbol, id).Error; err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "error deleting the symbol",
			"error":   err.Error(),
		})
	}

	//Cache handling
	cacheKey := fmt.Sprintf("%s_%d", claims.Issuer, symbol.ListID)

	listCache, existsOnCache := l.cache.Get(cacheKey)

	if existsOnCache {
		cacheAsList := listCache.(models.List)
		SymbolsCopy := make([]models.Symbol, len(cacheAsList.Symbols))
		copy(SymbolsCopy, cacheAsList.Symbols)

		//Delete a symbol
		for i, e := range SymbolsCopy {
			if e.ID == uint(id) {
				SymbolsCopy[i] = SymbolsCopy[len(SymbolsCopy)-1]
				SymbolsCopy = SymbolsCopy[:len(SymbolsCopy)-1]
			}
		}
		//-----------------
		cacheAsList.Symbols = SymbolsCopy

		l.cache.SetDefault(cacheKey, cacheAsList)
	}
	//-----------------

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
	//-------------

	//Mounts body
	var bodyData models.Symbol
	if err := c.BodyParser(&bodyData); err != nil {
		return err
	}

	if err := database.DB.Where("symbol=? AND list_id=?", bodyData.Symbol, listID).Find(&bodyData).Error; err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "error trying to find the symbol",
			"err":     err.Error(),
		})
	}
	bodyData.ListID = uint(listID)
	//-------------

	//validates max symbols amount
	var count int64
	database.DB.Model(&models.Symbol{}).Where("list_id=?", listID).Count(&count)

	if count >= MaxSimbolsAmount {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "list exceeded maximum symbol amount",
			"amount":  MaxSimbolsAmount,
		})
	}
	//-------------

	//validates inexistence
	if bodyData.ID != 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Symbol already exists in this list",
			"data":    bodyData,
		})
	}
	//-------------

	if err := database.DB.Create(&bodyData).Error; err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "error creating the symbol",
			"error":   err.Error(),
		})
	}

	//Cache Handling
	cacheKey := fmt.Sprintf("%s_%d", claims.Issuer, bodyData.ListID)

	listCache, existsOnCache := l.cache.Get(cacheKey)

	if existsOnCache {
		cacheAsList := listCache.(models.List)
		SymbolsCopy := make([]models.Symbol, len(cacheAsList.Symbols))
		copy(SymbolsCopy, cacheAsList.Symbols)

		cacheAsList.Symbols = append(SymbolsCopy, bodyData)
		l.cache.SetDefault(cacheKey, cacheAsList)
	}
	//--------------------
	return c.JSON(bodyData)
}

func (l *listServiceImpl) getCacheLists(userID string) ([]models.List, bool) {
	var lists []models.List

	ids, foundUser := l.cache.Get(userID)
	if foundUser {
		idsAsSlice := ids.([]string)

		for _, e := range idsAsSlice {
			listCache, _ := l.cache.Get(userID + "_" + e)
			listCacheAsList := listCache.(models.List)

			lists = append(lists, listCacheAsList)
		}
	}

	return lists, foundUser
}

func (l *listServiceImpl) readAllDefaultLists() []models.List { //totest
	var lists []models.List

	listsCache, existsOnCache := l.cache.Get("default")

	if !existsOnCache {
		database.DB.Where("is_default=1").Find(&lists)

		for i := range lists {
			lists[i].PopulateWithSymbols(database.DB)
		}
		l.cache.SetDefault("default", lists) //Set("default", lists, DurationX)

	} else {
		lists = listsCache.([]models.List)
	}

	return lists
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
