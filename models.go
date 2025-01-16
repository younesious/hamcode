package main

import (
	"encoding/json"
	"time"
)

const (
	DailyRate    = 100.0
	OvertimeRate = 20.0
	DelayPenalty = 10.0
	CheckIn      = "09:00:00"
	CheckOut     = "17:00:00"
)

type GirinofReport struct {
	Engineer             string
	TotalDelays          int
	TotalEarlyDepartures int
}

type SalaryReport struct {
	Engineer           string  `json:"engineer"`
	TotalDaysPresent   int     `json:"total_days_present"`
	TotalOvertimeHours float64 `json:"total_overtime_hours"`
	TotalDelayHours    float64 `json:"total_delay_hours"`
	TotalSalary        float64 `json:"total_salary"`
}

type Attendance struct {
	ID        uint      `gorm:"primaryKey"`
	Engineer  string    `gorm:"size:100;not null"`
	Date      time.Time `gorm:"not null"`
	CheckIn   time.Time
	CheckOut  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *Attendance) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        uint   `json:"ID"`
		Engineer  string `json:"Engineer"`
		Date      string `json:"Date"`
		CheckIn   string `json:"CheckIn"`
		CheckOut  string `json:"CheckOut"`
		CreatedAt string `json:"CreatedAt"`
		UpdatedAt string `json:"UpdatedAt"`
	}{
		ID:        a.ID,
		Engineer:  a.Engineer,
		Date:      a.Date.Format("2006-01-02"),
		CheckIn:   a.CheckIn.UTC().Format("2006-01-02 15:04:05"),
		CheckOut:  a.CheckOut.UTC().Format("2006-01-02 15:04:05"),
		CreatedAt: a.CreatedAt.UTC().Format("2006-01-02 15:04:05"),
		UpdatedAt: a.UpdatedAt.UTC().Format("2006-01-02 15:04:05"),
	})
}

type Input struct {
	Date     *time.Time `json:"date"`
	CheckIn  *time.Time `json:"check_in"`
	CheckOut *time.Time `json:"check_out"`
}

func (i *Input) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Date     string `json:"date"`
		CheckIn  string `json:"check_in,omitempty"`
		CheckOut string `json:"check_out,omitempty"`
	}{
		Date:     i.Date.UTC().Format("2006-01-02"),
		CheckIn:  i.CheckIn.UTC().Format("2006-01-02 15:04:05"),
		CheckOut: i.CheckOut.UTC().Format("2006-01-02 15:04:05"),
	})
}
