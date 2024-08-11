package model

type UserBand struct {
	ID         int `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     int
	BandID     int
	IsFavorite bool `gorm:"default:false"`
}

func AddBandMember(member *UserBand) error {
	if err := db.Create(member).Error; err != nil {
		return err
	}
	return nil
}

func GetBandMembers(BandID int) ([]User, error) {
	var members []User

	err := db.Joins("join user_bands on user_bands.user_id = users.id").
		Where("user_bands.band_id = ?", BandID).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil

}

func DeleteBandMember(uid int, bandId int) error {
	err := db.Where("user_bands.user_id = ? AND user_bands.band_id = ?", uid, bandId).
		Delete(&UserBand{}).Error
	return err
}

func FavoriteBand(uid int, bandId int) error {
	err := db.Model(&UserBand{}).
		Where("user_bands.user_id = ? AND user_bands.band_id = ?", uid, bandId).
		Update("is_favorite", true).Error
	return err
}

func RemoveFavoriteBand(uid int, bandId int) error {
	err := db.Model(&UserBand{}).
		Where("user_bands.user_id = ? AND user_bands.band_id = ?", uid, bandId).
		Update("is_favorite", false).Error
	return err
}
