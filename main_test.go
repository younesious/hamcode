package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type MyTestRepository struct {
	DB *gorm.DB
}

func (r *MyTestRepository) myCreateAttendance(attendance *Attendance) error {
	return r.DB.Create(attendance).Error
}

func (r *MyTestRepository) myCreateProgrammer(names ...string) (*Programmer, error) {
	name := "Younes Mahmoudi"
	if len(names) > 0 {
		name = names[0]
	}
	programmer := Programmer{
		Name: name,
	}
	if err := repo.DB.Create(&programmer).Error; err != nil {
		return nil, err
	}
	return &programmer, nil
}

func (r *MyTestRepository) myIsExistDateAndProgrammerID(id uint, date time.Time) (*Attendance, bool) {
	var attendance Attendance
	if err := r.DB.Where("programmer_id = ? AND date = ?", id, date).First(&attendance).Error; err != nil {
		return nil, false
	}
	return &attendance, true
}

func seedAttendanceData(t *testing.T, repo *MyTestRepository, programmerID uint) {
	// Helper function to parse time strings
	parseTime := func(timeStr string) time.Time {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			t.Fatalf("Failed to parse time: %v", err)
		}
		return parsedTime
	}

	now := time.Now().UTC()

	attendances := []Attendance{
		// Late check-in
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -30).Truncate(24 * time.Hour), // 30 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -30).Format("2006-01-02") + " 09:30:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -30).Format("2006-01-02") + " 17:00:00"),
		},
		// Overtime
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -29).Truncate(24 * time.Hour), // 29 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -29).Format("2006-01-02") + " 09:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -29).Format("2006-01-02") + " 18:30:00"),
		},
		// Early departure
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -28).Truncate(24 * time.Hour), // 28 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -28).Format("2006-01-02") + " 08:50:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -28).Format("2006-01-02") + " 16:45:00"),
		},
		// Slight delay
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -27).Truncate(24 * time.Hour), // 27 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -27).Format("2006-01-02") + " 09:10:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -27).Format("2006-01-02") + " 17:00:00"),
		},
		// On time
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -26).Truncate(24 * time.Hour), // 26 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -26).Format("2006-01-02") + " 09:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -26).Format("2006-01-02") + " 17:00:00"),
		},
		// Delay and early departure
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -25).Truncate(24 * time.Hour), // 25 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -25).Format("2006-01-02") + " 09:15:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -25).Format("2006-01-02") + " 16:50:00"),
		},
		// Major delay
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -24).Truncate(24 * time.Hour), // 24 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -24).Format("2006-01-02") + " 10:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -24).Format("2006-01-02") + " 17:30:00"),
		},
		// Overtime
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -23).Truncate(24 * time.Hour), // 23 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -23).Format("2006-01-02") + " 09:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -23).Format("2006-01-02") + " 19:00:00"),
		},
		// Early departure
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -22).Truncate(24 * time.Hour), // 22 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -22).Format("2006-01-02") + " 08:55:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -22).Format("2006-01-02") + " 16:30:00"),
		},
		// Normal
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -21).Truncate(24 * time.Hour), // 21 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -21).Format("2006-01-02") + " 09:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -21).Format("2006-01-02") + " 17:00:00"),
		},
		// Delay and minor overtime
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -20).Truncate(24 * time.Hour), // 20 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -20).Format("2006-01-02") + " 09:20:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -20).Format("2006-01-02") + " 17:10:00"),
		},
		// Slight delay and overtime
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -19).Truncate(24 * time.Hour), // 19 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -19).Format("2006-01-02") + " 09:05:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -19).Format("2006-01-02") + " 18:00:00"),
		},
		// Normal
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -18).Truncate(24 * time.Hour), // 18 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -18).Format("2006-01-02") + " 09:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -18).Format("2006-01-02") + " 17:00:00"),
		},
		// Late and early departure
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -17).Truncate(24 * time.Hour), // 17 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -17).Format("2006-01-02") + " 09:45:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -17).Format("2006-01-02") + " 15:30:00"),
		},
		// Delay and minor overtime
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -16).Truncate(24 * time.Hour), // 16 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -16).Format("2006-01-02") + " 09:10:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -16).Format("2006-01-02") + " 17:15:00"),
		},
		// Normal
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -15).Truncate(24 * time.Hour), // 15 days ago
			CheckIn:      parseTime(now.AddDate(0, 0, -15).Format("2006-01-02") + " 09:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -15).Format("2006-01-02") + " 17:00:00"),
		},
		// Day 14: On time
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -14).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -14).Format("2006-01-02") + " 09:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -14).Format("2006-01-02") + " 17:00:00"),
		},
		// Day 13: 30 minutes late
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -13).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -13).Format("2006-01-02") + " 09:30:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -13).Format("2006-01-02") + " 17:00:00"),
		},
		// Day 12: 1-hour overtime
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -12).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -12).Format("2006-01-02") + " 09:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -12).Format("2006-01-02") + " 18:00:00"),
		},
		// Day 11: 20 minutes late, 10 minutes early departure
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -11).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -11).Format("2006-01-02") + " 09:20:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -11).Format("2006-01-02") + " 16:50:00"),
		},
		// Day 10: Early departure (45 minutes early)
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -10).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -10).Format("2006-01-02") + " 09:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -10).Format("2006-01-02") + " 16:15:00"),
		},
		// Day 9: Slight delay (10 minutes late)
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -9).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -9).Format("2006-01-02") + " 09:10:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -9).Format("2006-01-02") + " 17:00:00"),
		},
		// Day 8: Major delay (1 hour late)
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -8).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -8).Format("2006-01-02") + " 10:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -8).Format("2006-01-02") + " 17:30:00"),
		},
		// Day 7: Overtime (1.5 hours overtime)
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -7).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -7).Format("2006-01-02") + " 09:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -7).Format("2006-01-02") + " 18:30:00"),
		},
		// Day 6: Early departure (30 minutes early)
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -6).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -6).Format("2006-01-02") + " 08:55:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -6).Format("2006-01-02") + " 16:30:00"),
		},
		// Day 5: Normal (on time)
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -5).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -5).Format("2006-01-02") + " 09:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -5).Format("2006-01-02") + " 17:00:00"),
		},
		// Day 4: Delay and minor overtime
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -4).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -4).Format("2006-01-02") + " 09:10:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -4).Format("2006-01-02") + " 17:15:00"),
		},
		// Day 3: Late and early departure (45 minutes late, 1.5 hours early)
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -3).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -3).Format("2006-01-02") + " 09:45:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -3).Format("2006-01-02") + " 15:30:00"),
		},
		// Day 2: Delay and overtime (10 minutes late, 30 minutes overtime)
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -2).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -2).Format("2006-01-02") + " 09:10:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -2).Format("2006-01-02") + " 17:30:00"),
		},
		// Day 1: Normal (on time)
		{
			ProgrammerID: programmerID,
			Date:         now.AddDate(0, 0, -1).Truncate(24 * time.Hour),
			CheckIn:      parseTime(now.AddDate(0, 0, -1).Format("2006-01-02") + " 09:00:00"),
			CheckOut:     parseTime(now.AddDate(0, 0, -1).Format("2006-01-02") + " 17:00:00"),
		},
	}

	// Insert the attendance records into the database
	for _, att := range attendances {
		err := repo.myCreateAttendance(&att)
		assert.NoError(t, err)
	}
}

