package handlers

import (
	"fmt"
	"strconv"

	"github.com/Cherillkin/Notepad/models"
	"github.com/Cherillkin/Notepad/utils"
	"github.com/gofiber/fiber/v2"
)

type SharedListHandler struct {
	service  models.SharedListService
	producer *utils.Producer
}

func (h *SharedListHandler) SharedList(ctx *fiber.Ctx) error {
	listID, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid list ID",
		})
	}

	var body struct {
		UserID uint `json:"user_id"`
	}
	if err := ctx.BodyParser(&body); err != nil || body.UserID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "User ID is required",
		})
	}

	if err := h.service.SharedList(ctx.Context(), uint(listID), body.UserID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	userID := ctx.Locals("userId").(uint)
	msg := fmt.Sprintf("User %d shared list %d with user %d", userID, listID, body.UserID)
	_ = h.producer.SendMessage("list shared", msg)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "List shared successfully",
	})
}

func (h *SharedListHandler) GetSharedLists(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userId").(uint)

	lists, err := h.service.GetSharedLists(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   lists,
	})
}

func NewSharedListHandler(router fiber.Router, service models.SharedListService, producer *utils.Producer) {
	handler := &SharedListHandler{
		service:  service,
		producer: producer,
	}
	router.Post("/:id/share", handler.SharedList)
	router.Get("/shared", handler.GetSharedLists)
}
