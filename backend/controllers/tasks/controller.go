package tasks

import "github.com/gofiber/fiber/v2"

type TaskService interface {
	Create(c *fiber.Ctx) error
	Read(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type taskServiceImpl struct {
	secret string
}

func NewTaskService(secret string) TaskService {
	return &taskServiceImpl{
		secret: secret,
	}
}

func (t *taskServiceImpl) Create(c *fiber.Ctx) error {
	return nil
}
func (t *taskServiceImpl) Read(c *fiber.Ctx) error {
	return nil
}
func (t *taskServiceImpl) Update(c *fiber.Ctx) error {
	return nil
}
func (t *taskServiceImpl) Delete(c *fiber.Ctx) error {
	return nil
}
