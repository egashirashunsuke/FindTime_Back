package model

type UserBandDTO struct {
	BandID     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	BandName   string `json:"name"`
	IsFavorite bool   `json:"is_favorite"`
}
