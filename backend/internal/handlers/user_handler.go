package handlers

import (
	"net/http"

	"ragify-backend/internal/utils"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	// In a real implementation, this would contain service dependencies
	// For now, we'll use a placeholder
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// CreateUser handles user registration
func (h *UserHandler) CreateUser(c echo.Context) error {
	// Request structure
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	if req.Email == "" || req.Password == "" {
		return utils.BadRequest(c, "Email and password are required")
	}

	// In a real implementation, this would create a user in the database
	user := map[string]interface{}{
		"id":    1,
		"name":  req.Name,
		"email": req.Email,
	}

	return utils.SendSuccess(c, "User created successfully", user, http.StatusCreated)
}

// Login handles user login
func (h *UserHandler) Login(c echo.Context) error {
	// Request structure
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request format")
	}

	if req.Email == "" || req.Password == "" {
		return utils.BadRequest(c, "Email and password are required")
	}

	// In a real implementation, this would verify credentials and generate a JWT token
	response := map[string]interface{}{
		"token": "sample_jwt_token",
		"user": map[string]interface{}{
			"id":    1,
			"name":  "Sample User",
			"email": req.Email,
		},
	}

	return utils.SendSuccess(c, "Login successful", response, http.StatusOK)
}
