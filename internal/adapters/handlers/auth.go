package handlers

import (
	"net/http"
	"tienda-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Cedula    string `json:"cedula" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

// Register godoc
// @Summary Register a new customer
// @Description Creates a new customer account with their details
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body RegisterInput true "Customer Registration Data"
// @Success 201 {object} map[string]interface{} "Customer registered successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal error or already exists"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.RegisterCustomer(input.FirstName, input.LastName, input.Cedula, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Customer registered successfully", "user": gin.H{
		"id":         user.ID,
		"first_name": user.CustomerData.FirstName,
		"last_name":  user.CustomerData.LastName,
		"email":      user.Email,
		"role":       user.Role,
	}})
}

type SetupPasswordInput struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) SetupPassword(c *gin.Context) {
	var input SetupPasswordInput
	if err := h.checkInput(c, &input); err != nil {
		return
	}

	if err := h.authService.SetupPassword(input.Token, input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña configurada con éxito"})
}

type GoogleAuthInput struct {
	Token string `json:"token" binding:"required"`
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var input GoogleAuthInput
	if err := h.checkInput(c, &input); err != nil {
		return
	}

	loginRes, err := h.authService.GoogleLogin(input.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loginRes)
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary User Login
// @Description Authenticates a user and returns Access/Refresh tokens plus user details
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body LoginInput true "Login Credentials"
// @Success 200 {object} domain.LoginResponse "Returns login response with tokens and user details"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginRes, err := h.authService.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Handle MFA case
	if loginRes.MfaToken == "MFA_REQUIRED" {
		c.JSON(http.StatusAccepted, loginRes) // 202 Accepted
		return
	}

	c.JSON(http.StatusOK, loginRes)
}

type Setup2FAResponse struct {
	Secret    string `json:"secret"`
	QRCodeURL string `json:"qr_code_url"`
}

func (h *AuthHandler) Setup2FA(c *gin.Context) {
	userID := c.GetUint("userID")
	secret, qrURL, err := h.authService.Setup2FA(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Setup2FAResponse{
		Secret:    secret,
		QRCodeURL: qrURL,
	})
}

type Activate2FAInput struct {
	Code string `json:"code" binding:"required"`
}

func (h *AuthHandler) Activate2FA(c *gin.Context) {
	var input Activate2FAInput
	if err := h.checkInput(c, &input); err != nil {
		return
	}

	userID := c.GetUint("userID")
	if err := h.authService.Activate2FA(userID, input.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA activado con éxito"})
}

type Verify2FAInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

func (h *AuthHandler) Verify2FA(c *gin.Context) {
	var input Verify2FAInput
	if err := h.checkInput(c, &input); err != nil {
		return
	}

	loginRes, err := h.authService.FinalizeLogin(input.Email, input.Password, input.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loginRes)
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh godoc
// @Summary Refresh Access Token
// @Description Uses a refresh token to obtain a new pair of Access and Refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body RefreshInput true "Refresh Token"
// @Success 200 {object} domain.TokenPair "Returns new token pair"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized or expired"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var input RefreshInput
	if err := h.checkInput(c, &input); err != nil {
		return
	}

	tokenPair, err := h.authService.RefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

type VerifyEmailInput struct {
	Token string `json:"token" binding:"required"`
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var input VerifyEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.VerifyEmail(input.Token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Correo electrónico verificado con éxito"})
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var input ForgotPasswordInput
	if err := h.checkInput(c, &input); err != nil {
		return
	}

	if err := h.authService.RequestPasswordReset(input.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Si el correo existe, se ha enviado un enlace de recuperación"})
}

type ResetPasswordInput struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var input ResetPasswordInput
	if err := h.checkInput(c, &input); err != nil {
		return
	}

	if err := h.authService.ResetPassword(input.Token, input.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada con éxito"})
}

// Helper to avoid repetition
func (h *AuthHandler) checkInput(c *gin.Context, input interface{}) error {
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	return nil
}
