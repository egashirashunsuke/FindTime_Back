package router

import (
	"FindTime-Server/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AddBandHandler(c echo.Context) error {
	uid := userIDFromToken(c)
	if user := model.FindUser(&model.User{ID: uid}); user.ID == 0 {
		return echo.ErrNotFound
	}

	band := new(model.Band)

	if err := c.Bind(band); err != nil {
		return err
	}

	if err := model.AddBand(band); err != nil {
		return err
	}

	member := &model.UserBand{
		BandID: band.ID,
		UserID: uid,
	}
	if err := model.AddBandMember(member); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, band)

}

func GetBandHandler(c echo.Context) error {
	uid := userIDFromToken(c)
	if user := model.FindUser(&model.User{ID: uid}); user.ID == 0 {
		return echo.ErrNotFound
	}

	bands, err := model.GetUserBandWithFavorite(uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, bands)

}