func clearTables(t *testing.T) {
	t.Helper()

	err := testDB.Where("1=1").Delete(&Attendance{}).Error
	assert.NoError(t, err)

	err = testDB.Where("1=1").Delete(&Programmer{}).Error
	assert.NoError(t, err)
}

func TestGetConnection(t *testing.T) {
	db := GetConnection()
	assert.NotNil(t, db)
}

func TestMain(m *testing.M) {
	var err error
	testDB = GetConnection()

	testRepo = &Repository{DB: testDB}
	repo = testRepo

	code := m.Run()

	sqlDB, err := testDB.DB()
	if err == nil {
		sqlDB.Close()
	}

	os.Exit(code)
}

// Test CreateProgrammerHandler
func TestCreateProgrammerHandler(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	tests := []struct {
		name           string
		input          ProgrammerInput
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid programmer creation",
			input:          ProgrammerInput{Name: "Younes Mahmoudi"},
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

				var dbProg Programmer
				err = testDB.First(&dbProg, resp.ID).Error
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Name, dbProg.Name)

				var count int64
				testDB.Model(&Programmer{}).Count(&count)
				assert.Equal(t, int64(1), count)
			}
		})
	}
}

// Test CreateAttendanceHandler
func TestCreateAttendanceHandler(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer()
	assert.NoError(t, err)

	tests := []struct {
		name           string
		attendance     Attendance
		expectedStatus int
		expectError    bool
		errorMessage   string
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
				CheckIn:      time.Now(),
				CheckOut:     time.Now().Add(8 * time.Hour),
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
			errorMessage:   "Programmer not found",
		},
		{
			name: "Invalid - CheckIn after CheckOut",
			attendance: Attendance{
				ProgrammerID: prog.ID,
				Date:         time.Now(),
				CheckIn:      time.Now().Add(2 * time.Hour),
				CheckOut:     time.Now(),
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			errorMessage:   "CheckOut cannot be before CheckIn",
		},
		{
			name: "Invalid - duplicate date",
			attendance: Attendance{
				ProgrammerID: prog.ID,
				Date:         time.Now().Truncate(24 * time.Hour),
				CheckIn:      time.Now().Truncate(24 * time.Hour).Add(9 * time.Hour),
				CheckOut:     time.Now().Truncate(24 * time.Hour).Add(17 * time.Hour),
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
			errorMessage:   "duplicate key",
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

				var programmer Programmer
				testDB.First(&programmer, resp.ProgrammerID)
				assert.Equal(t, "Younes Mahmoudi", programmer.Name)

				found, exists := repo.myIsExistDateAndProgrammerID(resp.ProgrammerID, resp.Date)
				assert.True(t, exists)
				assert.Equal(t, resp.CheckIn.UTC().Format(time.RFC3339), found.CheckIn.UTC().Format(time.RFC3339))
			}

			if tt.errorMessage != "" {
				body, _ := io.ReadAll(w.Body)
				assert.Contains(t, string(body), tt.errorMessage)
			}
		})
	}
}

