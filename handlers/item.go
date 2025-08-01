package handlers

import (
	"strconv"

	"github.com/Cherillkin/Notepad/models"
	"github.com/gofiber/fiber/v2"
)

type ItemHandler struct {
	service models.ItemService
}

func (h *ItemHandler) CreateItem(ctx *fiber.Ctx) error {
	listID := ctx.Locals("listId").(uint)

	var itemData models.Item
	if err := ctx.BodyParser(&itemData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	newItem, err := h.service.CreateItem(ctx.Context(), listID, &itemData)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   newItem,
	})
}

func (h *ItemHandler) GetListItems(ctx *fiber.Ctx) error {
	listID := ctx.Locals("listId").(uint)

	items, err := h.service.GetListItems(ctx.Context(), listID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "succees",
		"data":   items,
	})
}

func (h *ItemHandler) GetItem(ctx *fiber.Ctx) error {
	listID := ctx.Locals("listId").(uint)

	itemID, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid item ID",
		})
	}

	item, err := h.service.GetItem(ctx.Context(), listID, uint(itemID))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   item,
	})
}

func (h *ItemHandler) DeleteItem(ctx *fiber.Ctx) error {
	listID := ctx.Locals("listId").(uint)

	itemID, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid Item ID",
		})
	}

	if err = h.service.DeleteItem(ctx.Context(), listID, uint(itemID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   "Item delete",
	})
}

func (h *ItemHandler) UpdateItem(ctx *fiber.Ctx) error {
	listID := ctx.Locals("listId").(uint)

	itemID, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid Item ID",
		})
	}

	var updateData models.Item
	if err := ctx.BodyParser(&updateData); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	updateItem, err := h.service.UpdateItem(ctx.Context(), listID, uint(itemID), &updateData)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   updateItem,
	})
}

func NewItemHandler(router fiber.Router, service models.ItemService) {
	handler := &ItemHandler{
		service: service,
	}

	router.Post("/", handler.CreateItem)
	router.Get("/", handler.GetListItems)
	router.Get("/:id", handler.GetItem)
	router.Delete("/:id", handler.DeleteItem)
	router.Put("/:id", handler.UpdateItem)
}
