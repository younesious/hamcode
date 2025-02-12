package app

import (
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
	ProgrammerID uint      // TODO Implement
	Date         time.Time // TODO Implement
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

type AttendanceInput struct {
	Date     time.Time  `json:"date"`
	CheckIn  *time.Time `json:"check_in,omitempty"`
	CheckOut *time.Time `json:"check_out,omitempty"`
}

type ProgrammerInput struct {
	Name string `json:"name"`
}

func (a *Attendance) UnmarshalJSON(data []byte) error {
	// TODO Implement
	return nil
}

func (a Attendance) MarshalJSON() ([]byte, error) {
	// TODO Implement
	return nil, nil
}

func (ai *AttendanceInput) UnmarshalJSON(data []byte) error {
	// TODO Implement
	return nil
}
