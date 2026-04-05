package hospital_repository

import (
	"agnos-backend/internal/models"
	"agnos-backend/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) port.IHospitalRepository {
	return &repository{db: db}
}

func (r *repository) FindByName(name string) (*models.Hospital, error) {
	var hospital models.Hospital
	err := r.db.Where("name = ?", name).First(&hospital).Error
	if err != nil {
		return nil, err
	}
	return &hospital, nil
}

func (r *repository) FindByID(id uuid.UUID) (*models.Hospital, error) {
	var hospital models.Hospital
	err := r.db.Where("id = ?", id).First(&hospital).Error
	if err != nil {
		return nil, err
	}
	return &hospital, nil
}
