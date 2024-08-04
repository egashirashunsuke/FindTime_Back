package router

import (
	"FindTime-Server/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func AddBandMemberHandler(c echo.Context) error {
	uid := userIDFromToken(c)
	if user := model.FindUser(&model.User{ID: uid}); user.ID == 0 {
		return echo.ErrNotFound
	}

	bandID, err := strconv.Atoi(c.Param("bandID"))
	if err != nil {
		return echo.ErrNotFound
	}

	member := &model.UserBand{
		BandID: bandID,
		UserID: uid,
	}

	if err := model.AddBandMember(member); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, member)

}

func GetBandMembersHandler(c echo.Context) error {
	uid := userIDFromToken(c)
	if user := model.FindUser(&model.User{ID: uid}); user.ID == 0 {
		return echo.ErrNotFound
	}

	bandID, err := strconv.Atoi(c.Param("bandID"))
	if err != nil {
		return echo.ErrNotFound
	}

	members, err := model.GetBandMembers(bandID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, members)

}

func LeaveBandMemberHandler(c echo.Context) error {
	uid := userIDFromToken(c)
	if user := model.FindUser(&model.User{ID: uid}); user.ID == 0 {
		return echo.ErrNotFound
	}

	bandID, err := strconv.Atoi(c.Param("bandID"))
	if err != nil {
		return echo.ErrNotFound
	}

	if err := model.DeleteBandMember(uid, bandID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)

}
