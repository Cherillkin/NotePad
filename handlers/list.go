package handlers

import (
	"strconv"

	"github.com/Cherillkin/Notepad/models"
	"github.com/gofiber/fiber/v2"
)

type ListHandler struct {
	service models.ListService
}

func (h *ListHandler) CreateList(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userId").(uint)

	var listData models.List
	if err := ctx.BodyParser(&listData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	newList, err := h.service.CreateList(ctx.Context(), userID, &listData)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   newList,
	})
}

func (h *ListHandler) GetUserLists(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userId").(uint)

	lists, err := h.service.GetUserLists(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   lists,
	})
}

func (h *ListHandler) GetList(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userId").(uint)

	listID, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid list ID",
		})
	}

	list, err := h.service.GetList(ctx.Context(), userID, uint(listID))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   list,
	})
}

func (h *ListHandler) DeleteList(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userId").(uint)

	listID, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"success": "Invalid list ID",
		})
	}

	err = h.service.DeleteList(ctx.Context(), userID, uint(listID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   "List delete",
	})
}

func (h *ListHandler) UpdateList(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userId").(uint)

	listID, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid list ID",
		})
	}

	var updateData models.List
	if err := ctx.BodyParser(&updateData); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	updateList, err := h.service.UpdateList(ctx.Context(), userID, uint(listID), &updateData)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   updateList,
	})
}

func NewListHandler(router fiber.Router, service models.ListService) {
	handler := &ListHandler{
		service: service,
	}

	router.Post("/list", handler.CreateList)
	router.Get("/lists", handler.GetUserLists)
	router.Get("/list/:id", handler.GetList)
	router.Delete("/list/:id", handler.DeleteList)
	router.Put("/list/:id", handler.UpdateList)
}
