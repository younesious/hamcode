package main

import (
	"log"

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
