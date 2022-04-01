package lists

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

func (l *listServiceImpl) Create(c *fiber.Ctx) error {
	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var bodyData models.List
	if err := c.BodyParser(&bodyData); err != nil {
		return err
	}

	user_id, err := strconv.Atoi(claims.Issuer)
	bodyData.UserID = uint(user_id)

	var count int64
	database.DB.Model(&models.List{}).Where("user_id=?", bodyData.UserID).Count(&count)

	if count >= MaxListsAmount {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "user exceeded maximum list amount",
			"amount":  MaxListsAmount,
		})
	}

	if err := database.DB.Create(&bodyData).Error; err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "error finding the list",
			"error":   err.Error(),
		})
	}

	l.cache.Delete(fmt.Sprintf("%d", user_id))
	l.cache.SetDefault(fmt.Sprintf("%d_%d", user_id, bodyData.ID), bodyData)

	return c.JSON(bodyData)
}

func (l *listServiceImpl) Read(c *fiber.Ctx) error {
	id := c.Query("id")

	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	user_id, err := strconv.Atoi(claims.Issuer)

	listCache, existsOnCache := l.cache.Get(fmt.Sprintf("%d_%s", user_id, id))

	if !existsOnCache {
		var list models.List
		if err := database.DB.Where("user_id=? AND id=?", user_id, id).Find(&list).Error; err != nil {
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

		l.cache.SetDefault(fmt.Sprintf("%d_%s", user_id, id), list)

		return c.JSON(list)
	}

	return c.JSON(listCache)
}

func (l *listServiceImpl) ReadAll(c *fiber.Ctx) error {
	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	user_id, err := strconv.Atoi(claims.Issuer)

	listsCache, existsOnCache := l.cache.Get(fmt.Sprintf("%d", user_id))

	if !existsOnCache {
		var lists []models.List
		if err := database.DB.Where("user_id=?", user_id).Find(&lists).Error; err != nil {
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
		}

		l.cache.SetDefault(fmt.Sprintf("%d", user_id), lists)

		return c.JSON(lists)
	}

	return c.JSON(listsCache)
}

func (l *listServiceImpl) Update(c *fiber.Ctx) error {
	id := c.Query("id")

	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var bodyData models.List
	if err := c.BodyParser(&bodyData); err != nil {
		return err
	}

	user_id, err := strconv.Atoi(claims.Issuer)

	var list models.List
	if err := database.DB.Where("user_id=? AND id=?", user_id, id).Find(&list).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error finding the list",
			"error":   err.Error(),
		})
	}

	list.IsDefault = bodyData.IsDefault
	list.Name = bodyData.Name
	list.UserID = bodyData.UserID
	if bodyData.Symbols != nil {
		var symbol models.Symbol
		database.DB.Where("list_id=?", id).Delete(&symbol)
		// var symbolIds []uint
		// for i, e := range bodyData.Symbols {
		// 	database.DB.Where("symbol=? AND list_id=?", e.Symbol, e.ListID).Find(&symbol)

		// 	if symbol.ID != 0 {
		// 		bodyData.Symbols[i].ID = symbol.ID
		// 		symbolIds = append(symbolIds, symbol.ID)
		// 	}
		// }
		// database.DB.Where("list_id=? AND id NOT IN ?", list.ID, symbolIds).Delete(&models.Symbol{})

		list.Symbols = bodyData.Symbols
	}

	database.DB.Where("id=?", id).Save(&list)

	l.cache.Delete(fmt.Sprintf("%d", user_id))
	l.cache.Delete(fmt.Sprintf("%d_%s", user_id, id))

	return c.JSON(list)
}

func (l *listServiceImpl) Delete(c *fiber.Ctx) error {
	id := c.Query("id")

	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	user_id, err := strconv.Atoi(claims.Issuer)

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

	l.cache.Delete(fmt.Sprintf("%d", user_id))
	l.cache.Delete(fmt.Sprintf("%d_%s", user_id, id))

	return c.JSON(fiber.Map{
		"message": "List successfully deleted",
	})
}

func (l *listServiceImpl) DeleteSymbol(c *fiber.Ctx) error {
	id := c.Query("id")

	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	user_id, err := strconv.Atoi(claims.Issuer)

	var symbol models.Symbol
	database.DB.Where("id=?", id).Find(&symbol)

	if symbol.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Could not found symbol by the given id",
			"id":      id,
		})
	}

	if err := database.DB.Delete(&symbol, id).Error; err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "error deleting the symbol",
			"error":   err.Error(),
		})
	}

	l.cache.Delete(fmt.Sprintf("%d", user_id))
	l.cache.Delete(fmt.Sprintf("%d_%d", user_id, symbol.ListID))

	return c.JSON(fiber.Map{
		"message": "Symbol successfully deleted",
	})
}

const MaxSimbolsAmount = 50

func (l *listServiceImpl) CreateSymbol(c *fiber.Ctx) error {
	claims, err := middleware.Auth(l.secret, c)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	user_id, err := strconv.Atoi(claims.Issuer)

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

	l.cache.Delete(fmt.Sprintf("%d", user_id))
	l.cache.Delete(fmt.Sprintf("%d_%d", user_id, bodyData.ListID))
	return c.JSON(bodyData)
}
