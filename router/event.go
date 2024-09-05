package router

import (
	"FindTime-Server/model"
	"log"
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

	events, err := model.GetEvents(&model.Event{UID: uid})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, events)
}

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
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
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
	MemberIDs   []int `json:"member_ids"`
	MinDuration int   `json:"min_duration"`
}

type AvailabilitiyResponse struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func GetBandMembersFreeTimeHandler(c echo.Context) error {
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
		events, err := model.GetEvents(&model.Event{UID: memberID})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
		}
		memberEvents[memberID] = events
	}

	calculateCommonAvailabilities := calculateCommonAvailbbility(memberEvents, req.MemberIDs, req.MinDuration)

	return c.JSON(http.StatusOK, calculateCommonAvailabilities)
}

func calculateCommonAvailbbility(memberEvents map[int][]model.Event, memberIDs []int, minDuration int) []AvailabilitiyResponse {
	timeFormat := "2006-01-02 15:04:05"

	parseEvents := make(map[int][]struct {
		StartTime time.Time
		EndTime   time.Time
	})

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
		Members   []int
	}

	if len(memberIDs) > 0 && len(parseEvents[memberIDs[0]]) > 0 {
		for _, slot := range parseEvents[memberIDs[0]] {
			commonTimeSlots = append(commonTimeSlots, struct {
				StartTime time.Time
				EndTime   time.Time
				Members   []int
			}{
				StartTime: slot.StartTime,
				EndTime:   slot.EndTime,
				Members:   []int{memberIDs[0]},
			})
		}
	} else {
		return []AvailabilitiyResponse{}
	}

	for _, memberID := range memberIDs[1:] {
		var newCommonTimeSlots []struct {
			StartTime time.Time
			EndTime   time.Time
			Members   []int
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
						Members   []int
					}{
						StartTime: start,
						EndTime:   end,
						Members:   append(commonSlot.Members, memberID),
					})
				}
			}
		}
		commonTimeSlots = newCommonTimeSlots
	}

	var filteredCommonAvailabilities []AvailabilitiyResponse
	minDurationDuration := float64(minDuration) // minDuration を float64 に変換して分単位に

	for _, slot := range commonTimeSlots {
		duration := slot.EndTime.Sub(slot.StartTime).Minutes()
		log.Printf("Slot Start: %v, Slot End: %v, Duration: %v\n", slot.StartTime, slot.EndTime, duration)
		if duration >= minDurationDuration {
			log.Printf("Slot meets the criteria: Start %v, End %v\n", slot.StartTime, slot.EndTime)
			filteredCommonAvailabilities = append(filteredCommonAvailabilities, AvailabilitiyResponse{
				StartTime: slot.StartTime.Format(timeFormat),
				EndTime:   slot.EndTime.Format(timeFormat),
			})
		}
	}

	return filteredCommonAvailabilities
}
