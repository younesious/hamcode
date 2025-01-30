package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var testRepo *Repository

// Test DB Connection
func TestGetConnection(t *testing.T) {
	db := GetConnection()
	assert.NotNil(t, db)

	// Test singleton pattern
	db2 := GetConnection()
	assert.Equal(t, db, db2)
}

func TestMain(m *testing.M) {
	var err error
	testDB = GetConnection()

	testDB.AutoMigrate(&Programmer{}, &Attendance{})

	// Initialize repository
	testRepo = &Repository{DB: testDB}
	repo = testRepo

	// Run tests
	code := m.Run()

	// Cleanup
	sqlDB, err := testDB.DB()
	if err == nil {
		sqlDB.Close()
	}

	os.Exit(code)
}

// Helper function to create a test programmer
func createTestProgrammer(t *testing.T) *Programmer {
	prog, err := testRepo.CreateProgrammer("Test Programmer")
	assert.NoError(t, err)
	return prog
}

// Test CreateProgrammerHandler
func TestCreateProgrammerHandler(t *testing.T) {
	tests := []struct {
		name           string
		input          ProgrammerInput
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid programmer creation",
			input:          ProgrammerInput{Name: "John Doe"},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name:           "Invalid input - empty name",
			input:          ProgrammerInput{Name: ""},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/programmer", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			CreateProgrammerHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if !tt.expectError {
				var resp Programmer
				err := json.NewDecoder(w.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Name, resp.Name)
			}
		})
	}
}

// Test CreateAttendanceHandler
func TestCreateAttendanceHandler(t *testing.T) {
	prog := createTestProgrammer(t)

	tests := []struct {
		name           string
		attendance     Attendance
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Valid attendance creation",
			attendance: Attendance{
				ProgrammerID: prog.ID,
				Date:         time.Now(),
				CheckIn:      time.Now(),
				CheckOut:     time.Now().Add(8 * time.Hour),
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "Invalid - non-existent programmer",
			attendance: Attendance{
				ProgrammerID: 9999,
				Date:         time.Now(),
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.attendance)
			req := httptest.NewRequest("POST", "/attendance", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			CreateAttendanceHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if !tt.expectError {
				var resp Attendance
				err := json.NewDecoder(w.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.attendance.ProgrammerID, resp.ProgrammerID)
			}
		})
	}
}

// Test GetAttendanceHandler
func TestGetAttendanceHandler(t *testing.T) {
	prog := createTestProgrammer(t)
	attendance := Attendance{
		ProgrammerID: prog.ID,
		Date:         time.Now(),
		CheckIn:      time.Now(),
		CheckOut:     time.Now().Add(8 * time.Hour),
	}
	err := testRepo.CreateAttendance(&attendance)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		programmerID   string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid attendance fetch",
			programmerID:   fmt.Sprint(prog.ID),
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid programmer ID",
			programmerID:   "invalid",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Non-existent programmer",
			programmerID:   "9999",
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/attendance/"+tt.programmerID, nil)
			w := httptest.NewRecorder()

			GetAttendanceHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if !tt.expectError {
				var resp []Attendance
				err := json.NewDecoder(w.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.NotEmpty(t, resp)
			}
		})
	}
}

