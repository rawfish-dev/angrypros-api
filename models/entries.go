package models

import "time"

type AngerTier struct {
	Id        int64
	Label     string `gorm:"not null;uniqueindex"`
	RageLevel int    // Keep order explicit to prevent reliance on ID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Entry struct {
	Id          int64
	TextContent string
	RageLevel   int
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// References
	UserId               int64     `gorm:"not null"`
	User                 User      `gorm:"foreignKey:UserId"`
	CountryIsoAlpha2Code string    `gorm:"not null"`
	Country              Country   `gorm:"foreignKey:CountryIsoAlpha2Code"`
	AngryTierId          int64     `gorm:"not null"`
	AngerTier            AngerTier `gorm:"foreignKey:AngryTierId"`
}
