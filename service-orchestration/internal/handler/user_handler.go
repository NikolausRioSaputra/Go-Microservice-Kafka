package handler

import (
	"net/http"
	"service-orchestration/m/internal/domain"
	"service-orchestration/m/internal/middleware"
	"service-orchestration/m/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(uc usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: uc,
	}
}

// StoreNewUser creates a new user and returns a JWT token
func (h *UserHandler) StoreNewUser(c *gin.Context) {
	var userRequest domain.User
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Create user
	// createdUser, err := h.userUseCase.CreateUser(c.Request.Context(), &userRequest)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	// 	return
	// }

	// Generate JWT token for the created user
	token, err := middleware.GenerateToken() // Here you might want to pass some user details to the token generator
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return user data and token
	c.JSON(http.StatusOK, gin.H{
		// "user":  createdUser,
		"token": token,
	})
}
