package patient_service

import (
	"agnos-backend/internal/models"
	"agnos-backend/internal/port"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type patientService struct {
	patientRepo  port.IPatientRepository
	hospitalRepo port.IHospitalRepository
	hisClient    port.IHisClient
}

func New(
	patientRepo port.IPatientRepository,
	hospitalRepo port.IHospitalRepository,
	hisClient port.IHisClient,
) port.IPatientService {
	return &patientService{
		patientRepo:  patientRepo,
		hospitalRepo: hospitalRepo,
		hisClient:    hisClient,
	}
}

func (s *patientService) Search(hospitalID uuid.UUID, filter models.SearchRequest) ([]models.Patient, error) {
	if filter.NationalID != "" || filter.PassportID != "" {
		searchID := filter.NationalID
		if searchID == "" {
			searchID = filter.PassportID
		}

		hospital, err := s.hospitalRepo.FindByID(hospitalID)
		if err != nil {
			log.Printf("search patient: find hospital error hospitalID=%s: %v", hospitalID, err)
			return nil, fmt.Errorf("find hospital: %w", err)
		}

		hisPatient, err := s.hisClient.FetchPatient(hospital.APIBase, searchID)
		if err != nil {
			log.Printf("search patient: his client error searchID=%s: %v", searchID, err)
			return nil, fmt.Errorf("fetch patient: %w", err)
		}

		patientData := models.Patient{
			HospitalID:   hospitalID,
			FirstNameTH:  hisPatient.FirstNameTH,
			MiddleNameTH: hisPatient.MiddleNameTH,
			LastNameTH:   hisPatient.LastNameTH,
			FirstNameEN:  hisPatient.FirstNameEN,
			MiddleNameEN: hisPatient.MiddleNameEN,
			LastNameEN:   hisPatient.LastNameEN,
			DateOfBirth:  hisPatient.DateOfBirth,
			PatientHN:    hisPatient.PatientHN,
			NationalID:   hisPatient.NationalID,
			PassportID:   hisPatient.PassportID,
			PhoneNumber:  hisPatient.PhoneNumber,
			Email:        hisPatient.Email,
			Gender:       hisPatient.Gender,
		}

		if err := s.patientRepo.Upsert(patientData); err != nil {
			log.Printf("search patient: upsert error searchID=%s: %v", searchID, err)
		}
	}

	patients, err := s.patientRepo.Search(hospitalID, filter)
	if err != nil {
		log.Printf("search patient: db search error hospitalID=%s: %v", hospitalID, err)
		return nil, err
	}

	return patients, nil
}
