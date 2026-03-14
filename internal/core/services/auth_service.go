package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"
	"tienda-backend/pkg/utils"

	"github.com/pquerna/otp/totp"
)

type authService struct {
	userRepo         ports.UserRepository
	refreshTokenRepo ports.RefreshTokenRepository
	emailService     ports.EmailService
	jwtSecret        string
}

func NewAuthService(
	userRepo ports.UserRepository,
	refreshTokenRepo ports.RefreshTokenRepository,
	emailService ports.EmailService,
	jwtSecret string,
) ports.AuthService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		emailService:     emailService,
		jwtSecret:        jwtSecret,
	}
}

func (s *authService) generateRandomToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *authService) RegisterCustomer(firstName, lastName, cedula, email string) (*domain.User, error) {
	verificationToken := s.generateRandomToken()

	user := &domain.User{
		Name:  firstName + " " + lastName,
		Email: email,
		// Password field is empty initially for standard registration until SetupPassword
		Role:              "customer",
		IsVerified:        false,
		VerificationToken: verificationToken,
		CustomerData: &domain.CustomerProfile{
			FirstName: firstName,
			LastName:  lastName,
			Cedula:    cedula,
		},
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user, possibly email or cedula already exists")
	}

	// Send verification email
	go s.emailService.SendVerificationEmail(user.Email, verificationToken)

	return user, nil
}

func (s *authService) Login(email, password string) (*domain.LoginResponse, error) {
	log.Printf("[DEBUG] Login attempt for email: %s", email)
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("[DEBUG] User not found: %s", email)
		return nil, errors.New("invalid email or password")
	}

	// 1. Check if account is locked
	if user.LockUntil != nil && user.LockUntil.After(time.Now()) {
		log.Printf("[DEBUG] Account locked for: %s", email)
		return nil, errors.New("tu cuenta está bloqueada temporalmente por demasiados intentos fallidos. Inténtalo más tarde")
	}

	// 2. Check if verified
	if !user.IsVerified {
		log.Printf("[DEBUG] User not verified: %s", email)
		return nil, errors.New("por favor verifica tu correo electrónico antes de iniciar sesión")
	}

	log.Println("[DEBUG] Checking password hash...")
	if !utils.CheckPasswordHash(password, user.Password) {
		log.Println("[DEBUG] Password hash mismatch")

		// Handle failed attempts
		user.FailedAttempts++
		if user.FailedAttempts >= 5 {
			lockTime := time.Now().Add(15 * time.Minute)
			user.LockUntil = &lockTime
		}
		s.userRepo.Update(user)

		return nil, errors.New("invalid email or password")
	}

	// 3. Check if 2FA is enabled
	if user.TwoFactorEnabled {
		log.Printf("[DEBUG] 2FA required for user: %s", email)
		// We return a specific structure or error that indicates MFA is needed
		// For now, let's return a partial response or a specific error message
		return &domain.LoginResponse{
			UserID:   user.ID,
			Email:    user.Email,
			Role:     user.Role,
			MfaToken: "MFA_REQUIRED", // Signal to frontend
		}, nil
	}

	// 4. Reset failed attempts on success
	if user.FailedAttempts > 0 {
		user.FailedAttempts = 0
		user.LockUntil = nil
		s.userRepo.Update(user)
	}

	log.Println("[DEBUG] Generating token pair...")
	tokenPair, err := utils.GenerateTokenPair(user, s.jwtSecret)
	if err != nil {
		log.Printf("[DEBUG] Token generation failed: %v", err)
		return nil, errors.New("failed to generate tokens")
	}

	// Store refresh token
	claims, _ := utils.ValidateToken(tokenPair.RefreshToken, s.jwtSecret)

	rt := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: claims.ExpiresAt.Time,
	}

	if err := s.refreshTokenRepo.Create(rt); err != nil {
		return nil, errors.New("failed to store refresh token")
	}

	return &domain.LoginResponse{
		UserID:       user.ID,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		Email:        user.Email,
		Role:         user.Role,
		Customer:     user.CustomerData,
		Seller:       user.SellerData,
	}, nil
}

func (s *authService) SetupPassword(token, password string) error {
	user, err := s.userRepo.FindByVerificationToken(token)
	if err != nil {
		return errors.New("token de configuración inválido o expirado")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return errors.New("error al procesar la contraseña")
	}

	user.Password = hashedPassword
	user.IsVerified = true
	user.VerificationToken = ""

	return s.userRepo.Update(user)
}

func (s *authService) GoogleLogin(googleToken string) (*domain.LoginResponse, error) {
	// 1. Verify Google token via Google API
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + googleToken)
	if err != nil {
		return nil, errors.New("error validando el token de Google")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("token de Google inválido")
	}

	var googleObj map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&googleObj); err != nil {
		return nil, errors.New("error decodificando respuesta de Google")
	}

	email, ok := googleObj["email"].(string)
	if !ok || email == "" {
		return nil, errors.New("no se pudo obtener el email de Google")
	}

	// 2. Find or Create User
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Create new user if not exists
		firstName := ""
		lastName := ""
		if name, ok := googleObj["given_name"].(string); ok {
			firstName = name
		}
		if name, ok := googleObj["family_name"].(string); ok {
			lastName = name
		}

		user = &domain.User{
			Name:       firstName + " " + lastName,
			Email:      email,
			Role:       "customer",
			IsVerified: true, // Auto-verified by Google
			CustomerData: &domain.CustomerProfile{
				FirstName: firstName,
				LastName:  lastName,
				Cedula:    "GOOGLE-" + s.generateRandomToken()[:8], // Placehoder for cedula
			},
		}

		if err := s.userRepo.Create(user); err != nil {
			return nil, errors.New("error registrando usuario desde Google")
		}
	} else if !user.IsVerified {
		// If user exists but is not verified, auto-verify them since Google confirmed their email
		user.IsVerified = true
		user.VerificationToken = ""
		s.userRepo.Update(user)
	}

	// 3. Check if account is locked
	if user.LockUntil != nil && user.LockUntil.After(time.Now()) {
		return nil, errors.New("tu cuenta está bloqueada temporalmente")
	}

	// 4. Generate tokens
	tokenPair, err := utils.GenerateTokenPair(user, s.jwtSecret)
	if err != nil {
		return nil, errors.New("error generando tokens")
	}

	// 5. Store refresh token
	claims, _ := utils.ValidateToken(tokenPair.RefreshToken, s.jwtSecret)
	rt := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: claims.ExpiresAt.Time,
	}

	if err := s.refreshTokenRepo.Create(rt); err != nil {
		return nil, errors.New("error guardando el refresh token")
	}

	return &domain.LoginResponse{
		UserID:       user.ID,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		Email:        user.Email,
		Role:         user.Role,
		Customer:     user.CustomerData,
		Seller:       user.SellerData,
	}, nil
}

