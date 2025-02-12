package app

import (
	"log"
	"net/http"
)

func StartServer() {
	db := GetConnection()
	Repo = &Repository{DB: db}

	mux := http.NewServeMux()
	RegisterRoutes(mux)

	log.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /attendance", CreateAttendanceHandler)
	mux.HandleFunc("GET /attendance/{programmer_id}", GetAttendanceHandler)
	mux.HandleFunc("PUT /attendance/{programmer_id}", UpdateAttendanceHandler)
	mux.HandleFunc("DELETE /attendance/{programmer_id}", DeleteAttendanceHandler)
	mux.HandleFunc("DELETE /attendance/{programmer_id}/{date}", DeleteOneDayAttendanceHandler)
	mux.HandleFunc("GET /attendance", GetAllAttendanceHandler)

	mux.HandleFunc("GET /girinof/{programmer_id}/{date}", GetGirinofReportHandler)
	mux.HandleFunc("GET /report/monthly/{programmer_id}/{start_date}/{end_date}", GetMonthlyReportHandler)
	mux.HandleFunc("GET /salary/{programmer_id}", GetSalaryHandler)

	mux.HandleFunc("POST /programmer", CreateProgrammerHandler)
}
