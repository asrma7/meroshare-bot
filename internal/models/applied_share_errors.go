package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppliedShareError struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index"`
	AppliedShareID uuid.UUID `gorm:"type:uuid;not null;index"`
	Message        string    `gorm:"not null"`
	Seen           bool      `gorm:"default:false"`
	CreatedAt      time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt      time.Time `gorm:"type:timestamptz;default:now()"`
}

func (u *AppliedShareError) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
