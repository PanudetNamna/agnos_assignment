package port

import "agnos-backend/internal/models"

//go:generate mockgen -source=his_client.go -destination=./mocks/mock_his_client.go -package=mocks

type IHisClient interface {
	FetchPatient(apiBase, id string) (*models.HISPatientResponse, error)
}
