package repository

import (
	"tienda-backend/internal/core/domain"

	"gorm.io/gorm"
)

type AuditLogPostgresRepository struct {
	db *gorm.DB
}

func NewAuditLogPostgresRepository(db *gorm.DB) *AuditLogPostgresRepository {
	return &AuditLogPostgresRepository{db: db}
}

func (r *AuditLogPostgresRepository) Create(log *domain.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditLogPostgresRepository) GetByEntity(entityName string, entityID uint) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	err := r.db.Preload("User").
		Where("entity_name = ? AND entity_id = ?", entityName, entityID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

func (r *AuditLogPostgresRepository) GetAll() ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	err := r.db.Preload("User").Order("created_at DESC").Find(&logs).Error
	return logs, err
}
