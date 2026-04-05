package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Hospital struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	APIBase   string    `gorm:"type:varchar(512)" json:"api_base"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *Hospital) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

func (h *Hospital) TableName() string {
	return "hospitals"
}
