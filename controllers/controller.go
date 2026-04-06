package sucontrollers

import (
	"errors"

	surepository "github.com/WelintonJunior/superUtil/repository"
	"github.com/WelintonJunior/superUtil/validate"
	"github.com/gofiber/fiber/v2"
)

// Bloqueio de payload grande para evitar ataques de negação de serviço (DoS) com corpos de requisição excessivamente grandes.
const maxRequestBodyBytes = 1 << 20 // 1 MiB

type SuperUtilController[T any] struct {
	superUtilRepository surepository.SuperUtilRepository[T]
}

func NewSuperUtilController[T any](superUtilRepository surepository.SuperUtilRepository[T]) *SuperUtilController[T] {
	return &SuperUtilController[T]{superUtilRepository: superUtilRepository}
}

func (c *SuperUtilController[T]) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var req T

		if err := ensureJSONRequest(ctx); err != nil {
			fiberErr, ok := err.(*fiber.Error)
			if ok {
				return ctx.Status(fiberErr.Code).JSON(fiber.Map{"error": fiberErr.Message})
			}

			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

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
		id, err := parseIDParam(ctx)
		if err != nil {
			fiberErr, ok := err.(*fiber.Error)
			if ok {
				return ctx.Status(fiberErr.Code).JSON(fiber.Map{"error": fiberErr.Message})
			}

			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid id param"})
		}

		item, err := c.superUtilRepository.GetByID(id)
		if err != nil {
			if errors.Is(err, surepository.ErrNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
			}

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

		if err := ensureJSONRequest(ctx); err != nil {
			fiberErr, ok := err.(*fiber.Error)
			if ok {
				return ctx.Status(fiberErr.Code).JSON(fiber.Map{"error": fiberErr.Message})
			}

			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

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
		id, err := parseIDParam(ctx)
		if err != nil {
			fiberErr, ok := err.(*fiber.Error)
			if ok {
				return ctx.Status(fiberErr.Code).JSON(fiber.Map{"error": fiberErr.Message})
			}

			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid id param"})
		}

		item, err := c.superUtilRepository.GetByID(id)
		if err != nil {
			if errors.Is(err, surepository.ErrNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve item"})
		}

		if err := c.superUtilRepository.DeleteByID(id); err != nil {
			if errors.Is(err, surepository.ErrNotFound) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete item"})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Item deleted successfully", "data": item})
	}
}
