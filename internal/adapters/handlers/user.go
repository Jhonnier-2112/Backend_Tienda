package handlers

import (
	"net/http"
	"tienda-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetAllUsers godoc
// @Summary List all users
// @Description Gets a list of all users and their profile data
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.User
// @Failure 500 {object} map[string]string "Failed to fetch users"
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

type CreateSellerInput struct {
	CompanyName string `json:"company_name" binding:"required"`
	ContactName string `json:"contact_name" binding:"required"`
	NIT         string `json:"nit" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
}

// CreateSeller godoc
// @Summary Create a new Seller
// @Description Allows an Admin to create a new Seller with NIT and Company Name
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body CreateSellerInput true "Seller Registration Data"
// @Success 201 {object} map[string]interface{} "Seller created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /users/seller [post]
func (h *UserHandler) CreateSeller(c *gin.Context) {
	var input CreateSellerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateSeller(input.CompanyName, input.ContactName, input.NIT, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Seller created successfully", "user": gin.H{
		"id":           user.ID,
		"company_name": user.SellerData.CompanyName,
		"nit":          user.SellerData.NIT,
		"email":        user.Email,
		"role":         user.Role,
	}})
}

type CreateAdminInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// CreateAdmin godoc
// @Summary Create a new Admin
// @Description Allows an existing Admin to create a new Admin account
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body CreateAdminInput true "Admin Registration Data"
// @Success 201 {object} map[string]interface{} "Admin created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /users/admin [post]
func (h *UserHandler) CreateAdmin(c *gin.Context) {
	var input CreateAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateAdmin(input.Name, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Admin created successfully", "user": gin.H{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	}})
}
