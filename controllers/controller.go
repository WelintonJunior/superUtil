package sucontrollers

import (
	"strconv"

	surepository "github.com/WelintonJunior/superUtil/repository"
	"github.com/WelintonJunior/superUtil/validate"
	"github.com/gofiber/fiber/v2"
)

type SuperUtilController[T any] struct {
	superUtilRepository surepository.SuperUtilRepository[T]
}

func NewSuperUtilController[T any](superUtilRepository surepository.SuperUtilRepository[T]) *SuperUtilController[T] {
	return &SuperUtilController[T]{superUtilRepository: superUtilRepository}
}

func (c *SuperUtilController[T]) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var req T

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if err := validate.ValidateStruct(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.superUtilRepository.Create(&req); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create item"})
		}

		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Item created successfully", "data": req})
	}
}

func (c *SuperUtilController[T]) GetByID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		idParam := ctx.Params("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid id param"})
		}

		item, err := c.superUtilRepository.GetByID(uint(id))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve item"})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Item retrieved successfully", "data": item})
	}
}

func (c *SuperUtilController[T]) GetAll() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		items, err := c.superUtilRepository.GetAll()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve items"})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Items retrieved successfully", "data": items})
	}
}

func (c *SuperUtilController[T]) Update() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var req T

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if err := validate.ValidateStruct(&req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.superUtilRepository.Update(&req); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update item"})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Item updated successfully", "data": req})
	}
}

func (c *SuperUtilController[T]) DeleteByID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		idParam := ctx.Params("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid id param"})
		}

		item, err := c.superUtilRepository.GetByID(uint(id))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve item"})
		}

		if err := c.superUtilRepository.DeleteByID(uint(id)); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete item"})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Item deleted successfully", "data": item})
	}
}
