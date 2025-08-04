package handlers

import (
	"strconv"

	"github.com/Cherillkin/Notepad/models"
	"github.com/gofiber/fiber/v2"
)

type SharedListHandler struct {
	service models.SharedListService
}

func (h *SharedListHandler) SharedList(c *fiber.Ctx) error {
	listID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid list ID",
		})
	}

	var body struct {
		UserID uint `json:"user_id"`
	}
	if err := c.BodyParser(&body); err != nil || body.UserID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "User ID is required",
		})
	}

	if err := h.service.SharedList(c.Context(), uint(listID), body.UserID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "List shared successfully",
	})
}

func (h *SharedListHandler) GetSharedLists(c *fiber.Ctx) error {
	userID := c.Locals("userId").(uint)

	lists, err := h.service.GetSharedLists(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   lists,
	})
}

func NewSharedListHandler(router fiber.Router, service models.SharedListService) {
	handler := &SharedListHandler{
		service: service,
	}
	router.Post("/:id/share", handler.SharedList)
	router.Get("/shared", handler.GetSharedLists)
}