// Test UpdateAttendanceHandler
func TestUpdateAttendanceHandler(t *testing.T) {
	prog := createTestProgrammer(t)
	attendance := Attendance{
		ProgrammerID: prog.ID,
		Date:         time.Now(),
		CheckIn:      time.Now(),
		CheckOut:     time.Now().Add(8 * time.Hour),
	}
	err := testRepo.CreateAttendance(&attendance)
	assert.NoError(t, err)

	newCheckIn := time.Now().Add(-1 * time.Hour)
	tests := []struct {
		name           string
		programmerID   string
		input          AttendanceInput
		expectedStatus int
	}{
		{
			name:         "Valid update",
			programmerID: fmt.Sprint(prog.ID),
			input: AttendanceInput{
				Date:    attendance.Date,
				CheckIn: &newCheckIn,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid programmer ID",
			programmerID:   "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("PUT", "/attendance/"+tt.programmerID, bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			UpdateAttendanceHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// Test DeleteAttendanceHandler
func TestDeleteAttendanceHandler(t *testing.T) {
	prog := createTestProgrammer(t)
	attendance := Attendance{
		ProgrammerID: prog.ID,
		Date:         time.Now(),
		CheckIn:      time.Now(),
		CheckOut:     time.Now().Add(8 * time.Hour),
	}
	err := testRepo.CreateAttendance(&attendance)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		programmerID   string
		expectedStatus int
	}{
		{
			name:           "Valid delete",
			programmerID:   fmt.Sprint(prog.ID),
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "Invalid programmer ID",
			programmerID:   "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/attendance/"+tt.programmerID, nil)
			w := httptest.NewRecorder()

			DeleteAttendanceHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// Test DeleteOneDayAttendanceHandler
func TestDeleteOneDayAttendanceHandler(t *testing.T) {
	prog := createTestProgrammer(t)
	today := time.Now()

	// Create test attendance record
	attendance := Attendance{
		ProgrammerID: prog.ID,
		Date:         today,
		CheckIn:      today.Add(9 * time.Hour),
		CheckOut:     today.Add(17 * time.Hour),
	}
	err := testRepo.CreateAttendance(&attendance)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		programmerID   string
		date           string
		expectedStatus int
	}{
		{
			name:           "Valid deletion",
			programmerID:   fmt.Sprint(prog.ID),
			date:           today.Format("2006-01-02"),
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "Invalid programmer ID",
			programmerID:   "invalid",
			date:           today.Format("2006-01-02"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid date format",
			programmerID:   fmt.Sprint(prog.ID),
			date:           "invalid-date",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Non-existent record",
			programmerID:   fmt.Sprint(prog.ID),
			date:           time.Now().AddDate(0, 0, -10).Format("2006-01-02"),
			expectedStatus: http.StatusInternalServerError, // or whatever status you return for non-existent records
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/attendance/%s/%s", tt.programmerID, tt.date), nil)
			w := httptest.NewRecorder()

			DeleteOneDayAttendanceHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusAccepted {
				// Verify the record was actually deleted
				found, exists := testRepo.IsExistDateAndProgrammerID(prog.ID, today)
				assert.False(t, exists)
				assert.Nil(t, found)
			}
		})
	}
}

// Test GetAllAttendanceHandler
func TestGetAllAttendanceHandler(t *testing.T) {
	// Create multiple programmers and attendance records
	prog1 := createTestProgrammer(t)
	prog2 := createTestProgrammer(t)
	today := time.Now()

	// Create attendance records for both programmers
	attendances := []Attendance{
		{
			ProgrammerID: prog1.ID,
			Date:         today,
			CheckIn:      today.Add(9 * time.Hour),
			CheckOut:     today.Add(17 * time.Hour),
		},
		{
			ProgrammerID: prog2.ID,
			Date:         today,
			CheckIn:      today.Add(8 * time.Hour),
			CheckOut:     today.Add(16 * time.Hour),
		},
		{
			ProgrammerID: prog1.ID,
			Date:         today.AddDate(0, 0, -1),
			CheckIn:      today.AddDate(0, 0, -1).Add(9 * time.Hour),
			CheckOut:     today.AddDate(0, 0, -1).Add(17 * time.Hour),
		},
	}

	for _, att := range attendances {
		err := testRepo.CreateAttendance(&att)
		assert.NoError(t, err)
	}

	tests := []struct {
		name           string
		expectedCount  int
		expectedStatus int
		clearDB        bool // Flag to test empty database scenario
	}{
		{
			name:           "Get all attendance records",
			expectedCount:  3,
			expectedStatus: http.StatusOK,
			clearDB:        false,
		},
		{
			name:           "Empty database",
			expectedCount:  0,
			expectedStatus: http.StatusOK,
			clearDB:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clearDB {
				// Clear the attendance records
				testDB.Where("1=1").Delete(&Attendance{})
			}

			req := httptest.NewRequest("GET", "/attendance", nil)
			w := httptest.NewRecorder()

			GetAllAttendanceHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response []Attendance
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(response))

			if !tt.clearDB {
				// Verify the response contains the expected data
				foundProg1 := false
				foundProg2 := false
				for _, att := range response {
					if att.ProgrammerID == prog1.ID {
						foundProg1 = true
					}
					if att.ProgrammerID == prog2.ID {
						foundProg2 = true
					}
				}
				assert.True(t, foundProg1, "Should contain attendance records for programmer 1")
				assert.True(t, foundProg2, "Should contain attendance records for programmer 2")
			}
		})
	}
}

// Test GetGirinofReportHandler
func TestGetGirinofReportHandler(t *testing.T) {
	prog := createTestProgrammer(t)
	attendance := Attendance{
		ProgrammerID: prog.ID,
		Date:         time.Now(),
		CheckIn:      time.Now().Add(30 * time.Minute), // 30 minutes late
		CheckOut:     time.Now().Add(7 * time.Hour),    // 1 hour early
	}
	err := testRepo.CreateAttendance(&attendance)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		programmerID   string
		date           string
		expectedStatus int
	}{
		{
			name:           "Valid report",
			programmerID:   fmt.Sprint(prog.ID),
			date:           attendance.Date.Format("2006-01-02"),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid date format",
			programmerID:   fmt.Sprint(prog.ID),
			date:           "invalid-date",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/girinof/%s/%s", tt.programmerID, tt.date), nil)
			w := httptest.NewRecorder()

			GetGirinofReportHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var report GirinofReport
				err := json.NewDecoder(w.Body).Decode(&report)
				assert.NoError(t, err)
				assert.NotZero(t, report.TotalDelayMinutes)
				assert.NotZero(t, report.TotalEarlyDepartures)
			}
		})
	}
}

