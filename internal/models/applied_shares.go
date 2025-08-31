package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppliedShare struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index"`
	AccountID      uuid.UUID `gorm:"type:uuid;not null;index"`
	CompanyName    string    `gorm:"not null"`
	CompanyShareID uint16    `gorm:"not null"`
	Scrip          string    `gorm:"not null"`
	AppliedKitta   string    `gorm:"not null"`
	ShareGroupName string    `gorm:"not null"`
	ShareTypeName  string    `gorm:"not null"`
	SubGroup       string    `gorm:"not null"`
	Status         string    `gorm:"type:varchar(20);default:'applied'"`
	CreatedAt      time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt      time.Time `gorm:"type:timestamptz;default:now()"`
}

func (u *AppliedShare) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
