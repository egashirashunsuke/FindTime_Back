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