// Test GetAttendanceHandler
func TestGetAttendanceHandler(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer()
	assert.NoError(t, err)

	attendance := Attendance{
		ProgrammerID: prog.ID,
		Date:         time.Now(),
		CheckIn:      time.Now(),
		CheckOut:     time.Now().Add(8 * time.Hour),
	}

	err = repo.myCreateAttendance(&attendance)
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

				found, exists := repo.myIsExistDateAndProgrammerID(prog.ID, attendance.Date)
				assert.True(t, exists)
				assert.Equal(t, found.ProgrammerID, resp[0].ProgrammerID)
				assert.Equal(t, found.Date.UTC().Format("2006-01-02"), resp[0].Date.UTC().Format("2006-01-02"))
				assert.Equal(t, found.CheckIn.UTC().Format("2006-01-02 15:04"), resp[0].CheckIn.UTC().Format("2006-01-02 15:04"))
				assert.Equal(t, found.CheckOut.UTC().Format("2006-01-02 15:04"), resp[0].CheckOut.UTC().Format("2006-01-02 15:04"))
			}
		})
	}
}

// TestUpdateAttendanceHandler
func TestUpdateAttendanceHandler(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer()
	assert.NoError(t, err)

	attendance := Attendance{
		ProgrammerID: prog.ID,
		Date:         time.Now(),
		CheckIn:      time.Now(),
		CheckOut:     time.Now().Add(8 * time.Hour),
	}
	err = repo.myCreateAttendance(&attendance)
	assert.NoError(t, err)

	newCheckIn := time.Now().Add(-1 * time.Hour)
	newCheckOut := newCheckIn.Add(-1 * time.Hour) // CheckOut before CheckIn

	tests := []struct {
		name           string
		programmerID   string
		input          AttendanceInput
		expectedStatus int
		errorMessage   string
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
		{
			name:           "Missing date",
			programmerID:   fmt.Sprint(prog.ID),
			input:          AttendanceInput{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "CheckIn > CheckOut",
			programmerID: fmt.Sprint(prog.ID),
			input: AttendanceInput{
				Date:     attendance.Date,
				CheckIn:  &newCheckIn,
				CheckOut: &newCheckOut,
			},
			expectedStatus: http.StatusBadRequest,
			errorMessage:   "CheckIn cannot be after CheckOut",
		},
		{
			name:         "Non-existent attendance record",
			programmerID: fmt.Sprint(prog.ID),
			input: AttendanceInput{
				Date:    time.Now().AddDate(0, 0, -10),
				CheckIn: &newCheckIn,
			},
			expectedStatus: http.StatusBadRequest,
			errorMessage:   "Programmer with the given date not recorded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := struct {
				Date     string  `json:"date"`
				CheckIn  *string `json:"check_in,omitempty"`
				CheckOut *string `json:"check_out,omitempty"`
			}{
				Date:     tt.input.Date.Format("2006-01-02"),
				CheckIn:  nil,
				CheckOut: nil,
			}

			if tt.input.CheckIn != nil {
				checkInStr := tt.input.CheckIn.UTC().Format("2006-01-02 15:04:05")
				input.CheckIn = &checkInStr
			}

			if tt.input.CheckOut != nil {
				checkOutStr := tt.input.CheckOut.UTC().Format("2006-01-02 15:04:05")
				input.CheckOut = &checkOutStr
			}

			body, err := json.Marshal(input)
			assert.NoError(t, err)

			req := httptest.NewRequest("PUT", "/attendance/"+tt.programmerID, bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			UpdateAttendanceHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.errorMessage != "" {
				body, _ := io.ReadAll(w.Body)
				assert.Contains(t, string(body), tt.errorMessage)
			}

			if tt.expectedStatus == http.StatusOK {
				updated, exists := repo.myIsExistDateAndProgrammerID(prog.ID, tt.input.Date)
				assert.True(t, exists)
				assert.NotNil(t, updated)

				if tt.input.CheckIn != nil {
					expectedCheckIn := tt.input.CheckIn.UTC().Format("2006-01-02 15:04")
					actualCheckIn := updated.CheckIn.UTC().Format("2006-01-02 15:04")
					assert.Equal(t, expectedCheckIn, actualCheckIn)
				} else {
					assert.Equal(t, attendance.CheckIn.UTC().Format("2006-01-02 15:04"), updated.CheckIn.UTC().Format("2006-01-02 15:04"))
				}

				if tt.input.CheckOut != nil {
					expectedCheckOut := tt.input.CheckOut.UTC().Format("2006-01-02 15:04")
					actualCheckOut := updated.CheckOut.UTC().Format("2006-01-02 15:04")
					assert.Equal(t, expectedCheckOut, actualCheckOut)
				} else {
					assert.Equal(t, attendance.CheckOut.UTC().Format("2006-01-02 15:04"), updated.CheckOut.UTC().Format("2006-01-02 15:04"))
				}
			}
		})
	}
}

