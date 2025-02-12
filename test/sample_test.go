package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"younesious/hamcode/app"
)

var sampleTestDB *gorm.DB

type MySampleTestRepository struct {
	DB *gorm.DB
}

func (r *MySampleTestRepository) mySampleCreateProgrammer(names ...string) (*app.Programmer, error) {
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

func sampleClearTables(t *testing.T) {
	t.Helper()

	tables := []string{"programmers", "attendances"}
	for _, table := range tables {
		err := sampleTestDB.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error
		assert.NoError(t, err)
	}
}

func TestSampleGetConnection(t *testing.T) {
	db := app.GetConnection()
	assert.NotNil(t, db)
}

func TestMain(m *testing.M) {
	var err error
	sampleTestDB = app.GetConnection()

	sampleTestDB.Migrator().DropTable(&app.Attendance{})
	sampleTestDB.AutoMigrate(&app.Attendance{})

	app.Repo = &app.Repository{DB: sampleTestDB}

	code := m.Run()

	sqlDB, err := sampleTestDB.DB()
	if err == nil {
		sqlDB.Close()
	}

	os.Exit(code)
}

// Test CreateProgrammerHandler
func TestSampleCreateProgrammerHandler(t *testing.T) {
	t.Cleanup(func() { sampleClearTables(t) })

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
				err = sampleTestDB.First(&dbProg, resp.ID).Error
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Name, dbProg.Name)

				var count int64
				sampleTestDB.Model(&app.Programmer{}).Count(&count)
				assert.Equal(t, int64(1), count)
			}
		})
	}
}

// TestRepository_CreateAttendance
func TestSampleRepository_CreateAttendance(t *testing.T) {
	t.Cleanup(func() { sampleClearTables(t) })

	repo := &MySampleTestRepository{DB: sampleTestDB}
	prog, err := repo.mySampleCreateProgrammer("Yahya Sinwar")
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

// TestAttendance_UnmarshalJSON
func TestSampleAttendance_UnmarshalJSON(t *testing.T) {
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
