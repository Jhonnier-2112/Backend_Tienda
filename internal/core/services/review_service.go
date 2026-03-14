package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"
)

type reviewService struct {
	reviewRepo   ports.ReviewRepository
	productRepo  ports.InventoryRepository
	imageStorage ports.ImageStorageService
}

func NewReviewService(reviewRepo ports.ReviewRepository, productRepo ports.InventoryRepository, imageStorage ports.ImageStorageService) ports.ReviewService {
	return &reviewService{
		reviewRepo:   reviewRepo,
		productRepo:  productRepo,
		imageStorage: imageStorage,
	}
}

func (s *reviewService) AddReview(productID, userID uint, rating int, comment string, imageFile *multipart.FileHeader) (*domain.Review, error) {
	// 1. Validar que el producto existe
	if _, err := s.productRepo.GetProductByID(productID); err != nil {
		return nil, errors.New("product not found")
	}

	// 2. Validar rating
	if rating < 1 || rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	var imageURL *string
	if imageFile != nil {
		// 3. Validar Seguridad de Imagen (Peso)
		// Límite: 2MB (2 * 1024 * 1024 bytes)
		if imageFile.Size > 2*1024*1024 {
			return nil, errors.New("image size exceeds 2MB limit")
		}

		// 4. Validar Seguridad de Imagen (Formato real via sniffing)
		src, err := imageFile.Open()
		if err != nil {
			return nil, fmt.Errorf("could not open image: %w", err)
		}
		defer src.Close()

		// Leer los primeros 512 bytes para detectar el tipo de contenido
		buffer := make([]byte, 512)
		if _, err := src.Read(buffer); err != nil {
			return nil, fmt.Errorf("could not read image header: %w", err)
		}

		contentType := http.DetectContentType(buffer)
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/webp": true,
		}

		if !allowedTypes[contentType] {
			return nil, fmt.Errorf("invalid image format: %s. Only JPEG, PNG, and WEBP are allowed", contentType)
		}

		// 5. Subir imagen
		uploadedURL, err := s.imageStorage.UploadImage(imageFile)
		if err != nil {
			return nil, fmt.Errorf("failed to upload review image: %w", err)
		}
		imageURL = &uploadedURL
	}

	review := &domain.Review{
		ProductID: productID,
		UserID:    userID,
		Rating:    rating,
		Comment:   comment,
		ImageURL:  imageURL,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *reviewService) GetProductReviews(productID uint) ([]domain.Review, error) {
	return s.reviewRepo.GetByProductID(productID)
}
