package models

import (
	"time"
)

type User struct {
	Id                     int64
	FirebaseUserId         string `gorm:"uniqueindex;not null"`
	Title                  string `gorm:"not null"`
	NormalisedEmailAddress string `gorm:"uniqueindex;not null"`
	CreatedAt              time.Time
	UpdatedAt              time.Time

	// References
	CountryIsoAlpha2Code string  `gorm:"not null"`
	Country              Country `gorm:"foreignKey:CountryIsoAlpha2Code"`
}

type Country struct {
	IsoAlpha2Code string `gorm:"primaryKey"`
	Name          string `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
