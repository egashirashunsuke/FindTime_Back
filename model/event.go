package model

import (
	_ "gorm.io/gorm"
)

type Event struct {
	UID       int    `json:"uid"`
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// 関数GetTasksは、引数はなく、戻り値は[]Event型（Event型のスライス）とerror型である
func GetEvents(event *Event) ([]Event, error) {

	// 空のタスクのスライスである、tasksを定義する
	var events []Event

	// tasksにDBのタスク全てを代入する。その操作の成否をerrと定義する(*4)
	err := db.Where(event).Find(&events).Error

	// tasksとerrを返す
	return events, err
}

func AddEvent(event *Event) error {

	if err := db.Create(event).Error; err != nil {
		return err
	}
	return nil
}

func ChangeEvent(event *Event) error {

	// DBのTaskテーブルからtaskIDと一致するidを探し、そのFinishedをtureにする(*3)
	err := db.Model(&Event{}).Where("id = ?", event.ID).Updates(map[string]interface{}{"StartTime": event.StartTime, "EndTime": event.EndTime}).Error
	return err
}

func DeleteEvent(event *Event) error {
	// DBのTaskテーブルからtaskIDと一致するidを探し、そのタスクを削除する
	err := db.Where("id = ?", event.ID).Delete(&Event{}).Error
	return err
}
