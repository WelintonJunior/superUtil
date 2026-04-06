package sucontrollers

import (
	"mime"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func parseIDParam(ctx *fiber.Ctx) (uint, error) {
	idParam := strings.TrimSpace(ctx.Params("id"))
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "Invalid id param")
	}

	return uint(id), nil
}

func ensureJSONRequest(ctx *fiber.Ctx) error {
	contentType := ctx.Get(fiber.HeaderContentType)
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return fiber.NewError(fiber.StatusUnsupportedMediaType, "Invalid Content-Type header")
	}

	if mediaType != fiber.MIMEApplicationJSON {
		return fiber.NewError(fiber.StatusUnsupportedMediaType, "Content-Type must be application/json")
	}

	if len(ctx.Body()) > maxRequestBodyBytes {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "Request body too large")
	}

	return nil
}
