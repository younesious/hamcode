package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"younesious/hamcode/app"
)

var testDB *gorm.DB

type MyTestRepository struct {
	DB *gorm.DB
}

func (r *MyTestRepository) myCreateAttendance(attendance *app.Attendance) error {
	return r.DB.Create(attendance).Error
}

func (r *MyTestRepository) myCreateProgrammer(names ...string) (*app.Programmer, error) {
	name := "Younes Mahmoudi"
	if len(names) > 0 {
		name = names[0]
	}
	programmer := app.Programmer{
		Name: name,
	}
	if err := r.DB.Create(&programmer).Error; err != nil {
		return nil, err
	}
	return &programmer, nil
}

func (r *MyTestRepository) myIsExistDateAndProgrammerID(pid uint, date time.Time) (*app.Attendance, bool) {
	var attendance app.Attendance

	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	if err := r.DB.Where("programmer_id = ? AND date >= ? AND date < ?",
		pid, startOfDay, endOfDay).First(&attendance).Error; err != nil {
		return nil, false
	}

	return &attendance, true

}

func (r *MyTestRepository) myGetAttendancesByProgrammer(pid uint) ([]app.Attendance, error) {
	var attendances []app.Attendance

	err := r.DB.Preload("Programmer").
		Where("programmer_id = ?", pid).
		Order("date DESC").
		Limit(12).
		Find(&attendances).Error

	return attendances, err
}

