package repository

import (
	"tienda-backend/internal/core/domain"

	"gorm.io/gorm"
)

type UserPostgresRepository struct {
	db *gorm.DB
}

func NewUserPostgresRepository(db *gorm.DB) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

func (r *UserPostgresRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserPostgresRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Preload("CustomerData").Preload("SellerData").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserPostgresRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.Preload("CustomerData").Preload("SellerData").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserPostgresRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *UserPostgresRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Preload("CustomerData").Preload("SellerData").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserPostgresRepository) FindByVerificationToken(token string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("verification_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserPostgresRepository) FindByResetToken(token string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("reset_password_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