// Test DeleteAttendanceHandler
func TestDeleteAttendanceHandler(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer()
	assert.NoError(t, err)

	attendance := Attendance{
		ProgrammerID: prog.ID,
		Date:         time.Now(),
		CheckIn:      time.Now(),
		CheckOut:     time.Now().Add(8 * time.Hour),
	}
	err = repo.myCreateAttendance(&attendance)
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
		{
			name:           "Non-Existent record",
			programmerID:   "99999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/attendance/"+tt.programmerID, nil)
			w := httptest.NewRecorder()

			DeleteAttendanceHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify deletion
			if tt.expectedStatus == http.StatusAccepted {
				_, exists := repo.myIsExistDateAndProgrammerID(prog.ID, attendance.Date)
				assert.False(t, exists)
			}
		})
	}
}

// Test DeleteOneDayAttendanceHandler
func TestDeleteOneDayAttendanceHandler(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer()
	assert.NoError(t, err)

	today := time.Now()

	attendance := Attendance{
		ProgrammerID: prog.ID,
		Date:         today,
		CheckIn:      today.Add(9 * time.Hour),
		CheckOut:     today.Add(17 * time.Hour),
	}
	err = repo.myCreateAttendance(&attendance)
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
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/attendance/%s/%s", tt.programmerID, tt.date), nil)
			w := httptest.NewRecorder()

			DeleteOneDayAttendanceHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusAccepted {
				parsedDate, err := time.Parse("2006-01-02", tt.date)
				if err != nil {
					t.Fatalf("Failed to parse date: %v", err)
				}

				found, exists := repo.myIsExistDateAndProgrammerID(prog.ID, parsedDate)
				assert.False(t, exists)
				assert.Nil(t, found)
			}
		})
	}
}