func (repo *MyTestRepository) myGetAllAttendance() ([]app.Attendance, error) {
	var attendances []app.Attendance
	err := repo.DB.Find(&attendances).Error
	return attendances, err
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

	attendances := []app.Attendance{
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

	for _, att := range attendances {
		err := repo.myCreateAttendance(&att)
		assert.NoError(t, err)
	}
}

func clearTables(t *testing.T) {
	t.Helper()

	tables := []string{"programmers", "attendances"}
	for _, table := range tables {
		err := testDB.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error
		assert.NoError(t, err)
	}
}

func TestGetConnection(t *testing.T) {
	db := app.GetConnection()
	assert.NotNil(t, db)
}

func TestMain(m *testing.M) {
	var err error
	testDB = app.GetConnection()

	testDB.Migrator().DropTable(&app.Attendance{})
	testDB.AutoMigrate(&app.Attendance{})

	app.Repo = &app.Repository{DB: testDB}

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
		input          app.ProgrammerInput
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid programmer creation",
			input:          app.ProgrammerInput{Name: "Younes Mahmoudi"},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name:           "Invalid input - empty name",
			input:          app.ProgrammerInput{Name: ""},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/programmer", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			app.CreateProgrammerHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if !tt.expectError {
				var resp app.Programmer
				err := json.NewDecoder(w.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Name, resp.Name)

				var dbProg app.Programmer
				err = testDB.First(&dbProg, resp.ID).Error
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Name, dbProg.Name)

				var count int64
				testDB.Model(&app.Programmer{}).Count(&count)
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
		attendance     app.Attendance
		expectedStatus int
		expectError    bool
		errorMessage   string
	}{
		{
			name: "Valid attendance creation",
			attendance: app.Attendance{
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
			attendance: app.Attendance{
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
			attendance: app.Attendance{
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
			attendance: app.Attendance{
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

			app.CreateAttendanceHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if !tt.expectError {
				var resp app.Attendance
				err := json.NewDecoder(w.Body).Decode(&resp)
				assert.NoError(t, err)

				assert.Equal(t, tt.attendance.ProgrammerID, resp.ProgrammerID)

				var programmer app.Programmer
				testDB.First(&programmer, resp.ProgrammerID)
				assert.Equal(t, "Younes Mahmoudi", programmer.Name)

				found, exists := repo.myIsExistDateAndProgrammerID(resp.ProgrammerID, resp.Date)
				if found != nil {
					assert.True(t, exists)
					assert.Equal(t, resp.CheckIn.UTC().Format(time.RFC3339), found.CheckIn.UTC().Format(time.RFC3339))
				}
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

	attendance := app.Attendance{
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

			app.GetAttendanceHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if !tt.expectError {
				var resp []app.Attendance
				err := json.NewDecoder(w.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.NotEmpty(t, resp)

				found, exists := repo.myIsExistDateAndProgrammerID(prog.ID, attendance.Date)
				if found != nil && len(resp) > 0 {
					assert.True(t, exists)
					assert.Equal(t, found.ProgrammerID, resp[0].ProgrammerID)
					assert.Equal(t, found.Date.UTC().Format("2006-01-02"), resp[0].Date.UTC().Format("2006-01-02"))
					assert.Equal(t, found.CheckIn.UTC().Format("2006-01-02 15:04"), resp[0].CheckIn.UTC().Format("2006-01-02 15:04"))
					assert.Equal(t, found.CheckOut.UTC().Format("2006-01-02 15:04"), resp[0].CheckOut.UTC().Format("2006-01-02 15:04"))
				}
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

	attendance := app.Attendance{
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
		input          app.AttendanceInput
		expectedStatus int
		errorMessage   string
	}{
		{
			name:         "Valid update",
			programmerID: fmt.Sprint(prog.ID),
			input: app.AttendanceInput{
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
			input:          app.AttendanceInput{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "CheckIn > CheckOut",
			programmerID: fmt.Sprint(prog.ID),
			input: app.AttendanceInput{
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
			input: app.AttendanceInput{
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

			app.UpdateAttendanceHandler(w, req)

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

	attendance := app.Attendance{
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

			app.DeleteAttendanceHandler(w, req)

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

	attendance := app.Attendance{
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

			app.DeleteOneDayAttendanceHandler(w, req)

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
	prog2, err := repo.myCreateProgrammer("Yahya Sinwar")
	assert.NoError(t, err)

	today := time.Now()

	attendances := []app.Attendance{
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

			app.GetAllAttendanceHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response []app.Attendance
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

			app.GetGirinofReportHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var report app.GirinofReport
				err := json.NewDecoder(w.Body).Decode(&report)
				assert.NoError(t, err, "Expected valid JSON response")
				assert.Equal(t, tt.programmerID, report.ProgrammerID, "unexpected programmerID")
				assert.Equal(t, tt.expectedDelayMinutes, report.TotalDelayMinutes, "unexpected total delay minutes")
				assert.Equal(t, tt.expectedEarlyDepartureMins, report.TotalEarlyDepartures, "unexpected total early departure minutes")
			}
		})
	}
}

// Test GetMonthlyReportHandler
func TestGetMonthlyReportHandler(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer()
	assert.NoError(t, err)

	seedAttendanceData(t, repo, prog.ID)

	now := time.Now().UTC()
	startDate := now.AddDate(0, 0, -30).Truncate(24 * time.Hour)
	endDate := now.Truncate(24 * time.Hour)

	tests := []struct {
		name           string
		programmerID   string
		startDate      string
		endDate        string
		expectedStatus int
		expectedReport *app.MonthlySummary
		expectedErrMsg string
	}{
		{
			name:           "Valid report",
			programmerID:   fmt.Sprint(prog.ID),
			startDate:      startDate.Format("2006-01-02"),
			endDate:        endDate.Format("2006-01-02"),
			expectedStatus: http.StatusOK,
			expectedReport: &app.MonthlySummary{
				ProgrammerID:         fmt.Sprint(prog.ID),
				TotalDaysPresent:     30,
				TotalOvertimeMinutes: 550,
				TotalDelayMinutes:    380,
				TotalEarlyDepartures: 320,
			},
		},
		{
			name:           "Invalid start date format",
			programmerID:   fmt.Sprint(prog.ID),
			startDate:      "invalid-start-date",
			endDate:        endDate.Format("2006-01-02"),
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Invalid start date format",
		},
		{
			name:           "Invalid end date format",
			programmerID:   fmt.Sprint(prog.ID),
			startDate:      startDate.Format("2006-01-02"),
			endDate:        "invalid-end-date",
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Invalid end date format",
		},
		{
			name:           "End date before start date",
			programmerID:   fmt.Sprint(prog.ID),
			startDate:      endDate.Format("2006-01-02"),
			endDate:        startDate.Format("2006-01-02"),
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "End date cannot be before start date",
		},
		{
			name:           "Missing start date",
			programmerID:   fmt.Sprint(prog.ID),
			startDate:      "",
			endDate:        endDate.Format("2006-01-02"),
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Start date is required",
		},
		{
			name:           "Missing end date",
			programmerID:   fmt.Sprint(prog.ID),
			startDate:      startDate.Format("2006-01-02"),
			endDate:        "",
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "End date is required",
		},
		{
			name:           "Missing programmer ID",
			programmerID:   "",
			startDate:      startDate.Format("2006-01-02"),
			endDate:        endDate.Format("2006-01-02"),
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Programmer ID is required",
		},
		{
			name:           "Invalid programmer ID",
			programmerID:   "invalid-id",
			startDate:      startDate.Format("2006-01-02"),
			endDate:        endDate.Format("2006-01-02"),
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Invalid programmer ID",
		},
		{
			name:           "Programmer not found",
			programmerID:   "99999",
			startDate:      startDate.Format("2006-01-02"),
			endDate:        endDate.Format("2006-01-02"),
			expectedStatus: http.StatusNotFound,
			expectedErrMsg: "Programmer not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/report/monthly/%s/%s/%s",
				tt.programmerID, tt.startDate, tt.endDate), nil)
			w := httptest.NewRecorder()

			app.GetMonthlyReportHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var report app.MonthlySummary
				err := json.NewDecoder(w.Body).Decode(&report)
				assert.NoError(t, err, "Failed to decode response body")

				assert.Equal(t, tt.expectedReport.ProgrammerID, report.ProgrammerID, "unexpected programmerID")
				assert.Equal(t, tt.expectedReport.TotalDaysPresent, report.TotalDaysPresent, "unexpected total day present")
				assert.Equal(t, tt.expectedReport.TotalOvertimeMinutes, report.TotalOvertimeMinutes, "unexpected total overtime minutes.")
				assert.Equal(t, tt.expectedReport.TotalDelayMinutes, report.TotalDelayMinutes, "unexpected total delay minutes")
				assert.Equal(t, tt.expectedReport.TotalEarlyDepartures, report.TotalEarlyDepartures, "unexpected total early departure")
			} else {
				responseBody := strings.TrimSpace(w.Body.String())
				assert.Contains(t, responseBody, tt.expectedErrMsg, "unexpected error message")
			}
		})
	}
}

// Test GetSalaryHandler
func TestGetSalaryHandler(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer()
	assert.NoError(t, err)

	prog2, err := repo.myCreateProgrammer()
	assert.NoError(t, err)

	seedAttendanceData(t, repo, prog.ID)

	expectedSalary := (float64(30) * app.DailyRate) + (float64(550) * app.OvertimeRate) - (float64(380) * app.DelayPenalty)

	tests := []struct {
		name           string
		programmerID   string
		expectedStatus int
		expectedReport *app.SalaryReport
		expectedErrMsg string
	}{
		{
			name:           "Valid salary calculation",
			programmerID:   fmt.Sprint(prog.ID),
			expectedStatus: http.StatusOK,
			expectedReport: &app.SalaryReport{
				Name:                 "Younes Mahmoudi",
				TotalDaysPresent:     30,
				TotalOvertimeMinutes: 550,
				TotalDelayMinutes:    380,
				TotalEarlyDepartures: 320,
				TotalSalary:          expectedSalary,
			},
		},
		{
			name:           "Missing programmer ID",
			programmerID:   "",
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Programmer ID is required",
		},
		{
			name:           "Invalid programmer ID",
			programmerID:   "invalid-id",
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Invalid programmer ID",
		},
		{
			name:           "Programmer not found",
			programmerID:   "999999",
			expectedStatus: http.StatusNotFound,
			expectedErrMsg: "Programmer not found",
		},
		{
			name:           "Zero attendance records",
			programmerID:   fmt.Sprint(prog2.ID),
			expectedStatus: http.StatusOK,
			expectedReport: &app.SalaryReport{
				Name:                 "",
				TotalDaysPresent:     0,
				TotalOvertimeMinutes: 0,
				TotalDelayMinutes:    0,
				TotalEarlyDepartures: 0,
				TotalSalary:          0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/salary/"+tt.programmerID, nil)
			w := httptest.NewRecorder()

			app.GetSalaryHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var report app.SalaryReport
				err := json.NewDecoder(w.Body).Decode(&report)
				assert.NoError(t, err, "Failed to decode response body")

				if tt.expectedReport.TotalDaysPresent > 0 {
					assert.Equal(t, tt.expectedReport.Name, report.Name, "unexpected name")
					assert.Equal(t, tt.expectedReport.TotalDaysPresent, report.TotalDaysPresent, "unexpected total days present")
					assert.Equal(t, tt.expectedReport.TotalOvertimeMinutes, report.TotalOvertimeMinutes, "unexpected total overtimei minutes")
					assert.Equal(t, tt.expectedReport.TotalDelayMinutes, report.TotalDelayMinutes, "unexpected total delay minutes")
					assert.Equal(t, tt.expectedReport.TotalEarlyDepartures, report.TotalEarlyDepartures, "unexpected total early departure")
					assert.Equal(t, tt.expectedReport.TotalSalary, report.TotalSalary, "unexpected total salary")
				} else {
					assert.Zero(t, report.TotalDaysPresent, "total day present should be zero")
					assert.Zero(t, report.TotalOvertimeMinutes, "total overtime minutes should be zero")
					assert.Zero(t, report.TotalDelayMinutes, "total delay minutes should be zero")
					assert.Zero(t, report.TotalSalary, "total salary should be zero")
				}
			} else {
				responseBody := strings.TrimSpace(w.Body.String())
				assert.Contains(t, responseBody, tt.expectedErrMsg, "unexpected error message")

			}

		})
	}
}

func TestRepository_CreateProgrammer(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	prog, err := app.Repo.CreateProgrammer("Test Developer")
	assert.NoError(t, err)
	assert.NotNil(t, prog)
	if prog != nil {
		assert.NotZero(t, prog.ID)
		assert.Equal(t, "Test Developer", prog.Name)
	}
}

func TestRepository_GetProgrammerByID(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer("Yahya Sinwar")
	assert.NoError(t, err)

	found, err := app.Repo.GetProgrammerByID(prog.ID)
	assert.NoError(t, err)
	if found != nil {
		assert.Equal(t, prog.ID, found.ID)
		assert.Equal(t, prog.Name, found.Name)
	}
	_, err = app.Repo.GetProgrammerByID(99999)
	assert.Error(t, err)
}

func TestRepository_CreateAttendance(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer("Yahya Sinwar")
	assert.NoError(t, err)

	attendance := &app.Attendance{
		ProgrammerID: prog.ID,
		Date:         time.Now(),
		CheckIn:      time.Now(),
		CheckOut:     time.Now().Add(8 * time.Hour),
	}

	err = app.Repo.CreateAttendance(attendance)
	assert.NoError(t, err)
	assert.NotZero(t, attendance.ID)

	err = app.Repo.CreateAttendance(attendance)
	assert.Error(t, err)
}

func TestRepository_GetAttendancesByProgrammer(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer("Yahya Sinwar")
	assert.NoError(t, err)

	attendance := &app.Attendance{
		ProgrammerID: prog.ID,
		Date:         time.Now(),
		CheckIn:      time.Now(),
		CheckOut:     time.Now().Add(8 * time.Hour),
	}
	repo.myCreateAttendance(attendance)

	records, err := app.Repo.GetAttendancesByProgrammer(prog.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, records)
	if len(records) > 0 {
		assert.Equal(t, prog.ID, records[0].ProgrammerID)
	}
}

func TestRepository_GetAllAttendance(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer("Yahya Sinwar")
	assert.NoError(t, err)

	attendance := &app.Attendance{
		ProgrammerID: prog.ID,
		Date:         time.Now(),
		CheckIn:      time.Now(),
		CheckOut:     time.Now().Add(8 * time.Hour),
	}
	repo.myCreateAttendance(attendance)

	records, err := repo.myGetAllAttendance()
	assert.NoError(t, err)
	assert.NotEmpty(t, records)
}

func TestRepository_UpdateAttendance(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer("Yahya Sinwar")
	assert.NoError(t, err)

	attendance := &app.Attendance{
		ProgrammerID: prog.ID,
		Date:         time.Now(),
		CheckIn:      time.Now(),
		CheckOut:     time.Now().Add(8 * time.Hour),
	}
	repo.myCreateAttendance(attendance)

	newCheckOut := time.Now().Add(9 * time.Hour)
	attendance.CheckOut = newCheckOut
	err = app.Repo.UpdateAttendance(attendance)
	assert.NoError(t, err)

	updated, _ := app.Repo.GetAttendancesByProgrammer(prog.ID)
	if len(updated) > 0 {
		assert.Equal(t, newCheckOut.Unix(), updated[0].CheckOut.Unix())
	}
}

func TestRepository_DeleteAttendance(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer("Yahya Sinwar")
	assert.NoError(t, err)

	attendance := &app.Attendance{
		ProgrammerID: prog.ID,
		Date:         time.Now(),
		CheckIn:      time.Now(),
		CheckOut:     time.Now().Add(8 * time.Hour),
	}
	repo.myCreateAttendance(attendance)

	err = app.Repo.DeleteAttendance(prog.ID)
	assert.NoError(t, err)

	records, _ := repo.myGetAttendancesByProgrammer(prog.ID)
	assert.Empty(t, records)
}

func TestRepository_DeleteOneDayAttendance(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer("Yahya Sinwar")
	assert.NoError(t, err)

	today := time.Now()
	attendance := &app.Attendance{
		ProgrammerID: prog.ID,
		Date:         today,
		CheckIn:      today,
		CheckOut:     today.Add(8 * time.Hour),
	}
	repo.myCreateAttendance(attendance)

	err = app.Repo.DeleteOneDayAttendance(prog.ID, today)
	assert.NoError(t, err)

	records, _ := repo.myGetAttendancesByProgrammer(prog.ID)
	assert.Empty(t, records)
}

func TestRepository_IsExistDateAndProgrammerID(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer("Yahya Sinwar")
	assert.NoError(t, err)

	today := time.Now()
	attendance := &app.Attendance{
		ProgrammerID: prog.ID,
		Date:         today,
		CheckIn:      today,
		CheckOut:     today.Add(8 * time.Hour),
	}
	repo.myCreateAttendance(attendance)

	found, exists := app.Repo.IsExistDateAndProgrammerID(prog.ID, today)
	assert.True(t, exists)
	assert.NotNil(t, found)

	found, exists = app.Repo.IsExistDateAndProgrammerID(prog.ID, today.AddDate(0, 0, 1))
	assert.False(t, exists)
	assert.Nil(t, found)
}

func TestRepository_GetGirinofReport(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer("Yahya Sinwar")
	assert.NoError(t, err)

	today := time.Now().UTC().Truncate(24 * time.Hour)
	attendance := &app.Attendance{
		ProgrammerID: prog.ID,
		Date:         today,
		CheckIn:      today.Add(10 * time.Hour),
		CheckOut:     today.Add(16 * time.Hour),
	}
	repo.myCreateAttendance(attendance)

	report, err := app.Repo.GetGirinofReport(prog.ID, today)
	assert.NoError(t, err)
	assert.NotNil(t, report)
	if report != nil {
		assert.NotZero(t, report.TotalDelayMinutes)
		assert.NotZero(t, report.TotalEarlyDepartures)
	}
}

func TestRepository_GetMonthlyReport(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer("Yahya Sinwar")
	assert.NoError(t, err)

	startDate := time.Now().AddDate(0, 0, -5)
	endDate := time.Now()

	for i := 0; i < 5; i++ {
		attendance := &app.Attendance{
			ProgrammerID: prog.ID,
			Date:         startDate.AddDate(0, 0, i),
			CheckIn:      startDate.AddDate(0, 0, i).Add(9 * time.Hour),
			CheckOut:     startDate.AddDate(0, 0, i).Add(18 * time.Hour),
		}
		repo.myCreateAttendance(attendance)
	}

	report, err := app.Repo.GetMonthlyReport(uint(prog.ID), startDate, endDate)
	assert.NoError(t, err)
	assert.NotNil(t, report)
	if report != nil {
		assert.Equal(t, 5, report.TotalDaysPresent)
		assert.NotZero(t, report.TotalOvertimeMinutes)
	}
}

func TestRepository_CalculateSalary(t *testing.T) {
	t.Cleanup(func() { clearTables(t) })

	repo := &MyTestRepository{DB: testDB}
	prog, err := repo.myCreateProgrammer("Yahya Sinwar")
	assert.NoError(t, err)

	startDate := time.Now().AddDate(0, 0, -5)
	endDate := time.Now()

	for i := 0; i < 5; i++ {
		attendance := &app.Attendance{
			ProgrammerID: prog.ID,
			Date:         startDate.AddDate(0, 0, i),
			CheckIn:      startDate.AddDate(0, 0, i).Add(9 * time.Hour),
			CheckOut:     startDate.AddDate(0, 0, i).Add(18 * time.Hour),
		}
		repo.myCreateAttendance(attendance)
	}

	report, err := app.Repo.CalculateSalary(uint(prog.ID), startDate, endDate)
	assert.NoError(t, err)
	assert.NotNil(t, report)
	if report != nil {
		assert.Equal(t, 5, report.TotalDaysPresent)
		assert.NotZero(t, report.TotalSalary)
		assert.NotZero(t, report.TotalOvertimeMinutes)

		expectedBase := float64(report.TotalDaysPresent) * app.DailyRate
		expectedOvertime := float64(report.TotalOvertimeMinutes) * app.OvertimeRate
		expectedDelay := float64(report.TotalDelayMinutes) * app.DelayPenalty
		expectedTotal := expectedBase + expectedOvertime - expectedDelay
		assert.Equal(t, expectedTotal, report.TotalSalary)
	}
}

// TestAttendance_UnmarshalJSON
func TestAttendance_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		errContains string
		expected    app.Attendance
	}{
		{
			name: "valid input",
			input: `{
				"date": "2024-01-01",
				"check_in": "2024-01-01 09:00:00",
				"check_out": "2024-01-01 17:00:00",
				"programmer_id": 1
			}`,
			wantErr: false,
			expected: app.Attendance{
				ProgrammerID: 1,
				Date:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				CheckIn:      time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				CheckOut:     time.Date(2024, 1, 1, 17, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "invalid date format",
			input: `{
				"date": "01-01-2024",
				"check_in": "2024-01-01 09:00:00",
				"check_out": "2024-01-01 17:00:00",
				"programmer_id": 1
			}`,
			wantErr:     true,
			errContains: "parsing time \"01-01-2024\" as \"2006-01-02\"",
		},
		{
			name: "invalid check_in format",
			input: `{
				"date": "2024-01-01",
				"check_in": "09:00:00",
				"check_out": "2024-01-01 17:00:00",
				"programmer_id": 1
			}`,
			wantErr:     true,
			errContains: "parsing time \"09:00:00\" as \"2006-01-02 15:04:05\"",
		},
		{
			name: "missing programmer_id",
			input: `{
				"date": "2024-01-01",
				"check_in": "2024-01-01 09:00:00",
				"check_out": "2024-01-01 17:00:00"
			}`,
			wantErr: false,
			expected: app.Attendance{
				ProgrammerID: 0,
				Date:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				CheckIn:      time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				CheckOut:     time.Date(2024, 1, 1, 17, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a app.Attendance
			err := json.Unmarshal([]byte(tt.input), &a)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" && err != nil {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ProgrammerID, a.ProgrammerID)
				assert.True(t, tt.expected.Date.Equal(a.Date), "Date mismatch")
				assert.True(t, tt.expected.CheckIn.Equal(a.CheckIn), "CheckIn mismatch")
				assert.True(t, tt.expected.CheckOut.Equal(a.CheckOut), "CheckOut mismatch")
			}
		})
	}
}

// TestAttendance_MarshalJSON
func TestAttendance_MarshalJSON(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	testTime := time.Date(2024, 1, 1, 9, 0, 0, 0, loc)

	tests := []struct {
		name     string
		input    app.Attendance
		expected string
	}{
		{
			name: "valid marshaling",
			input: app.Attendance{
				ID:           1,
				ProgrammerID: 1,
				Date:         testTime,
				CheckIn:      testTime,
				CheckOut:     testTime.Add(8 * time.Hour),
			},
			expected: `{
				"id":1,
				"programmer_id":1,
				"date":"2024-01-01",
				"check_in":"2024-01-01 00:00:00",
				"check_out":"2024-01-01 08:00:00"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			// Normalize JSON
			var expected, actual map[string]interface{}
			json.Unmarshal([]byte(tt.expected), &expected)
			json.Unmarshal(result, &actual)

			assert.Equal(t, expected, actual)
		})
	}
}

// TestAttendance_InputUnmarshalJSON
func TestAttendanceInput_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		errContains string
		expected    app.AttendanceInput
	}{
		{
			name: "valid input with all fields",
			input: `{
				"date": "2024-01-01",
				"check_in": "2024-01-01 09:00:00",
				"check_out": "2024-01-01 17:00:00"
			}`,
			wantErr: false,
			expected: app.AttendanceInput{
				Date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				CheckIn:  refTime(2024, 1, 1, 9, 0, 0),
				CheckOut: refTime(2024, 1, 1, 17, 0, 0),
			},
		},
		{
			name: "missing date field",
			input: `{
				"check_in": "2024-01-01 09:00:00",
				"check_out": "2024-01-01 17:00:00"
			}`,
			wantErr:     true,
			errContains: "date field is nessasery",
		},
		{
			name: "invalid date format",
			input: `{
				"date": "2024/01/01",
				"check_in": "2024-01-01 09:00:00"
			}`,
			wantErr:     true,
			errContains: "invalid date format",
		},
		{
			name: "invalid check_in format",
			input: `{
				"date": "2024-01-01",
				"check_in": "09:00:00"
			}`,
			wantErr:     true,
			errContains: "invalid check_in format",
		},
		{
			name: "optional check_out omitted",
			input: `{
				"date": "2024-01-01",
				"check_in": "2024-01-01 09:00:00"
			}`,
			wantErr: false,
			expected: app.AttendanceInput{
				Date:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				CheckIn: refTime(2024, 1, 1, 9, 0, 0),
			},
		},
		{
			name: "null check_in",
			input: `{
				"date": "2024-01-01",
				"check_in": null
			}`,
			wantErr: false,
			expected: app.AttendanceInput{
				Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ai app.AttendanceInput
			err := json.Unmarshal([]byte(tt.input), &ai)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" && err != nil {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {

				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Date.UTC(), ai.Date.UTC())
				assertTimesEqual(t, tt.expected.CheckIn, ai.CheckIn)
				assertTimesEqual(t, tt.expected.CheckOut, ai.CheckOut)

			}
		})
	}
}

// Helper functions
func refTime(year, month, day, hour, min, sec int) *time.Time {
	t := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
	return &t
}

func assertTimesEqual(t *testing.T, expected, actual *time.Time) {
	t.Helper()

	if expected == nil || actual == nil {
		assert.Equal(t, expected, actual)
		return
	}

	assert.True(t, expected.Equal(*actual), "expected: %v, actual: %v", expected, actual)
}
