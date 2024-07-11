package router

import (
	"FindTime-Server/model"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func GetEventsHandler(c echo.Context) error {
	uid := userIDFromToken(c)
	if user := model.FindUser(&model.User{ID: uid}); user.ID == 0 {
		return echo.ErrNotFound
	}

	// model(package)の関数GetTasksを実行し、戻り値をtasks,errと定義する。
	events, err := model.GetEvents(&model.Event{UID: uid})

	// errが空でない時は StatusBadRequest(*5) を返す
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	// StasusOK と tasksを返す
	return c.JSON(http.StatusOK, events)
}

// ReqTask型は文字列のNameをパラメーターとして持つ
type ReqEvent struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func AddEventHandler(c echo.Context) error {

	event := new(model.Event)

	if err := c.Bind(event); err != nil {
		return err
	}

	uid := userIDFromToken(c)
	if user := model.FindUser(&model.User{ID: uid}); user.ID == 0 {
		return echo.ErrNotFound
	}

	event.UID = uid

	if err := model.AddEvent(event); err != nil {
		return err
	}

	// StastsOK と 追加されたtaskを返す
	return c.JSON(http.StatusOK, event)
}

func ChangeEventHandler(c echo.Context) error {
	uid := userIDFromToken(c)
	if user := model.FindUser(&model.User{ID: uid}); user.ID == 0 {
		return echo.ErrNotFound
	}

	eventID, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		return echo.ErrNotFound
	}

	event := new(model.Event)

	if err := c.Bind(event); err != nil {
		return err
	}

	event.ID = eventID

	err = model.ChangeEvent(event)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Reauest")
	}

	return c.NoContent(http.StatusOK)

}

func DeleteEventHandler(c echo.Context) error {
	uid := userIDFromToken(c)
	if user := model.FindUser(&model.User{ID: uid}); user.ID == 0 {
		return echo.ErrNotFound
	}

	eventId, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		return echo.ErrNotFound
	}

	if err := model.DeleteEvent(&model.Event{ID: eventId, UID: uid}); err != nil {
		return echo.ErrNotFound
	}

	return c.NoContent(http.StatusNoContent)

}

type AvailabilitiyRequest struct {
	MemberIDs []int `json:"member_ids"`
}

type AvailabilitiyResponse struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func GetBandMembersFreeTimeHandler(c echo.Context) error {
	// bandID := c.Param("bandID")

	var req AvailabilitiyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if len(req.MemberIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No member IDs provided",
		})
	}

	memberEvents := make(map[int][]model.Event)
	for _, memberID := range req.MemberIDs {
		var events []model.Event

		events, err := model.GetEvents(&model.Event{UID: memberID})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
		}
		memberEvents[memberID] = events
	}

	calculateCommonAvailbbilities := calculateCommonAvailbbility(memberEvents, req.MemberIDs)

	return c.JSON(http.StatusOK, calculateCommonAvailbbilities)
}

// 共通時間を計算する関数
func calculateCommonAvailbbility(memberEvents map[int][]model.Event, memberIDs []int) []AvailabilitiyResponse {

	parseEvents := make(map[int][]struct {
		StartTime time.Time
		EndTime   time.Time
	})

	timeFormat := "2006-01-02 15:04:05"

	for memberID, events := range memberEvents {
		for _, event := range events {
			startTime, _ := time.Parse(timeFormat, event.StartTime)
			endTime, _ := time.Parse(timeFormat, event.EndTime)
			parseEvents[memberID] = append(parseEvents[memberID],
				struct {
					StartTime time.Time
					EndTime   time.Time
				}{
					StartTime: startTime,
					EndTime:   endTime,
				})
		}
	}

	var commonTimeSlots []struct {
		StartTime time.Time
		EndTime   time.Time
	}
	if len(memberIDs) > 0 && len(parseEvents[memberIDs[0]]) > 0 {
		commonTimeSlots = append(commonTimeSlots, parseEvents[memberIDs[0]]...)
	} else {
		return []AvailabilitiyResponse{}
	}

	for _, memberID := range memberIDs[1:] {
		var newCommonTimeSlots []struct {
			StartTime time.Time
			EndTime   time.Time
		}
		for _, commonSlot := range commonTimeSlots {
			for _, memberSlot := range parseEvents[memberID] {
				if commonSlot.StartTime.Before(memberSlot.EndTime) &&
					commonSlot.EndTime.After(memberSlot.StartTime) {
					start := commonSlot.StartTime
					if memberSlot.StartTime.After(start) {
						start = memberSlot.StartTime
					}
					end := commonSlot.EndTime
					if memberSlot.EndTime.Before(end) {
						end = memberSlot.EndTime
					}
					newCommonTimeSlots = append(newCommonTimeSlots, struct {
						StartTime time.Time
						EndTime   time.Time
					}{
						StartTime: start,
						EndTime:   end,
					})
				}
			}
		}
		commonTimeSlots = newCommonTimeSlots
	}
	var commonAvailabilities []AvailabilitiyResponse
	for _, slot := range commonTimeSlots {
		commonAvailabilities = append(commonAvailabilities, AvailabilitiyResponse{
			StartTime: slot.StartTime.Format(timeFormat),
			EndTime:   slot.EndTime.Format(timeFormat),
		})
	}
	return commonAvailabilities
}