func (s *authService) Setup2FA(userID uint) (string, string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", "", err
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "VirtualStore",
		AccountName: user.Email,
	})
	if err != nil {
		return "", "", err
	}

	// Store secret temporarily but don't enable yet
	user.TwoFactorSecret = key.Secret()
	if err := s.userRepo.Update(user); err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

func (s *authService) Activate2FA(userID uint, code string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if user.TwoFactorSecret == "" {
		return errors.New("2FA not set up")
	}

	valid := totp.Validate(code, user.TwoFactorSecret)
	if !valid {
		return errors.New("código de verificación inválido")
	}

	user.TwoFactorEnabled = true
	return s.userRepo.Update(user)
}

func (s *authService) Verify2FA(userID uint, code string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if !user.TwoFactorEnabled {
		return errors.New("2FA not enabled")
	}

	valid := totp.Validate(code, user.TwoFactorSecret)
	if !valid {
		return errors.New("código de verificación inválido")
	}

	return nil
}

func (s *authService) FinalizeLogin(email, password, code string) (*domain.LoginResponse, error) {
	log.Printf("[DEBUG] Finalizing login for email: %s", email)
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 1. Double check password (security measure)
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// 2. Verify MFA code
	valid := totp.Validate(code, user.TwoFactorSecret)
	if !valid {
		return nil, errors.New("código de autenticación inválido")
	}

	// 3. Reset failed attempts
	if user.FailedAttempts > 0 {
		user.FailedAttempts = 0
		user.LockUntil = nil
		s.userRepo.Update(user)
	}

	// 4. Generate tokens
	tokenPair, err := utils.GenerateTokenPair(user, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// Store refresh token
	claims, _ := utils.ValidateToken(tokenPair.RefreshToken, s.jwtSecret)
	rt := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: claims.ExpiresAt.Time,
	}
	s.refreshTokenRepo.Create(rt)

	return &domain.LoginResponse{
		UserID:       user.ID,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		Email:        user.Email,
		Role:         user.Role,
		Customer:     user.CustomerData,
		Seller:       user.SellerData,
	}, nil
}

func (s *authService) RefreshToken(token string) (*domain.TokenPair, error) {
	// 1. Validate the old token
	_, err := utils.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// 2. Check if it exists in DB and is not revoked
	rt, err := s.refreshTokenRepo.FindByToken(token)
	if err != nil {
		return nil, errors.New("refresh token not found or revoked")
	}

	// 3. Get user
	user, err := s.userRepo.FindByID(rt.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 4. Revoke the old token (Rotation)
	if err := s.refreshTokenRepo.Revoke(token); err != nil {
		return nil, errors.New("failed to revoke old token")
	}

	// 5. Generate new pair
	tokenPair, err := utils.GenerateTokenPair(user, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	// 6. Store new refresh token
	newClaims, err := utils.ValidateToken(tokenPair.RefreshToken, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to validate new refresh token")
	}

	newRt := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: newClaims.ExpiresAt.Time,
	}

	if err := s.refreshTokenRepo.Create(newRt); err != nil {
		return nil, errors.New("failed to store new refresh token")
	}

	return tokenPair, nil
}

func (s *authService) ValidateToken(tokenString string) (uint, string, error) {
	claims, err := utils.ValidateToken(tokenString, s.jwtSecret)
	if err != nil {
		return 0, "", errors.New("invalid or expired token")
	}

	return claims.UserID, claims.Role, nil
}

func (s *authService) VerifyEmail(token string) error {
	user, err := s.userRepo.FindByVerificationToken(token)
	if err != nil {
		return errors.New("token de verificación inválido o expirado")
	}

	user.IsVerified = true
	user.VerificationToken = ""
	return s.userRepo.Update(user)
}

func (s *authService) RequestPasswordReset(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// We don't want to leak if an email exists, but for this mock we can just return nil
		return nil
	}

	token := s.generateRandomToken()
	expires := time.Now().Add(1 * time.Hour)

	user.ResetPasswordToken = token
	user.ResetPasswordExpires = &expires

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	go s.emailService.SendPasswordResetEmail(user.Email, token)
	return nil
}

func (s *authService) ResetPassword(token, newPassword string) error {
	user, err := s.userRepo.FindByResetToken(token)
	if err != nil {
		return errors.New("token de recuperación inválido")
	}

	if user.ResetPasswordExpires == nil || user.ResetPasswordExpires.Before(time.Now()) {
		return errors.New("el token ha expirado")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.ResetPasswordToken = ""
	user.ResetPasswordExpires = nil
	user.FailedAttempts = 0
	user.LockUntil = nil

	return s.userRepo.Update(user)
}
