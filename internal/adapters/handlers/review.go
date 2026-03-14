package handlers

import (
	"net/http"
	"strconv"

	"tienda-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	service ports.ReviewService
}

func NewReviewHandler(service ports.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

// AddReview godoc
// @Summary Add a review to a product
// @Tags Reviews
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param rating formData int true "Rating (1-5)"
// @Param comment formData string true "Review comment"
// @Param image formData file false "Review image"
// @Success 201 {object} domain.Review
// @Router /products/{id}/reviews [post]
func (h *ReviewHandler) AddReview(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Obtener UserID del token JWT (inyectado por el middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	rating, _ := strconv.Atoi(c.PostForm("rating"))
	comment := c.PostForm("comment")
	file, _ := c.FormFile("image")

	review, err := h.service.AddReview(uint(productID), userID.(uint), rating, comment, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// GetProductReviews godoc
// @Summary List reviews for a product
// @Tags Reviews
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {array} domain.Review
// @Router /products/{id}/reviews [get]
func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	reviews, err := h.service.GetProductReviews(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}