// Test GetAllAttendanceHandler
func TestGetAllAttendanceHandler(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog1, err := repo.myCreateProgrammer()
	assert.NoError(t, err)
	prog2, err := repo.myCreateProgrammer("Yahya Sinwar") // TODO I think here we can change the signuture of createTestProgrammer func to have default value and accept optianol name for this scenario
	assert.NoError(t, err)

	today := time.Now()

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
		err := repo.myCreateAttendance(&att)
		assert.NoError(t, err)
	}

	tests := []struct {
		name           string
		expectedCount  int
		expectedStatus int
	}{
		{
			name:           "Get all attendance records",
			expectedCount:  3,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty database",
			expectedCount:  0,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Empty database" {
				clearTables(t)
			}

			req := httptest.NewRequest("GET", "/attendance", nil)
			w := httptest.NewRecorder()

			GetAllAttendanceHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response []Attendance
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(response))

			if tt.name != "Empty database" {
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

/*
	INSERT INTO attendances (programmer_id, date, check_in, check_out) VALUES
		(1, '2024-12-16', '2024-12-16 09:30:00', '2024-12-16 17:00:00'), -- Late check-in
		(1, '2024-12-17', '2024-12-17 09:00:00', '2024-12-17 18:30:00'), -- Overtime
		(1, '2024-12-18', '2024-12-18 08:50:00', '2024-12-18 16:45:00'), -- Early departure
		(1, '2024-12-19', '2024-12-19 09:10:00', '2024-12-19 17:00:00'), -- Slight delay
		(1, '2024-12-20', '2024-12-20 09:00:00', '2024-12-20 17:00:00'), -- On time
		(1, '2024-12-21', '2024-12-21 09:15:00', '2024-12-21 16:50:00'), -- Delay and early departure
		(1, '2024-12-22', '2024-12-22 10:00:00', '2024-12-22 17:30:00'), -- Major delay
		(1, '2024-12-23', '2024-12-23 09:00:00', '2024-12-23 19:00:00'), -- Overtime
		(1, '2024-12-24', '2024-12-24 08:55:00', '2024-12-24 16:30:00'), -- Early departure
		(1, '2024-12-25', '2024-12-25 09:00:00', '2024-12-25 17:00:00'), -- Normal
		(1, '2024-12-26', '2024-12-26 09:20:00', '2024-12-26 17:10:00'), -- Delay and minor overtime
		(1, '2024-12-27', '2024-12-27 09:05:00', '2024-12-27 18:00:00'), -- Slight delay and overtime
		(1, '2024-12-28', '2024-12-28 09:00:00', '2024-12-28 17:00:00'), -- Normal
		(1, '2024-12-29', '2024-12-29 09:45:00', '2024-12-29 15:30:00'), -- Late and early departure
		(1, '2024-12-30', '2024-12-30 09:10:00', '2024-12-30 17:15:00'), -- Delay and minor overtime
		(1, '2024-12-31', '2024-12-31 09:00:00', '2024-12-31 17:00:00'), -- Normal
		(1, '2025-01-01', '2025-01-01 08:55:00', '2025-01-01 16:55:00'), -- Early departure
		(1, '2025-01-02', '2025-01-02 09:30:00', '2025-01-02 17:10:00'), -- Delay and minor overtime
		(1, '2025-01-03', '2025-01-03 09:00:00', '2025-01-03 17:00:00'), -- Normal
		(1, '2025-01-04', '2025-01-04 09:20:00', '2025-01-04 16:50:00'), -- Delay and early departure
		(1, '2025-01-05', '2025-01-05 10:00:00', '2025-01-05 15:00:00'), -- Major delay and early departure
		(1, '2025-01-06', '2025-01-06 09:00:00', '2025-01-06 18:00:00'), -- Overtime
		(1, '2025-01-07', '2025-01-07 08:55:00', '2025-01-07 16:40:00'), -- Early departure
		(1, '2025-01-08', '2025-01-08 09:00:00', '2025-01-08 17:00:00'), -- Normal
		(1, '2025-01-09', '2025-01-09 09:40:00', '2025-01-09 17:30:00'), -- Delay and overtime
		(1, '2025-01-10', '2025-01-10 09:00:00', '2025-01-10 17:00:00'), -- Normal
		(1, '2025-01-11', '2025-01-11 09:05:00', '2025-01-11 17:15:00'), -- Slight delay and minor overtime
		(1, '2025-01-12', '2025-01-12 09:25:00', '2025-01-12 16:45:00'), -- Delay and early departure
		(1, '2025-01-13', '2025-01-13 09:00:00', '2025-01-13 17:00:00'), -- Normal
		(1, '2025-01-14', '2025-01-14 09:15:00', '2025-01-14 18:00:00'), -- Delay and overtime
		(1, '2025-01-15', '2025-01-15 09:00:00', '2025-01-15 17:00:00'); -- Normal
*/

// Test GetGirinofReportHandler
func TestGetGirinofReportHandler(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer()
	assert.NoError(t, err)

	seedAttendanceData(t, repo, prog.ID)

	now := time.Now().UTC()
	date30DaysAgo := now.AddDate(0, 0, -30).Truncate(24 * time.Hour).Format("2006-01-02")
	date28DaysAgo := now.AddDate(0, 0, -28).Truncate(24 * time.Hour).Format("2006-01-02")
	date25DaysAgo := now.AddDate(0, 0, -25).Truncate(24 * time.Hour).Format("2006-01-02")
	currentDate := now.Truncate(24 * time.Hour).Format("2006-01-02")
	tests := []struct {
		name                       string
		programmerID               string
		date                       string
		expectedStatus             int
		expectedDelayMinutes       float64
		expectedEarlyDepartureMins float64
	}{
		{
			name:                       "Valid report with 30 minutes delay",
			programmerID:               fmt.Sprint(prog.ID),
			date:                       date30DaysAgo,
			expectedStatus:             http.StatusOK,
			expectedDelayMinutes:       30,
			expectedEarlyDepartureMins: 0,
		},
		{
			name:                       "Valid report with 15 minutes early departure",
			programmerID:               fmt.Sprint(prog.ID),
			date:                       date28DaysAgo,
			expectedStatus:             http.StatusOK,
			expectedDelayMinutes:       0,
			expectedEarlyDepartureMins: 15,
		},
		{
			name:                       "Valid report with 15 minutes delay and 10 minutes early",
			programmerID:               fmt.Sprint(prog.ID),
			date:                       date25DaysAgo,
			expectedStatus:             http.StatusOK,
			expectedDelayMinutes:       15,
			expectedEarlyDepartureMins: 10,
		},
		{
			name:                       "Valid date with no attendance",
			programmerID:               fmt.Sprint(prog.ID),
			date:                       currentDate,
			expectedStatus:             http.StatusBadRequest,
			expectedDelayMinutes:       0,
			expectedEarlyDepartureMins: 0,
		},
		{
			name:                       "Invalid date format",
			programmerID:               fmt.Sprint(prog.ID),
			date:                       "invalid-date",
			expectedStatus:             http.StatusBadRequest,
			expectedDelayMinutes:       0,
			expectedEarlyDepartureMins: 0,
		},
		{
			name:           "Missing programmer id",
			programmerID:   "",
			date:           currentDate,
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
				assert.NoError(t, err, "Expected valid JSON response")
				assert.Equal(t, tt.programmerID, report.ProgrammerID, "unexpected programmerID")
				assert.Equal(t, tt.expectedDelayMinutes, report.TotalDelayMinutes, "unexpected total delay minutes")
				assert.Equal(t, tt.expectedEarlyDepartureMins, report.TotalEarlyDepartures, "unexpected total early departure minutes")
			}
		})
	}
}

