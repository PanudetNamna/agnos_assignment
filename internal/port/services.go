package port

import (
	"agnos-backend/internal/models"

	"github.com/google/uuid"
)

//go:generate mockgen -source=services.go -destination=./mocks/mock_services.go -package=mocks

type IPatientService interface {
	Search(hospitalID uuid.UUID, filter models.SearchRequest) ([]models.Patient, error)
}

type IStaffService interface {
	Create(username, password, hospitalName string) (*models.Staff, error)
	Login(username, password, hospitalName string) (string, error)
}
