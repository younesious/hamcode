package app

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (repo *Repository) CreateAttendance(attendance *Attendance) error {
	// TODO Implement
	return nil
}

func (repo *Repository) GetAttendancesByProgrammer(pid uint) ([]Attendance, error) {
	// TODO Implement
	return nil, nil
}

func (repo *Repository) GetAllAttendance() ([]Attendance, error) {
	// TODO Implement
	return nil, nil

}

func (repo *Repository) UpdateAttendance(attendance *Attendance) error {
	// TODO Implement
	return nil

}

func (repo *Repository) DeleteAttendance(pid uint) error {
	// TODO Implement
	return nil

}

func (repo *Repository) DeleteOneDayAttendance(pid uint, date time.Time) error {
	// TODO Implement
	return nil
}

func (repo *Repository) GetGirinofReport(pid uint, date time.Time) (*GirinofReport, error) {
	// TODO Implement
	return nil, nil
}

func (repo *Repository) GetMonthlyReport(pid uint, startDate, endDate time.Time) (*MonthlySummary, error) {
	// TODO Implement
	return nil, nil
}

func (repo *Repository) CalculateSalary(pid uint, startDate, endDate time.Time) (*SalaryReport, error) {
	// TODO Implement
	return nil, nil
}

func (repo *Repository) CreateProgrammer(name string) (*Programmer, error) {
	// TODO Implement
	return nil, nil
}

func (repo *Repository) GetProgrammerByID(id uint) (*Programmer, error) {
	// TODO Implement
	return nil, nil
}

func (repo *Repository) IsExistDateAndProgrammerID(pid uint, date time.Time) (*Attendance, bool) {
	// TODO Implement
	return nil, false
}
