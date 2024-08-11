package model

import (
	_ "gorm.io/gorm"
)

type Band struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
}

func AddBand(band *Band) error {

	if err := db.Create(band).Error; err != nil {
		return err
	}
	return nil

}

func GetBandsByUserID(UserID int) ([]Band, error) {
	var bands []Band

	err := db.Joins("join user_bands on user_bands.band_id = bands.id").
		Where("user_bands.user_id = ?", UserID).Find(&bands).Error
	if err != nil {
		return nil, err
	}
	return bands, nil
}

func GetUserBandWithFavorite(UserID int) ([]UserBandDTO, error) {
	var dtos []UserBandDTO

	err := db.Table("bands").
		Select("bands.id as band_id, bands.name as band_name, user_bands.user_id as user_id, user_bands.is_favorite").
		Joins("join user_bands on user_bands.band_id = bands.id").
		Where("user_bands.user_id = ?", UserID).
		Order("user_bands.is_favorite DESC").
		Scan(&dtos).Error

	if err != nil {
		return nil, err
	}

	return dtos, nil
}
