# Attendance Tracking System

This repository contains a backend solution for an attendance tracking system built in Go. The system is designed for HR teams to manage programmer attendance, generate detailed reports, and calculate salaries based on attendance records.

## Features

- **Programmer Management:**  
  Create new programmer records.
  
- **Attendance Management:**  
  Create, read, update, and delete attendance records with unique constraints (one record per programmer per day).

- **Reports:**  
  - **Daily Report (GirinofReport):** Calculates delays and early departures based on standard work hours (09:00–17:00).  
  - **Monthly Report (MonthlySummary):** Summarizes days present, overtime, delays, and early departures over a specified period.  
  - **Salary Calculation (SalaryReport):** Computes salary for the past 30 days using a formula based on daily rate, overtime rate, and delay penalty.

- **Custom JSON Handling:**  
  Custom `MarshalJSON`/`UnmarshalJSON` implementations ensure strict date and time formats:
  - `date`: `YYYY-MM-DD`
  - `check_in` / `check_out`: `YYYY-MM-DD HH:MM:SS`

## Project Structure

```shell
. 
├── app 
│ ├── db.go # Database connection using GORM (singleton pattern) 
│ ├── handlers.go # HTTP endpoint handlers (TODO Implement) 
│ ├── models.go # Data models and JSON conversion methods (TODO Implement) 
│ ├── repository.go # Repository methods for DB operations (TODO Implement) 
│ └── server.go 
├── go.mod 
├── go.sum 
├── main.go 
└── test 
    ├── main_test.go 
    └── sample_test.go
```

## Getting Started

### Prerequisites

- **Go:** Version 1.21 or higher  
- **Database:** A PostgreSQL database (configured in `app/db.go` via GORM)

### Installation

1. **Clone the Repository:**

```bash
git clone https://github.com/younesious/hamcode.git
cd hamcode
```

2.  **Configure the Database:**
    Update the connection settings in `app/db.go` as needed for your local environment.

3.  **Run the Application:**

```shell
go run main.go
```

The server will start and listen for HTTP requests.

### Running Tests

Execute the following command to run all tests:

```shell
go test ./...
```

API Endpoints
-------------

The backend exposes the following HTTP endpoints:

-   `POST /attendance`\
    Create a new attendance record.

-   `GET /attendance/{programmer_id}`\
    Retrieve all attendance records for a specific programmer.

-   `PUT /attendance/{programmer_id}`\
    Update an existing attendance record.

-   `DELETE /attendance/{programmer_id}`\
    Delete all attendance records for a programmer.

-   `DELETE /attendance/{programmer_id}/{date}`\
    Delete the attendance record for a specific date.

-   `GET /attendance`\
    Retrieve all attendance records.

-   `GET /girinof/{programmer_id}/{date}`\
    Generate a daily attendance report.

-   `GET /report/monthly/{programmer_id}/{start_date}/{end_date}`\
    Generate a monthly attendance summary.

-   `GET /salary/{programmer_id}`\
    Calculate salary based on attendance for the past 30 days.

-   `POST /programmer`\
    Create a new programmer record.

> **Important:** All date and time fields must follow the strict formats:
>
> -   `date`: `YYYY-MM-DD`
> -   `check_in` / `check_out`: `YYYY-MM-DD HH:MM:SS`

License
-------

This project is licensed under the MIT License.

Contact
-------

If you have any questions or encounter any issues, please open an issue in the repository.
