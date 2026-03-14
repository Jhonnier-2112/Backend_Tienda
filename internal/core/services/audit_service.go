package services

import (
	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/core/ports"
)

type auditLogService struct {
	repo ports.AuditLogRepository
}

func NewAuditLogService(repo ports.AuditLogRepository) ports.AuditLogService {
	return &auditLogService{repo: repo}
}

func (s *auditLogService) LogAction(userID uint, action, entityName string, entityID uint, oldValue, newValue string) error {
	log := &domain.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityName: entityName,
		EntityID:   entityID,
		OldValue:   oldValue,
		NewValue:   newValue,
	}
	return s.repo.Create(log)
}

func (s *auditLogService) GetProductHistory(productID uint) ([]domain.AuditLog, error) {
	return s.repo.GetByEntity("Product", productID)
}

func (s *auditLogService) GetGlobalHistory() ([]domain.AuditLog, error) {
	return s.repo.GetAll()
}
