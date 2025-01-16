package main

import (
	"encoding/json"
	"time"
)

type Attendance struct {
	ID        uint      `gorm:"primaryKey"`
	Engineer  string    `gorm:"size:100;not null"`
	Date      time.Time `gorm:"not null"`
	CheckIn   time.Time
	CheckOut  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *Attendance) UnmarshalJSON(data []byte) error {
	type Alias Attendance // Create an alias to avoid recursion

	aux := &struct {
		Date     string `json:"date"`
		CheckIn  string `json:"check_in"`
		CheckOut string `json:"check_out"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	// First, decode the JSON into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Manually parse the date, check_in, and check_out fields
	var err error
	if a.Date, err = time.Parse("2006-01-02", aux.Date); err != nil {
		return err
	}
	if a.CheckIn, err = time.Parse("2006-01-02 15:04:05", aux.CheckIn); err != nil {
		return err
	}
	if a.CheckOut, err = time.Parse("2006-01-02 15:04:05", aux.CheckOut); err != nil {
		return err
	}

	return nil
}

type Input struct {
	Date     *time.Time `json:"date"`
	CheckIn  *time.Time `json:"check_in"`
	CheckOut *time.Time `json:"check_out"`
}

/*
func (i *Input) UnmarshalJSON(data []byte) error {
	type Alias Input // Create an alias to avoid recursion

	aux := &struct {
		Date     *string `json:"date"`
		CheckIn  *string `json:"check_in"`
		CheckOut *string `json:"check_out"`
		*Alias
	}{
		Alias: (*Alias)(i),
	}

	// First, decode the JSON into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Manually parse the date, check_in, and check_out fields
	var err error

	if aux.Date != nil {
		if a.Date, err = time.Parse("2006-01-02", *aux.Date); err != nil {
			return err
		}
	}
	if aux.CheckIn != nil {
		if a.CheckIn, err = time.Parse("2006-01-02 15:04:05", *aux.CheckIn); err != nil {
			return err
		}
	}
	if aux.CheckOut != nil {
		if a.CheckOut, err = time.Parse("2006-01-02 15:04:05", *aux.CheckOut); err != nil {
			return err
		}
	}

	return nil
}
*/

func (i *Input) UnmarshalJSON(data []byte) error {
	// Create an auxiliary struct to handle optional fields
	aux := struct {
		Date     *string `json:"date"`
		CheckIn  *string `json:"check_in"`
		CheckOut *string `json:"check_out"`
	}{}

	// First, unmarshal the raw JSON into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse the date fields manually if they are not nil
	if aux.Date != nil {
		// Parse date into time.Time
		parsedDate, err := time.Parse("2006-01-02", *aux.Date)
		if err != nil {
			return err
		}
		// Assign the parsed time to the pointer field
		i.Date = &parsedDate
	}
	if aux.CheckIn != nil {
		// Parse check_in into time.Time
		parsedCheckIn, err := time.Parse("2006-01-02 15:04:05", *aux.CheckIn)
		if err != nil {
			return err
		}
		// Assign the parsed time to the pointer field
		i.CheckIn = &parsedCheckIn
	}
	if aux.CheckOut != nil {
		// Parse check_out into time.Time
		parsedCheckOut, err := time.Parse("2006-01-02 15:04:05", *aux.CheckOut)
		if err != nil {
			return err
		}
		// Assign the parsed time to the pointer field
		i.CheckOut = &parsedCheckOut
	}

	return nil
}
