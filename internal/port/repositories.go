package port

import (
	"agnos-backend/internal/models"

	"github.com/google/uuid"
)

//go:generate mockgen -source=repositories.go -destination=./mocks/mock_repositories.go -package=mocks

type IPatientRepository interface {
	Search(hospitalID uuid.UUID, filter models.SearchRequest) ([]models.Patient, error)
	Upsert(data models.Patient) error
}

type IStaffRepository interface {
	Create(staff *models.Staff) error
	FindByUsernameAndHospital(username string, hospitalID uuid.UUID) (*models.Staff, error)
}

type IHospitalRepository interface {
	FindByName(name string) (*models.Hospital, error)
	FindByID(id uuid.UUID) (*models.Hospital, error)
}
