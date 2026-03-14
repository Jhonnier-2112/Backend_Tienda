package repository

import (
	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"

	"gorm.io/gorm"
)

type emailLogRepository struct {
	db *gorm.DB
}

func NewEmailLogRepository(db *gorm.DB) ports.EmailLogRepository {
	return &emailLogRepository{db: db}
}

func (r *emailLogRepository) Create(log *domain.EmailLog) error {
	return r.db.Create(log).Error
}

func (r *emailLogRepository) GetAll() ([]domain.EmailLog, error) {
	var logs []domain.EmailLog
	err := r.db.Order("created_at desc").Find(&logs).Error
	return logs, err
}
