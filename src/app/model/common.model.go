package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;size:191"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:datetime;autoCreateTime:nano"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:datetime;autoUpdateTime:nano"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;type:datetime"`
}

func (m *Base) BeforeCreate(_ *gorm.DB) error {
	m.ID = uuid.New()

	return nil
}
