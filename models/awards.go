package models

type Award struct {
	Id          int64
	Label       string `gorm:"not null;uniqueindex"`
	ImageUrl    string `gorm:"not null;uniqueindex" `
	Description string `gorm:"not null"`

	// References
	RarityId int64  `gorm:"not null"`
	Rarity   Rarity `gorm:"foreignKey:RarityId"`
}

type Rarity struct {
	Id               int64
	Label            string  `gorm:"not null;uniqueindex"`
	PercentageChance float64 `gorm:"not null;uniqueindex"`
}
