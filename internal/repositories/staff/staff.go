package staff_repository

import (
	"agnos-backend/internal/models"
	"agnos-backend/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type staffRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) port.IStaffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) Create(staff *models.Staff) error {
	return r.db.Create(staff).Error
}

func (r *staffRepository) FindByUsernameAndHospital(username string, hospitalID uuid.UUID) (*models.Staff, error) {
	var staff models.Staff
	err := r.db.Where("username = ? AND hospital_id = ?", username, hospitalID).First(&staff).Error
	if err != nil {
		return nil, err
	}
	return &staff, nil
}
