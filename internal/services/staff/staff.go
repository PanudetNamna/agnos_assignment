package staff_service

import (
	"agnos-backend/internal/config"
	"agnos-backend/internal/models"
	"agnos-backend/internal/port"
	"agnos-backend/internal/utility"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type staffService struct {
	staffRepo    port.IStaffRepository
	hospitalRepo port.IHospitalRepository
	cfg          *config.AppConfig
}

func New(staffRepo port.IStaffRepository, hospitalRepo port.IHospitalRepository, cfg *config.AppConfig) port.IStaffService {
	return &staffService{
		staffRepo:    staffRepo,
		hospitalRepo: hospitalRepo,
		cfg:          cfg,
	}
}

func (s *staffService) Create(username, password, hospitalName string) (*models.Staff, error) {
	hospital, err := s.hospitalRepo.FindByName(hospitalName)
	if err != nil {
		log.Printf("create staff: hospital not found: %v", err)
		return nil, errors.New("hospital not found")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("create staff: hash password error: %v", err)
		return nil, err
	}

	staff := &models.Staff{
		Username:   username,
		Password:   string(hashed),
		HospitalID: hospital.ID,
	}

	if err := s.staffRepo.Create(staff); err != nil {
		log.Printf("create staff: insert db error: %v", err)
		return nil, err
	}

	return staff, nil
}

func (s *staffService) Login(username, password, hospitalName string) (string, error) {
	hospital, err := s.hospitalRepo.FindByName(hospitalName)
	if err != nil {
		log.Printf("login: hospital not found: %v", err)
		return "", errors.New("hospital not found")
	}

	staff, err := s.staffRepo.FindByUsernameAndHospital(username, hospital.ID)
	if err != nil {
		log.Printf("login: staff not found username=%s hospital=%s: %v", username, hospitalName, err)
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(password)); err != nil {
		log.Printf("login: wrong password username=%s: %v", username, err)
		return "", errors.New("invalid credentials")
	}

	signed, err := utility.GenerateToken(staff.ID, hospital.ID, s.cfg.Secrets.JwtSecretKey)
	if err != nil {
		log.Printf("login: generate token error: %v", err)
		return "", errors.New("generate token failed")
	}

	return signed, nil
}
