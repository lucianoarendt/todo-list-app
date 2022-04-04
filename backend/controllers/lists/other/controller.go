package other

import (
	"fmt"
	"strconv"

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

	//validates max list amount
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

func (l *listServiceImpl) ReadAll(c *fiber.Ctx) error {
	//Authenticates
	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	//---------------

	lists, existsOnCache := l.getCacheLists(claims.Issuer)

	if !existsOnCache {
		var lists []models.List
		if err := database.DB.Where("user_id=?", claims.Issuer).Find(&lists).Error; err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "error finding the list",
				"err":     err.Error(),
			})
		}

		var ids []string
		for i := range lists {
			if err := database.DB.Where("list_id=?", lists[i].ID).Find(&lists[i].Symbols).Error; err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(fiber.Map{
					"message": "error finding the symbols",
					"err":     err.Error(),
				})
			}
			ids = append(ids, strconv.Itoa(int(lists[i].ID)))
		}

		l.cache.SetDefault(claims.Issuer, ids)
		for _, list := range lists {
			l.cache.SetDefault(fmt.Sprintf("%s_%d", claims.Issuer, list.ID), list)
		}

		return c.JSON(lists)
	}

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

	//Updates list
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

	//Handles Cache
	ids, userFound := l.cache.Get(claims.Issuer)
	idsSlice := ids.([]string)
	if userFound {
		//Delete id from user cache
		for i, e := range idsSlice {
			if e == id {
				idsSlice[i] = idsSlice[len(idsSlice)-1]
				idsSlice = idsSlice[:len(idsSlice)-1]
				break
			}
		}
		l.cache.SetDefault(claims.Issuer, idsSlice)
		//-----------------------
	}
	l.cache.Delete(fmt.Sprintf("%s_%s", claims.Issuer, id))
	//-----------------

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
	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var bodyData models.Symbol
	if err := c.BodyParser(&bodyData); err != nil {
		return err
	}

	if err := database.DB.Where("symbol=? AND list_id=?", bodyData.Symbol, bodyData.ListID).Find(&bodyData).Error; err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "error trying to find the symbol",
			"err":     err.Error(),
		})
	}

	var count int64
	database.DB.Model(&models.Symbol{}).Where("list_id=?", bodyData.ListID).Count(&count)

	if count >= MaxSimbolsAmount {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "list exceeded maximum symbol amount",
			"amount":  MaxSimbolsAmount,
		})
	}

	if bodyData.ID != 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Symbol already exists in this list",
			"data":    bodyData,
		})
	}

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
