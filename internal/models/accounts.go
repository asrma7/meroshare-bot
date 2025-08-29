package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID             uuid.UUID `gorm:"type:uuid;not null;index"`
	Name               string    `gorm:"not null"`
	Email              string    `gorm:"not null"`
	Contact            string    `gorm:"not null"`
	ClientID           uint16    `gorm:"not null"`
	Username           string    `gorm:"uniqueIndex;not null;type:varchar(50)"`
	Password           string    `gorm:"not null"`
	BankID             string    `gorm:"not null"`
	CRNNumber          string    `gorm:"not null"`
	TransactionPIN     string    `gorm:"not null"`
	AccountTypeId      uint8     `gorm:"not null"`
	PreferredKitta     string    `gorm:"not null"`
	Demat              string    `gorm:"not null"`
	BOID               string    `gorm:"not null"`
	AccountNumber      string    `gorm:"not null"`
	CustomerId         uint32    `gorm:"not null"`
	AccountBranchId    uint32    `gorm:"not null"`
	DMATExpiryDate     string    `gorm:"type:varchar(50);not null"`
	ExpiredDate        time.Time `gorm:"type:timestamptz;not null"`
	PasswordExpiryDate time.Time `gorm:"type:timestamptz;not null"`
	CreatedAt          time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt          time.Time `gorm:"type:timestamptz;default:now()"`
	Status             string    `gorm:"type:varchar(20);default:'active'"`
}

func (u *Account) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
