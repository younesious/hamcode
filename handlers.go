package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var repo *Repository

func CreateAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	var attendance Attendance
	if err := json.NewDecoder(r.Body).Decode(&attendance); err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}

	if attendance.Date.IsZero() {
		http.Error(w, "Missing required field: date", http.StatusBadRequest)
		return
	}

	if !attendance.CheckOut.IsZero() && !attendance.CheckIn.IsZero() && attendance.CheckOut.Before(attendance.CheckIn) {
		http.Error(w, "CheckOut cannot be before CheckIn", http.StatusBadRequest)
		return
	}

	_, err := repo.GetProgrammerByID(attendance.ProgrammerID)
	if err != nil {
		http.Error(w, "Programmer not found! Please hire the programmer first.", http.StatusNotFound)
		return
	}

	if err := repo.CreateAttendance(&attendance); err != nil {
		errorMessage := "Failed to create attendance: " + err.Error()
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(attendance)
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
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if parts[2] == "" {
		http.Error(w, "Invalid URL path: programmer_id is required", http.StatusBadRequest)
		return
	}
	pidStr := parts[2]
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		http.Error(w, "Invalid Programmer ID", http.StatusBadRequest)
		return
	}

	attendances, err := repo.GetAttendancesByProgrammer(uint(pid))
	if err != nil {
		http.Error(w, "Error retrieving attendance records", http.StatusInternalServerError)
		return
	}

	if len(attendances) == 0 {
		http.Error(w, "No attendance records found for the given Programmer ID", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(attendances); err != nil {
		http.Error(w, "Failed to encode attendance records", http.StatusInternalServerError)
		return
	}
}

func UpdateAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if parts[2] == "" {
		http.Error(w, "Programmer ID is required", http.StatusBadRequest)
		return
	}
	pidStr := parts[2]
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		http.Error(w, "Invalid Programmer ID", http.StatusBadRequest)
		return
	}

	_, err = repo.GetProgrammerByID(uint(pid))
	if err != nil {
		http.Error(w, "Programmer not found", http.StatusNotFound)
		return
	}

	var input AttendanceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	existingAttendance, exist := repo.IsExistDateAndProgrammerID(uint(pid), input.Date)
	if !exist {
		http.Error(w, "Programmer with the given date not recorded. Please use CreateAttendanceHandler to create it!", http.StatusBadRequest)
		return
	}

	if input.CheckIn != nil {
		existingAttendance.CheckIn = *input.CheckIn
	}
	if input.CheckOut != nil {
		existingAttendance.CheckOut = *input.CheckOut
	}
	if input.CheckIn != nil && input.CheckOut != nil && input.CheckOut.Before(*input.CheckIn) {
		http.Error(w, "CheckIn cannot be after CheckOut", http.StatusBadRequest)
		return
	}

	if err := repo.UpdateAttendance(existingAttendance); err != nil {
		http.Error(w, "Failed to update attendance", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if parts[2] == "" {
		http.Error(w, "Programmer ID is required", http.StatusBadRequest)
		return
	}
	pidStr := parts[2]
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		http.Error(w, "Invalid Programmer ID", http.StatusBadRequest)
		return
	}

	if _, err := repo.GetProgrammerByID(uint(pid)); err != nil {
		http.Error(w, "Programmer not found", http.StatusNotFound)
		return
	}

	if err := repo.DeleteAttendance(uint(pid)); err != nil {
		http.Error(w, "Failed to delete attendance", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func DeleteOneDayAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if parts[2] == "" {
		http.Error(w, "Programmer ID is required", http.StatusBadRequest)
		return
	}
	pidStr := parts[2]
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		http.Error(w, "Invalid Programmer ID", http.StatusBadRequest)
		return
	}

	if parts[3] == "" {
		http.Error(w, "Date is required", http.StatusBadRequest)
		return
	}
	dateStr := parts[3]

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid Date", http.StatusBadRequest)
		return
	}

	if _, exists := repo.IsExistDateAndProgrammerID(uint(pid), date); !exists {
		http.Error(w, "No attendance record found for the given date", http.StatusNotFound)
		return
	}

	if err := repo.DeleteOneDayAttendance(uint(pid), date); err != nil {
		http.Error(w, "Failed to delete attendance", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func GetGirinofReportHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if parts[2] == "" {
		http.Error(w, "Programmer ID is required", http.StatusBadRequest)
		return
	}
	pidStr := parts[2]

	if parts[3] == "" {
		http.Error(w, "Date is required", http.StatusBadRequest)
		return
	}
	dateStr := parts[3]

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid Date", http.StatusBadRequest)
		return
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		http.Error(w, "Invalid Programmer ID", http.StatusBadRequest)
		return
	}
	if _, err := repo.GetProgrammerByID(uint(pid)); err != nil {
		http.Error(w, "Programmer not found", http.StatusNotFound)
		return
	}
	_, exist := repo.IsExistDateAndProgrammerID(uint(pid), date)
	if !exist {
		http.Error(w, "Programmer with the given date not recorded. Please use CreateAttendanceHandler to create it!", http.StatusBadRequest)
		return
	}

	report, err := repo.GetGirinofReport(pidStr, date)
	if err != nil {
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(report)
}

func GetMonthlyReportHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if parts[3] == "" {
		http.Error(w, "Programmer ID is required", http.StatusBadRequest)
		return
	}
	pidStr := parts[3]
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		http.Error(w, "Invalid programmer ID", http.StatusBadRequest)
		return
	}

	if parts[4] == "" {
		http.Error(w, "Start date is required", http.StatusBadRequest)
		return
	}
	startDateStr := parts[4]
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Invalid start date format", http.StatusBadRequest)
		return
	}
	if parts[5] == "" {
		http.Error(w, "End date is required", http.StatusBadRequest)
		return
	}
	endDateStr := parts[5]
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Invalid end date format", http.StatusBadRequest)
		return
	}

	if _, err := repo.GetProgrammerByID(uint(pid)); err != nil {
		http.Error(w, "Programmer not found", http.StatusNotFound)
		return
	}

	if endDate.Before(startDate) {
		http.Error(w, "End date cannot be before start date", http.StatusBadRequest)
		return
	}

	report, err := repo.GetMonthlyReport(pidStr, startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to generate monthly report", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(report)
}

func GetSalaryHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("programmer_id")
	if id == "" {
		http.Error(w, "Programmer ID is required", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC().Truncate(24 * time.Hour)
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	salary, err := repo.CalculateSalary(id, thirtyDaysAgo, now)
	if err != nil {
		http.Error(w, "Failed to calculate salary", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(salary)
}

func CreateProgrammerHandler(w http.ResponseWriter, r *http.Request) {
	var input ProgrammerInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	programmer, err := repo.CreateProgrammer(input.Name)
	if err != nil {
		http.Error(w, "Failed to create programmer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(programmer)
}