// Test GetMonthlyReportHandler
func TestGetMonthlyReportHandler(t *testing.T) {
	prog := createTestProgrammer(t)
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	// Create some attendance records
	for i := 0; i < 5; i++ {
		attendance := Attendance{
			ProgrammerID: prog.ID,
			Date:         startDate.AddDate(0, 0, i),
			CheckIn:      startDate.AddDate(0, 0, i).Add(9 * time.Hour),
			CheckOut:     startDate.AddDate(0, 0, i).Add(17 * time.Hour),
		}
		err := testRepo.CreateAttendance(&attendance)
		assert.NoError(t, err)
	}

	tests := []struct {
		name           string
		programmerID   string
		startDate      string
		endDate        string
		expectedStatus int
	}{
		{
			name:           "Valid report",
			programmerID:   fmt.Sprint(prog.ID),
			startDate:      startDate.Format("2006-01-02"),
			endDate:        endDate.Format("2006-01-02"),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid date format",
			programmerID:   fmt.Sprint(prog.ID),
			startDate:      "invalid-date",
			endDate:        endDate.Format("2006-01-02"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/report/monthly/%s/%s/%s",
				tt.programmerID, tt.startDate, tt.endDate), nil)
			w := httptest.NewRecorder()

			GetMonthlyReportHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var report MonthlySummary
				err := json.NewDecoder(w.Body).Decode(&report)
				assert.NoError(t, err)
				assert.Equal(t, 5, report.TotalDaysPresent)
			}
		})
	}
}

// Test GetSalaryHandler
func TestGetSalaryHandler(t *testing.T) {
	prog := createTestProgrammer(t)
	startDate := time.Now().AddDate(0, 0, -30)

	// Create some attendance records
	for i := 0; i < 5; i++ {
		attendance := Attendance{
			ProgrammerID: prog.ID,
			Date:         startDate.AddDate(0, 0, i),
			CheckIn:      startDate.AddDate(0, 0, i).Add(9 * time.Hour),
			CheckOut:     startDate.AddDate(0, 0, i).Add(18 * time.Hour), // 1 hour overtime
		}
		err := testRepo.CreateAttendance(&attendance)
		assert.NoError(t, err)
	}

	tests := []struct {
		name           string
		programmerID   string
		expectedStatus int
	}{
		{
			name:           "Valid salary calculation",
			programmerID:   fmt.Sprint(prog.ID),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing programmer ID",
			programmerID:   "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/salary/"+tt.programmerID, nil)
			w := httptest.NewRecorder()

			GetSalaryHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var report SalaryReport
				err := json.NewDecoder(w.Body).Decode(&report)
				assert.NoError(t, err)
				assert.Equal(t, 5, report.TotalDaysPresent)
				assert.NotZero(t, report.TotalSalary)
				assert.NotZero(t, report.TotalOvertimeMinutes)
			}
		})
	}
}

