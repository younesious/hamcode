package main

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (repo *Repository) CreateAttendance(attendance *Attendance) error {
	log.Printf("Saving attendance to DB: %+v\n", attendance)

	return repo.DB.Debug().Create(attendance).Error
}

func (repo *Repository) GetAttendance(id uint) (*Attendance, error) {
	var attendance Attendance
	err := repo.DB.First(&attendance, id).Error
	return &attendance, err
}

func (repo *Repository) GetAllAttendance() ([]Attendance, error) {
	var attendances []Attendance
	err := repo.DB.Find(&attendances).Error
	return attendances, err
}

func (repo *Repository) UpdateAttendance(attendance *Attendance) error {
	return repo.DB.Model(&Attendance{}).Where("id = ?", attendance.ID).
		Updates(map[string]interface{}{
			"date":      attendance.Date,
			"check_in":  attendance.CheckIn,
			"check_out": attendance.CheckOut,
		}).Error
}

func (repo *Repository) DeleteAttendance(id uint) error {
	return repo.DB.Delete(&Attendance{}, id).Error
}

func (repo *Repository) GetGirinofReport(id string, date time.Time) (*GirinofReport, error) {
	var report GirinofReport

	err := repo.DB.Raw(`
		SELECT engineer,
			SUM(CASE WHEN check_in > date_trunc('day', date) + interval '09:00:00' THEN 1 ELSE 0 END) AS total_delays,
			SUM(CASE WHEN check_out < date_trunc('day', date) + interval '17:00:00' THEN 1 ELSE 0 END) AS total_early_departures
		FROM 
			attendances
		WHERE 
			engineer = ? AND date = date_trunc('day', ?)
		GROUP BY 
			engineer
	`, id, date).Scan(&report).Error

	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (repo *Repository) GetMonthlyReport(id string, startDate, endDate time.Time) (*GirinofReport, error) {
	var report GirinofReport

	err := repo.DB.Raw(`
		SELECT engineer,
			COUNT(*) AS total_days_present,
			SUM(EXTRACT(EPOCH FROM check_out - check_in) / 3600) AS total_overtime_hours,
			SUM(CASE WHEN check_in > date_trunc('day', date) + interval '09:00:00' THEN 1 ELSE 0 END) AS total_delays,
			SUM(CASE WHEN check_out < date_trunc('day', date) + interval '17:00:00' THEN 1 ELSE 0 END) AS total_early_departures
		FROM 
			attendances
		WHERE 
			engineer = ? AND date BETWEEN date_trunc('day', ?) AND date_trunc('day', ?)
		GROUP BY 
			engineer
	`, id, startDate, endDate).Scan(&report).Error

	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (repo *Repository) CalculateSalary(id string, startDate, endDate time.Time) (*SalaryReport, error) {
	var report SalaryReport

	err := repo.DB.Raw(`
		SELECT engineer,
			COUNT(*) AS total_days_present,
			SUM(EXTRACT(EPOCH FROM CASE WHEN check_out > '18:00:00' THEN check_out - '18:00:00' ELSE '00:00:00' END) / 3600) AS total_overtime_hours,
			SUM(EXTRACT(EPOCH FROM CASE WHEN check_in > '09:00:00' THEN check_in - '09:00:00' ELSE '00:00:00' END) / 3600) AS total_delay_hours
		FROM 
			attendances
		WHERE 
			engineer = ? AND date BETWEEN ? AND ?
		GROUP BY 
			engineer
	`, id, startDate, endDate).Scan(&report).Error
	if err != nil {
		return nil, err
	}

	report.TotalSalary = (float64(report.TotalDaysPresent) * DailyRate) + (report.TotalOvertimeHours * OvertimeRate) - (report.TotalDelayHours * DelayPenalty)

	return &report, nil
}
