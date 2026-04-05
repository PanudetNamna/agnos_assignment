package patient_repository

import (
	"agnos-backend/internal/models"
	"agnos-backend/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) port.IPatientRepository {
	return &repository{db: db}
}

func (r *repository) Search(hospitalID uuid.UUID, filter models.SearchRequest) ([]models.Patient, error) {
	query := r.db.Where("hospital_id = ?", hospitalID)

	if filter.NationalID != "" {
		query = query.Where("national_id = ?", filter.NationalID)
	}
	if filter.PassportID != "" {
		query = query.Where("passport_id = ?", filter.PassportID)
	}
	if filter.FirstName != "" {
		query = query.Where("first_name_en ILIKE ? OR first_name_th ILIKE ?",
			"%"+filter.FirstName+"%", "%"+filter.FirstName+"%")
	}
	if filter.MiddleName != "" {
		query = query.Where("middle_name_en ILIKE ? OR middle_name_th ILIKE ?",
			"%"+filter.MiddleName+"%", "%"+filter.MiddleName+"%")
	}
	if filter.LastName != "" {
		query = query.Where("last_name_en ILIKE ? OR last_name_th ILIKE ?",
			"%"+filter.LastName+"%", "%"+filter.LastName+"%")
	}
	if filter.DateOfBirth != "" {
		query = query.Where("date_of_birth = ?", filter.DateOfBirth)
	}
	if filter.PhoneNumber != "" {
		query = query.Where("phone_number = ?", filter.PhoneNumber)
	}
	if filter.Email != "" {
		query = query.Where("email ILIKE ?", "%"+filter.Email+"%")
	}

	var patients []models.Patient
	if err := query.Find(&patients).Error; err != nil {
		return nil, err
	}
	return patients, nil
}

func (r *repository) Upsert(data models.Patient) error {
	return r.db.
		Where(models.Patient{
			HospitalID: data.HospitalID,
			NationalID: data.NationalID,
			PassportID: data.PassportID,
		}).
		Assign(data).
		FirstOrCreate(&data).Error
}