/*
// Test GetMonthlyReportHandler
func TestGetMonthlyReportHandler(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer()
	assert.NoError(t, err)

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
	} // TODO add test case for endDate before start date and check err message

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
				assert.Equal(t, 5, report.TotalDaysPresent) // TODO also check total_early_departure_minutes and Equal for their expected value
				// TODO here is the result of this query in DB: (programmer_id, total_days_present, total_overtime_minutes, total_delay_minutes, total_early_departure_minutes) = (1,28,410,360,300)
			}
		})
	}
}

/*
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
		} // also here help from helper function I explained in girinof
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
			} // TODO also check for Equality of fileds
			// TODO here is the result of this query in DB: (name, total_days_present, total_overtime_minutes, total_delay_minutes, total_early_departure_minutes, total_salary) =
			// TODO (Younes Mahmoudi,28,410,360,300, // TODO do it placeholder I fill it)

		})
	}
}

// Repository Tests // TODO and I thikn with todo comments I add to test the value in database after each endpoint test
// I think TestRepository is useless and not make sense, are you agree with me? and in integration test with endpoint I test every thing and test cverage is 100%, right?
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

// Model tests
func TestAttendance_Marshaling(t *testing.T) {
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	checkIn := date.Add(9 * time.Hour)
	checkOut := date.Add(17 * time.Hour)

	a := Attendance{
		ProgrammerID: 1,
		Date:         date,
		CheckIn:      checkIn,
		CheckOut:     checkOut,
	}

	// Test Marshal
	data, err := json.Marshal(a)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"date":"2024-01-01"`)              // I think here Equality is more make sense vs contain
	assert.Contains(t, string(data), `"check_in":"2024-01-01 09:00:00"`) // I think here Equality is more make sense vs contain

	// Test Unmarshal
	var a2 Attendance
	err = json.Unmarshal(data, &a2)
	assert.NoError(t, err)
	assert.True(t, date.Equal(a2.Date))
	assert.True(t, checkIn.Equal(a2.CheckIn))
}

func TestAttendanceInput_UnmarshalJSON(t *testing.T) {
	payload := `{
		"date": "2024-01-01",
		"check_in": "2024-01-01 09:00:00",
		"check_out": null
	}`

	var input AttendanceInput
	err := json.Unmarshal([]byte(payload), &input)
	assert.NoError(t, err)
	assert.NotNil(t, input.CheckIn)
	assert.Nil(t, input.CheckOut)
	// TODO check for exact format we expected and value of that
}
*/
