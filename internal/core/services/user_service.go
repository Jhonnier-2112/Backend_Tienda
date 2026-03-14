package services

import (
	"errors"
	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"
	"tienda-backend/pkg/utils"
)

type userService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) ports.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateSeller(companyName, contactName, nit, email, password string) (*domain.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &domain.User{
		Email:    email,
		Password: hashedPassword,
		Role:     "seller",
		SellerData: &domain.SellerProfile{
			CompanyName: companyName,
			ContactName: contactName,
			NIT:         nit,
		},
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create seller, possibly email or NIT already exists")
	}

	return user, nil
}

func (s *userService) CreateAdmin(name, email, password string) (*domain.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &domain.User{
		Email:    email,
		Password: hashedPassword,
		Role:     "admin",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create admin, possibly email already exists")
	}

	return user, nil
}

func (s *userService) GetAllUsers() ([]domain.User, error) {
	return s.userRepo.GetAll()
}
