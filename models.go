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

func (a *Attendance) UnmarshalJSON(data []byte) error {
	type Alias Attendance

	aux := &struct {
		Date     string `json:"date"`
		CheckIn  string `json:"check_in"`
		CheckOut string `json:"check_out"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	if a.Date, err = time.Parse("2006-01-02", aux.Date); err != nil {
		return err
	}
	if a.CheckIn, err = time.Parse("2006-01-02 15:04:05", aux.CheckIn); err != nil {
		return err
	}
	if a.CheckOut, err = time.Parse("2006-01-02 15:04:05", aux.CheckOut); err != nil {
		return err
	}

	return nil
}

type Input struct {
	Date     *time.Time `json:"date"`
	CheckIn  *time.Time `json:"check_in"`
	CheckOut *time.Time `json:"check_out"`
}

func (i *Input) UnmarshalJSON(data []byte) error {
	aux := struct {
		Date     *string `json:"date"`
		CheckIn  *string `json:"check_in"`
		CheckOut *string `json:"check_out"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Date != nil {
		parsedDate, err := time.Parse("2006-01-02", *aux.Date)
		if err != nil {
			return err
		}
		i.Date = &parsedDate
	}
	if aux.CheckIn != nil {
		parsedCheckIn, err := time.Parse("2006-01-02 15:04:05", *aux.CheckIn)
		if err != nil {
			return err
		}
		i.CheckIn = &parsedCheckIn
	}
	if aux.CheckOut != nil {
		parsedCheckOut, err := time.Parse("2006-01-02 15:04:05", *aux.CheckOut)
		if err != nil {
			return err
		}
		i.CheckOut = &parsedCheckOut
	}

	return nil
}
