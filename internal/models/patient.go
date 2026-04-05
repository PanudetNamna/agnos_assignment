package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SearchRequest struct {
	NationalID  string `json:"national_id"`
	PassportID  string `json:"passport_id"`
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type SearchResponse struct {
	Count    int       `json:"count"`
	Patients []Patient `json:"patients"`
}

type PatientSearchFilter struct {
	NationalID  string
	PassportID  string
	FirstName   string
	MiddleName  string
	LastName    string
	DateOfBirth string
	PhoneNumber string
	Email       string
}

type Patient struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	HospitalID   uuid.UUID `gorm:"type:uuid;not null;index" json:"hospital_id"`
	Hospital     Hospital  `gorm:"foreignKey:HospitalID" json:"hospital,omitempty"`
	FirstNameTH  string    `gorm:"type:varchar(255)" json:"first_name_th"`
	MiddleNameTH string    `gorm:"type:varchar(255)" json:"middle_name_th"`
	LastNameTH   string    `gorm:"type:varchar(255)" json:"last_name_th"`
	FirstNameEN  string    `gorm:"type:varchar(255)" json:"first_name_en"`
	MiddleNameEN string    `gorm:"type:varchar(255)" json:"middle_name_en"`
	LastNameEN   string    `gorm:"type:varchar(255)" json:"last_name_en"`
	DateOfBirth  string    `gorm:"type:varchar(10)" json:"date_of_birth"`
	PatientHN    string    `gorm:"type:varchar(100)" json:"patient_hn"`
	NationalID   string    `gorm:"type:varchar(20);index" json:"national_id"`
	PassportID   string    `gorm:"type:varchar(50);index" json:"passport_id"`
	PhoneNumber  string    `gorm:"type:varchar(20)" json:"phone_number"`
	Email        string    `gorm:"type:varchar(255)" json:"email"`
	Gender       string    `gorm:"type:char(1)" json:"gender"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (p *Patient) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (h *Patient) TableName() string {
	return "patients"
}
