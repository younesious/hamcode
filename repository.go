package main

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (repo *Repository) CreateAttendance(attendance *Attendance) error {
	return repo.DB.Debug().Create(attendance).Error
}

func (repo *Repository) GetAttendancesByProgrammer(pid uint) ([]Attendance, error) {
	var attendances []Attendance

	err := repo.DB.Preload("Programmer").
		Where("programmer_id = ?", pid).
		Order("date DESC").
		Limit(12).
		Find(&attendances).Error

	return attendances, err
}

func (repo *Repository) GetAllAttendance() ([]Attendance, error) {
	var attendances []Attendance
	err := repo.DB.Find(&attendances).Error
	return attendances, err
}

func (repo *Repository) UpdateAttendance(attendance *Attendance) error {
	return repo.DB.Model(&Attendance{}).Where("id = ?", attendance.ID).
		Updates(map[string]interface{}{
			"date":      attendance.Date.UTC(),
			"check_in":  attendance.CheckIn.UTC(),
			"check_out": attendance.CheckOut.UTC(),
		}).Error
}

func (repo *Repository) DeleteAttendance(pid uint) error {
	return repo.DB.Where("programmer_id = ?", pid).Delete(&Attendance{}).Error
}

func (repo *Repository) DeleteOneDayAttendance(pid uint, date time.Time) error {
	return repo.DB.Where("programmer_id = ? AND date = ?", pid, date).Delete(&Attendance{}).Error
}

func (repo *Repository) GetGirinofReport(pid uint, date time.Time) (*GirinofReport, error) {
	var report GirinofReport

	err := repo.DB.Raw(`
		SELECT 
			programmer_id,
			SUM(EXTRACT(EPOCH FROM (CASE WHEN check_in > date_trunc('day', date::timestamp) + interval '09:00:00' 
				THEN check_in - (date_trunc('day', date::timestamp) + interval '09:00:00') 
				ELSE interval '00:00:00' END)) / 60) AS total_delay_minutes,
			SUM(EXTRACT(EPOCH FROM (CASE WHEN check_out < date_trunc('day', date::timestamp) + interval '17:00:00' 
				THEN (date_trunc('day', date::timestamp) + interval '17:00:00') - check_out 
				ELSE interval '00:00:00' END)) / 60) AS total_early_departure_minutes
		FROM 
			attendances
		WHERE 
			programmer_id = ? AND date = date_trunc('day', ?::timestamp)
		GROUP BY 
			programmer_id
	`, pid, date).Scan(&report).Error

	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (repo *Repository) GetMonthlyReport(pid uint, startDate, endDate time.Time) (*MonthlySummary, error) {
	var report MonthlySummary

	err := repo.DB.Raw(`
		SELECT 
			programmer_id,
			COUNT(*) AS total_days_present,
			SUM(EXTRACT(EPOCH FROM (CASE WHEN check_out > date_trunc('day', date::timestamp) + interval '17:00:00' 
				THEN check_out - (date_trunc('day', date::timestamp) + interval '17:00:00') 
				ELSE interval '00:00:00' END)) / 60) AS total_overtime_minutes,
			SUM(EXTRACT(EPOCH FROM (CASE WHEN check_in > date_trunc('day', date::timestamp) + interval '09:00:00' 
				THEN check_in - (date_trunc('day', date::timestamp) + interval '09:00:00') 
				ELSE interval '00:00:00' END)) / 60) AS total_delay_minutes,
			SUM(EXTRACT(EPOCH FROM (CASE WHEN check_out < date_trunc('day', date::timestamp) + interval '17:00:00' 
				THEN (date_trunc('day', date::timestamp) + interval '17:00:00') - check_out 
				ELSE interval '00:00:00' END)) / 60) AS total_early_departure_minutes
		FROM 
			attendances
		WHERE 
			programmer_id = ? AND date BETWEEN date_trunc('day', ?::timestamp) AND date_trunc('day', ?::timestamp)
		GROUP BY 
			programmer_id
	`, pid, startDate, endDate).Scan(&report).Error

	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (repo *Repository) CalculateSalary(pid uint, startDate, endDate time.Time) (*SalaryReport, error) {
	var report SalaryReport

	err := repo.DB.Raw(`
		SELECT 
			p.name,
			COUNT(*) AS total_days_present,
			SUM(EXTRACT(EPOCH FROM (CASE WHEN a.check_out > date_trunc('day', a.date::timestamp) + interval '17:00:00' 
				THEN a.check_out - (date_trunc('day', a.date::timestamp) + interval '17:00:00') 
				ELSE interval '00:00:00' END)) / 60) AS total_overtime_minutes,
			SUM(EXTRACT(EPOCH FROM (CASE WHEN a.check_in > date_trunc('day', a.date::timestamp) + interval '09:00:00' 
				THEN a.check_in - (date_trunc('day', a.date::timestamp) + interval '09:00:00') 
				ELSE interval '00:00:00' END)) / 60) AS total_delay_minutes,
			SUM(EXTRACT(EPOCH FROM (CASE WHEN a.check_out < date_trunc('day', a.date::timestamp) + interval '17:00:00' 
				THEN (date_trunc('day', a.date::timestamp) + interval '17:00:00') - a.check_out 
				ELSE interval '00:00:00' END)) / 60) AS total_early_departure_minutes
		FROM 
			attendances a
		JOIN public.programmers p on p.id = a.programmer_id
		WHERE 
			a.programmer_id = ? AND a.date BETWEEN date_trunc('day', ?::timestamp) AND date_trunc('day', ?::timestamp)
		GROUP BY 
			p.name
	`, pid, startDate, endDate).Scan(&report).Error

	if err != nil {
		return nil, err
	}

	report.TotalSalary = (float64(report.TotalDaysPresent) * DailyRate) + (float64(report.TotalOvertimeMinutes) * OvertimeRate) -
		(float64(report.TotalDelayMinutes) * DelayPenalty)

	return &report, nil
}

func (repo *Repository) CreateProgrammer(name string) (*Programmer, error) {
	programmer := Programmer{
		Name: name,
	}
	if err := repo.DB.Create(&programmer).Error; err != nil {
		return nil, err
	}
	return &programmer, nil
}

func (repo *Repository) GetProgrammerByID(id uint) (*Programmer, error) {
	var programmer Programmer
	err := repo.DB.First(&programmer, id).Error
	if err != nil {
		return nil, err
	}
	return &programmer, nil
}

func (repo *Repository) IsExistDateAndProgrammerID(pid uint, date time.Time) (*Attendance, bool) {
	var attendance Attendance

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.UTC().Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	if err := repo.DB.Where("programmer_id = ? AND date >= ? AND date < ?", pid, startOfDay, endOfDay).First(&attendance).Error; err != nil {
		return nil, false
	}

	return &attendance, true
}
