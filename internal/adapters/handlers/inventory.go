package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service ports.InventoryService
}

func NewCategoryHandler(service ports.InventoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// CreateCategory godoc
// @Summary Create a Category
// @Description Creates a new product category (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body domain.Category true "Category Data"
// @Success 201 {object} domain.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input domain.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := h.service.CreateCategory(input.Name, input.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, category)
}

// GetCategories godoc
// @Summary List all Categories
// @Tags Categories
// @Produce json
// @Success 200 {array} domain.Category
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// GetCategory godoc
// @Summary Get a Category by ID
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} domain.Category
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	category, err := h.service.GetCategoryByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.JSON(http.StatusOK, category)
}

// UpdateCategory godoc
// @Summary Update a Category
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param input body domain.Category true "Updated Category Data"
// @Success 200 {object} domain.Category
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	var input domain.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := h.service.UpdateCategory(uint(id), input.Name, input.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// DeleteCategory godoc
// @Summary Delete a Category
// @Tags Categories
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	if err := h.service.DeleteCategory(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// ─── Product Handler ───────────────────────────────────────────────────────

type ProductHandler struct {
	service         ports.InventoryService
	discountService ports.DiscountService
}

func NewProductHandler(service ports.InventoryService, discountService ports.DiscountService) *ProductHandler {
	return &ProductHandler{service: service, discountService: discountService}
}

// CreateProduct godoc
// @Summary Create a Product
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Product name"
// @Param description formData string false "Description"
// @Param sku formData string true "SKU (unique reference)"
// @Param price formData number true "Sale price"
// @Param cost_price formData number false "Cost price"
// @Param stock formData int false "Initial stock"
// @Param min_stock formData int false "Low-stock alert threshold (default: 5)"
// @Param category_id formData int true "Category ID"
// @Param image formData file false "Product image"
// @Success 201 {object} domain.Product
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")
	sku := c.PostForm("sku")
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
	costPrice, _ := strconv.ParseFloat(c.PostForm("cost_price"), 64)
	stock, _ := strconv.Atoi(c.PostForm("stock"))
	minStock, _ := strconv.Atoi(c.PostForm("min_stock"))
	categoryID, _ := strconv.Atoi(c.PostForm("category_id"))

	if minStock == 0 {
		minStock = 5
	}

	file, _ := c.FormFile("image")

	userID, _ := c.Get("userID")

	fmt.Printf("Creating product with category ID: %d\n", categoryID) // Debug log
	fmt.Printf("Received price: %v\n", price)                         // Debug log
	fmt.Printf("Received cost price: %v\n", costPrice)                // Debug log
	fmt.Printf("Received stock: %d\n", stock)                         // Debug log
	fmt.Printf("Received min stock: %d\n", minStock)                  // Debug log

	product, err := h.service.CreateProduct(userID.(uint), name, description, sku, price, costPrice, stock, minStock, uint(categoryID), file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

// GetProducts godoc
// @Summary List all Products
// @Tags Products
// @Produce json
// @Success 200 {array} domain.Product
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.service.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary Get a Product by ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} domain.Product
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update a Product
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} domain.Product
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")
	sku := c.PostForm("sku")
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
	costPrice, _ := strconv.ParseFloat(c.PostForm("cost_price"), 64)
	stock, _ := strconv.Atoi(c.PostForm("stock"))
	minStock, _ := strconv.Atoi(c.PostForm("min_stock"))
	categoryID, _ := strconv.Atoi(c.PostForm("category_id"))
	file, _ := c.FormFile("image")

	userID, _ := c.Get("userID")

	product, err := h.service.UpdateProduct(userID.(uint), uint(id), name, description, sku, price, costPrice, stock, minStock, uint(categoryID), file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a Product
// @Tags Products
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	if err := h.service.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// ─── Stock & Movements ─────────────────────────────────────────────────────

// AdjustStockInput defines the JSON body for stock adjustment.
type AdjustStockInput struct {
	Type     domain.MovementType `json:"type" binding:"required"`
	Quantity int                 `json:"quantity" binding:"required,min=1"`
	Note     string              `json:"note"`
}

// AdjustStock godoc
// @Summary Adjust product stock
// @Description Register a stock entry, exit or manual adjustment (Admin/Seller)
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param input body AdjustStockInput true "Movement data"
// @Success 201 {object} domain.InventoryMovement
// @Router /products/{id}/stock [post]
func (h *ProductHandler) AdjustStock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	var input AdjustStockInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	movement, err := h.service.AdjustStock(userID.(uint), uint(id), input.Type, input.Quantity, input.Note)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, movement)
}

// GetLowStock godoc
// @Summary List products with low stock
// @Tags Inventory
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Product
// @Router /products/low-stock [get]
func (h *ProductHandler) GetLowStock(c *gin.Context) {
	products, err := h.service.GetLowStockProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetMovements godoc
// @Summary List inventory movements for a product
// @Tags Inventory
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {array} domain.InventoryMovement
// @Router /products/{id}/movements [get]
func (h *ProductHandler) GetMovements(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	uid := uint(id)
	movements, err := h.service.GetMovements(&uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movements)
}

// GetAllMovements godoc
// @Summary List all inventory movements
// @Tags Inventory
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.InventoryMovement
// @Router /inventory/movements [get]
func (h *ProductHandler) GetAllMovements(c *gin.Context) {
	movements, err := h.service.GetMovements(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movements)
}

// GetFinalPrice godoc
// @Summary Get product final price with applied discount
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Router /products/{id}/price [get]
func (h *ProductHandler) GetFinalPrice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	finalPrice, discount, err := h.discountService.GetFinalPrice(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	product, _ := h.service.GetProductByID(uint(id))

	response := gin.H{
		"product_id":     id,
		"original_price": product.Price,
		"final_price":    finalPrice,
		"discount":       discount,
	}
	c.JSON(http.StatusOK, response)
}
