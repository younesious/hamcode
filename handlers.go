package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var repo *Repository

func CreateAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	var attendance Attendance
	if err := json.NewDecoder(r.Body).Decode(&attendance); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("Decoded attendance: %+v\n", attendance)

	if err := repo.CreateAttendance(&attendance); err != nil {
		http.Error(w, "Failed to create attendance", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetAllAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	attendances, err := repo.GetAllAttendance()
	if err != nil {
		http.Error(w, "Failed to fetch attendance records", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(attendances)
}

func GetAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	attendance, err := repo.GetAttendance(uint(id))
	if err != nil {
		http.Error(w, "Attendance not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(attendance)
}

func UpdateAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	attendance, err := repo.GetAttendance(uint(id))
	if err != nil {
		http.Error(w, "Attendance not found", http.StatusNotFound)
		return
	}

	var input Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.Date != nil {
		attendance.Date = *input.Date
	}
	if input.CheckIn != nil {
		attendance.CheckIn = *input.CheckIn
	}
	if input.CheckOut != nil {
		attendance.CheckOut = *input.CheckOut
	}
	/*
		attendance.Date = input.Date
		attendance.CheckIn = input.CheckIn
		attendance.CheckOut = input.CheckOut
	*/
	if err := repo.UpdateAttendance(attendance); err != nil {
		http.Error(w, "Failed to update attendance", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := repo.DeleteAttendance(uint(id)); err != nil {
		http.Error(w, "Failed to delete attendance", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
