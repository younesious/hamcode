package main

import (
	"log"
	"net/http"
)

func StartServer() {
	db := GetConnection()
	repo = &Repository{DB: db}

	mux := http.NewServeMux()
	RegisterRoutes(mux)

	log.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /attendance", CreateAttendanceHandler)
	mux.HandleFunc("GET /attendance/{id}", GetAttendanceHandler)
	mux.HandleFunc("PUT /attendance/{id}", UpdateAttendanceHandler)
	mux.HandleFunc("DELETE /attendance/{id}", DeleteAttendanceHandler)
	mux.HandleFunc("GET /attendance", GetAllAttendanceHandler)
}
