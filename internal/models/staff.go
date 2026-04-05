package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateStaffRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"`
}

type CreateStaffResponse struct {
	StaffID    string `json:"staff_id"`
	HospitalID string `json:"hospital_id"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token" `
}

type Staff struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Username   string    `gorm:"type:varchar(255);not null" json:"username"`
	Password   string    `gorm:"type:varchar(512);not null" json:"-"`
	HospitalID uuid.UUID `gorm:"type:uuid;not null;index" json:"hospital_id"`
	Hospital   Hospital  `gorm:"foreignKey:HospitalID" json:"hospital,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (s *Staff) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (h *Staff) TableName() string {
	return "staffs"
}
