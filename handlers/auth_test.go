package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cherillkin/Notepad/models"
	mock_models "github.com/Cherillkin/Notepad/models/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type Handler struct {
	repo models.AuthRepository
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func (h *Handler) Register(ctx *fiber.Ctx) error {
	var creds models.AuthCredentials
	if err := ctx.BodyParser(&creds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid input",
		})
	}

	user, err := h.repo.RegisterUser(ctx.Context(), &creds)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "failed to register user",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   UserResponse{ID: user.ID, Email: user.Email},
	})
}

func TestHandler_Register(t *testing.T) {
	type mockBehavior func(s *mock_models.MockAuthRepository, user models.AuthCredentials)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            models.AuthCredentials
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"Email": "test@example.com", "Password": "qwerty"}`,
			inputUser: models.AuthCredentials{
				Email:    "test@example.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_models.MockAuthRepository, user models.AuthCredentials) {
				s.EXPECT().RegisterUser(gomock.Any(), &user).Return(&models.User{ID: 1, Email: user.Email}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"success","data":{"id":1,"email":"test@example.com"}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_models.NewMockAuthRepository(c)
			testCase.mockBehavior(repo, testCase.inputUser)

			handler := &Handler{repo: repo}
			app := fiber.New()
			app.Post("/register", handler.Register)

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.JSONEq(t, testCase.expectedResponseBody, string(body))
		})
	}
}
