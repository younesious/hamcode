package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	DailyRate    = 100.0
	OvertimeRate = 20.0
	DelayPenalty = 10.0
)

type GirinofReport struct {
	ProgrammerID         string  `json:"programmer_id" gorm:"column:programmer_id"`
	TotalDelayMinutes    float64 `json:"total_delays" gorm:"column:total_delay_minutes"`
	TotalEarlyDepartures float64 `json:"total_early_departures" gorm:"column:total_early_departure_minutes"`
}

type MonthlySummary struct {
	ProgrammerID         string  `json:"programmer_id" gorm:"column:programmer_id"`
	TotalDaysPresent     int     `json:"total_days_present" gorm:"column:total_days_present"`
	TotalOvertimeMinutes float64 `json:"total_overtime_minutes" gorm:"column:total_overtime_minutes"`
	TotalDelayMinutes    float64 `json:"total_delay_minutes" gorm:"column:total_delay_minutes"`
	TotalEarlyDepartures float64 `json:"total_early_departure_minutes" gorm:"column:total_early_departure_minutes"`
}

type SalaryReport struct {
	Name                 string  `json:"name" gorm:"column:name"`
	TotalDaysPresent     int     `json:"total_days_present" gorm:"column:total_days_present"`
	TotalOvertimeMinutes float64 `json:"total_overtime_minutes" gorm:"column:total_overtime_minutes"`
	TotalDelayMinutes    float64 `json:"total_delay_minutes" gorm:"column:total_delay_minutes"`
	TotalEarlyDepartures float64 `json:"total_early_departure_minutes" gorm:"column:total_early_departure_minutes"`
	TotalSalary          float64 `json:"total_salary" gorm:"total_salary"`
}

type Attendance struct {
	ID           uint      `gorm:"primaryKey"`
	ProgrammerID uint      `gorm:"uniqueIndex:idx_programmer_date,unique,priority:1;not null"`
	Date         time.Time `gorm:"uniqueIndex:idx_programmer_date,unique,priority:2;not null"`
	CheckIn      time.Time
	CheckOut     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Programmer Programmer `gorm:"foreignKey:ProgrammerID"`
}

type Programmer struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (a *Attendance) UnmarshalJSON(data []byte) error {
	type Alias Attendance

	aux := &struct {
		Date         string `json:"date"`
		CheckIn      string `json:"check_in"`
		CheckOut     string `json:"check_out"`
		ProgrammerID uint   `json:"programmer_id"`
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

	a.ProgrammerID = aux.ProgrammerID

	return nil
}

func (a Attendance) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID           uint   `json:"id"`
		ProgrammerID uint   `json:"programmer_id"`
		Date         string `json:"date"`
		CheckIn      string `json:"check_in"`
		CheckOut     string `json:"check_out"`
	}{
		ID:           a.ID,
		ProgrammerID: a.ProgrammerID,
		Date:         a.Date.UTC().Format("2006-01-02"),
		CheckIn:      a.CheckIn.UTC().Format("2006-01-02 15:04:05"),
		CheckOut:     a.CheckOut.UTC().Format("2006-01-02 15:04:05"),
	})
}

type AttendanceInput struct {
	Date     time.Time  `json:"date"`
	CheckIn  *time.Time `json:"check_in"`
	CheckOut *time.Time `json:"check_out"`
}

func (i *AttendanceInput) UnmarshalJSON(data []byte) error {
	aux := struct {
		Date     string  `json:"date"`
		CheckIn  *string `json:"check_in"`
		CheckOut *string `json:"check_out"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Date == "" {
		return errors.New("date field is nessasery")
	}

	parsedDate, err := time.Parse("2006-01-02", aux.Date)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %v", err)
	}
	i.Date = parsedDate

	if aux.CheckIn != nil {
		parsedCheckIn, err := time.Parse("2006-01-02 15:04:05", *aux.CheckIn)
		if err != nil {
			return fmt.Errorf("invalid check_in format, expected YYYY-MM-DD HH:mm:ss: %v", err)
		}
		i.CheckIn = &parsedCheckIn
	}
	if aux.CheckOut != nil {
		parsedCheckOut, err := time.Parse("2006-01-02 15:04:05", *aux.CheckOut)
		if err != nil {
			return fmt.Errorf("invalid check_out format, expected YYYY-MM-DD HH:mm:ss: %v", err)
		}
		i.CheckOut = &parsedCheckOut
	}

	return nil
}

type ProgrammerInput struct {
	Name string `json:"name"`
}