// Repository Tests
func TestRepository(t *testing.T) {
	// Test CreateProgrammer
	t.Run("CreateProgrammer", func(t *testing.T) {
		prog, err := testRepo.CreateProgrammer("Test Developer")
		assert.NoError(t, err)
		assert.NotNil(t, prog)
		assert.NotZero(t, prog.ID)
		assert.Equal(t, "Test Developer", prog.Name)
	})

	// Test GetProgrammerByID
	t.Run("GetProgrammerByID", func(t *testing.T) {
		prog, err := testRepo.CreateProgrammer("Test Dev")
		assert.NoError(t, err)

		// Test successful retrieval
		found, err := testRepo.GetProgrammerByID(prog.ID)
		assert.NoError(t, err)
		assert.Equal(t, prog.ID, found.ID)
		assert.Equal(t, prog.Name, found.Name)

		// Test non-existent programmer
		_, err = testRepo.GetProgrammerByID(99999)
		assert.Error(t, err)
	})

	// Test CreateAttendance
	t.Run("CreateAttendance", func(t *testing.T) {
		prog, _ := testRepo.CreateProgrammer("Test Dev")
		attendance := &Attendance{
			ProgrammerID: prog.ID,
			Date:         time.Now(),
			CheckIn:      time.Now(),
			CheckOut:     time.Now().Add(8 * time.Hour),
		}

		err := testRepo.CreateAttendance(attendance)
		assert.NoError(t, err)
		assert.NotZero(t, attendance.ID)

		// Test duplicate entry
		err = testRepo.CreateAttendance(attendance)
		assert.Error(t, err) // Should fail due to unique constraint
	})

	// Test GetAttendancesByProgrammer
	t.Run("GetAttendancesByProgrammer", func(t *testing.T) {
		prog, _ := testRepo.CreateProgrammer("Test Dev")
		attendance := &Attendance{
			ProgrammerID: prog.ID,
			Date:         time.Now(),
			CheckIn:      time.Now(),
			CheckOut:     time.Now().Add(8 * time.Hour),
		}
		testRepo.CreateAttendance(attendance)

		records, err := testRepo.GetAttendancesByProgrammer(prog.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, records)
		assert.Equal(t, prog.ID, records[0].ProgrammerID)
	})

	// Test GetAllAttendance
	t.Run("GetAllAttendance", func(t *testing.T) {
		prog, _ := testRepo.CreateProgrammer("Test Dev")
		attendance := &Attendance{
			ProgrammerID: prog.ID,
			Date:         time.Now(),
			CheckIn:      time.Now(),
			CheckOut:     time.Now().Add(8 * time.Hour),
		}
		testRepo.CreateAttendance(attendance)

		records, err := testRepo.GetAllAttendance()
		assert.NoError(t, err)
		assert.NotEmpty(t, records)
	})

	// Test UpdateAttendance
	t.Run("UpdateAttendance", func(t *testing.T) {
		prog, _ := testRepo.CreateProgrammer("Test Dev")
		attendance := &Attendance{
			ProgrammerID: prog.ID,
			Date:         time.Now(),
			CheckIn:      time.Now(),
			CheckOut:     time.Now().Add(8 * time.Hour),
		}
		testRepo.CreateAttendance(attendance)

		newCheckOut := time.Now().Add(9 * time.Hour)
		attendance.CheckOut = newCheckOut
		err := testRepo.UpdateAttendance(attendance)
		assert.NoError(t, err)

		// Verify update
		updated, _ := testRepo.GetAttendancesByProgrammer(prog.ID)
		assert.Equal(t, newCheckOut.Unix(), updated[0].CheckOut.Unix())
	})

	// Test DeleteAttendance
	t.Run("DeleteAttendance", func(t *testing.T) {
		prog, _ := testRepo.CreateProgrammer("Test Dev")
		attendance := &Attendance{
			ProgrammerID: prog.ID,
			Date:         time.Now(),
			CheckIn:      time.Now(),
			CheckOut:     time.Now().Add(8 * time.Hour),
		}
		testRepo.CreateAttendance(attendance)

		err := testRepo.DeleteAttendance(prog.ID)
		assert.NoError(t, err)

		// Verify deletion
		records, _ := testRepo.GetAttendancesByProgrammer(prog.ID)
		assert.Empty(t, records)
	})

	// Test DeleteOneDayAttendance
	t.Run("DeleteOneDayAttendance", func(t *testing.T) {
		prog, _ := testRepo.CreateProgrammer("Test Dev")
		today := time.Now()
		attendance := &Attendance{
			ProgrammerID: prog.ID,
			Date:         today,
			CheckIn:      today,
			CheckOut:     today.Add(8 * time.Hour),
		}
		testRepo.CreateAttendance(attendance)

		err := testRepo.DeleteOneDayAttendance(prog.ID, today)
		assert.NoError(t, err)

		// Verify deletion
		records, _ := testRepo.GetAttendancesByProgrammer(prog.ID)
		assert.Empty(t, records)
	})

	// Test IsExistDateAndProgrammerID
	t.Run("IsExistDateAndProgrammerID", func(t *testing.T) {
		prog, _ := testRepo.CreateProgrammer("Test Dev")
		today := time.Now()
		attendance := &Attendance{
			ProgrammerID: prog.ID,
			Date:         today,
			CheckIn:      today,
			CheckOut:     today.Add(8 * time.Hour),
		}
		testRepo.CreateAttendance(attendance)

		// Test existing record
		found, exists := testRepo.IsExistDateAndProgrammerID(prog.ID, today)
		assert.True(t, exists)
		assert.NotNil(t, found)

		// Test non-existing record
		found, exists = testRepo.IsExistDateAndProgrammerID(prog.ID, today.AddDate(0, 0, 1))
		assert.False(t, exists)
		assert.Nil(t, found)
	})

	// Test GetGirinofReport
	t.Run("GetGirinofReport", func(t *testing.T) {
		prog, _ := testRepo.CreateProgrammer("Test Dev")
		today := time.Now()
		attendance := &Attendance{
			ProgrammerID: prog.ID,
			Date:         today,
			CheckIn:      today.Add(10 * time.Hour), // 1 hour late
			CheckOut:     today.Add(16 * time.Hour), // 1 hour early
		}
		testRepo.CreateAttendance(attendance)

		report, err := testRepo.GetGirinofReport(fmt.Sprint(prog.ID), today)
		assert.NoError(t, err)
		assert.NotNil(t, report)
		assert.NotZero(t, report.TotalDelayMinutes)
		assert.NotZero(t, report.TotalEarlyDepartures)
	})

	// Test GetMonthlyReport
	t.Run("GetMonthlyReport", func(t *testing.T) {
		prog, _ := testRepo.CreateProgrammer("Test Dev")
		startDate := time.Now().AddDate(0, 0, -5)
		endDate := time.Now()

		// Create some attendance records
		for i := 0; i < 5; i++ {
			attendance := &Attendance{
				ProgrammerID: prog.ID,
				Date:         startDate.AddDate(0, 0, i),
				CheckIn:      startDate.AddDate(0, 0, i).Add(9 * time.Hour),
				CheckOut:     startDate.AddDate(0, 0, i).Add(18 * time.Hour), // 1 hour overtime
			}
			testRepo.CreateAttendance(attendance)
		}

		report, err := testRepo.GetMonthlyReport(fmt.Sprint(prog.ID), startDate, endDate)
		assert.NoError(t, err)
		assert.NotNil(t, report)
		assert.Equal(t, 5, report.TotalDaysPresent)
		assert.NotZero(t, report.TotalOvertimeMinutes)
	})

	// Test CalculateSalary
	t.Run("CalculateSalary", func(t *testing.T) {
		prog, _ := testRepo.CreateProgrammer("Test Dev")
		startDate := time.Now().AddDate(0, 0, -5)
		endDate := time.Now()

		// Create some attendance records
		for i := 0; i < 5; i++ {
			attendance := &Attendance{
				ProgrammerID: prog.ID,
				Date:         startDate.AddDate(0, 0, i),
				CheckIn:      startDate.AddDate(0, 0, i).Add(9 * time.Hour),
				CheckOut:     startDate.AddDate(0, 0, i).Add(18 * time.Hour), // 1 hour overtime
			}
			testRepo.CreateAttendance(attendance)
		}

		report, err := testRepo.CalculateSalary(fmt.Sprint(prog.ID), startDate, endDate)
		assert.NoError(t, err)
		assert.NotNil(t, report)
		assert.Equal(t, 5, report.TotalDaysPresent)
		assert.NotZero(t, report.TotalSalary)
		assert.NotZero(t, report.TotalOvertimeMinutes)

		// Verify salary calculation
		expectedBase := float64(report.TotalDaysPresent) * DailyRate
		expectedOvertime := float64(report.TotalOvertimeMinutes) * OvertimeRate
		expectedTotal := expectedBase + expectedOvertime
		assert.Equal(t, expectedTotal, report.TotalSalary)
	})
}
