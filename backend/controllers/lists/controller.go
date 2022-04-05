package lists

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
	//--------------

	//validates max list amount
	userID, _ := strconv.Atoi(claims.Issuer)
	bodyData.UserID = uint(userID)

	var count int64
	database.DB.Model(&models.List{}).Where("user_id=?", claims.Issuer).Count(&count)

	if count >= MaxListsAmount {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "user exceeded maximum list amount",
			"amount":  MaxListsAmount,
		})
	}
	//-----------------------

	if err := database.DB.Create(&bodyData).Error; err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "error finding the list",
			"error":   err.Error(),
		})
	}

	//Handles cache
	l.cache.SetDefault(fmt.Sprintf("%s_%d", claims.Issuer, bodyData.ID), bodyData)

	return c.JSON(bodyData)
}

func (l *listServiceImpl) Read(c *fiber.Ctx) error {
	id := c.Query("id")

	//Authenticates
	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	//---------------

	listCache, existsOnCache := l.cache.Get(fmt.Sprintf("%s_%s", claims.Issuer, id))

	if !existsOnCache {
		var list models.List
		if err := database.DB.Where("user_id=? AND id=?", claims.Issuer, id).Find(&list).Error; err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "error finding the list",
				"error":   err.Error(),
			})
		}

		if list.ID == 0 {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "List Not Found",
			})
		}

		if err := database.DB.Where("list_id=?", list.ID).Find(&list.Symbols).Error; err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "error finding the symbols",
				"err":     err.Error(),
			})
		}

		l.cache.SetDefault(fmt.Sprintf("%s_%s", claims.Issuer, id), list)

		return c.JSON(list)
	}

	return c.JSON(listCache)
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
	//---------------

	_, existsOnCache := l.cache.Get(claims.Issuer)

	var lists []models.List
	if !existsOnCache {
		if err := database.DB.Where("user_id=?", claims.Issuer).Find(&lists).Error; err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "error finding the list",
				"err":     err.Error(),
			})
		}

		for i := range lists {
			if err := database.DB.Where("list_id=?", lists[i].ID).Find(&lists[i].Symbols).Error; err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"message": "error finding the symbols",
					"err":     err.Error(),
				})
			}
			l.cache.SetDefault(fmt.Sprintf("%s_%d", claims.Issuer, lists[i].ID), lists[i])
		}
		l.cache.SetDefault(claims.Issuer, true)
	} else {
		lists = l.getCacheLists(claims.Issuer)
	}

	lists = append(lists, l.readAllDefaultLists()...)
	return c.JSON(lists)
}

func (l *listServiceImpl) Update(c *fiber.Ctx) error {
	id := c.Query("id")

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

	//Gets list and validates inexistence
	var list models.List
	if err := database.DB.Where("user_id=? AND id=?", claims.Issuer, id).Find(&list).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error finding the list",
			"error":   err.Error(),
		})
	}

	if list.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "List Not Found",
		})
	}
	//----------------

	//Updates list fields
	list.IsDefault = bodyData.IsDefault
	list.Name = bodyData.Name
	list.UserID = bodyData.UserID
	if bodyData.Symbols != nil {
		var symbol models.Symbol
		database.DB.Where("list_id=?", id).Delete(&symbol)

		list.Symbols = bodyData.Symbols
	} else {
		cacheList, existsOnCache := l.cache.Get(fmt.Sprintf("%s_%s", claims.Issuer, id))
		if existsOnCache {
			cacheAsList := cacheList.(models.List)
			list.Symbols = cacheAsList.Symbols
		} else {
			if err := database.DB.Where("list_id=?", list.ID).Find(&list.Symbols).Error; err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"message": "error finding the symbols",
					"err":     err.Error(),
				})
			}
		}
	}

	database.DB.Where("id=?", id).Save(&list)
	//-----------------

	l.cache.SetDefault(fmt.Sprintf("%s_%s", claims.Issuer, id), list)

	return c.JSON(list)
}

func (l *listServiceImpl) Delete(c *fiber.Ctx) error {
	id := c.Query("id")

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
	database.DB.Where("id=?", id).Find(&list)
	if list.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Could not found list by the given id",
			"id":      id,
		})
	}

	if err := database.DB.Delete(&list, id).Error; err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "error deleting the list",
			"error":   err.Error(),
		})
	}

	l.cache.Delete(fmt.Sprintf("%s_%s", claims.Issuer, id))

	return c.JSON(fiber.Map{
		"message": "List successfully deleted",
	})
}

func (l *listServiceImpl) DeleteSymbol(c *fiber.Ctx) error {
	id := c.Query("id")

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
		symbolID, _ := strconv.Atoi(id)
		for i, e := range SymbolsCopy {
			if e.ID == uint(symbolID) {
				SymbolsCopy[i] = SymbolsCopy[len(SymbolsCopy)-1]
				SymbolsCopy = SymbolsCopy[:len(SymbolsCopy)-1]
				break
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

func (l *listServiceImpl) getCacheLists(userID string) []models.List {
	itemsMap := l.cache.Items()
	var lists []models.List

	for k := range itemsMap {
		if strings.Contains(k, fmt.Sprintf("%s_", userID)) {
			lists = append(lists, itemsMap[k].Object.(models.List))
		}
	}

	return lists
}

func (l *listServiceImpl) readAllDefaultLists() []models.List { //totest
	var lists []models.List

	listsCache, existsOnCache := l.cache.Get("default")

	if !existsOnCache {
		database.DB.Where("is_default=1").Find(&lists)

		for i := range lists {
			database.DB.Where("list_id=?", lists[i].ID).Find(&lists[i].Symbols)
		}
		l.cache.SetDefault("default", lists) //Set("default", lists, DurationX)

	} else {
		lists = listsCache.([]models.List)
	}

	return lists
}
