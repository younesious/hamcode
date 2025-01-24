تیم همکاران‌سیستم تصمیم گرفته سیستم ردیابی ساعات ورود و خروج خود را با قابلیت شخصی‌سازی برای تیم *HR* طراحی کند. همکاران که به **شدت** شرکت **منظمی** است و روی نظم کارمندان خود نیز حساس است از شما می‌خواهد برای ورود و استخدام به این شرکت و کمک به تیم منابع انسانی، بک‌اند سیستم ترکینگ ورود و خروج برنامه‌نویسان برای تهیه گزارش‌ها و محاسبه حقوق‌شان را پیاده‌سازی کنید.

# جزئیات پروژه

پروژه اولیه را از %problem.initial_project% دانلود کنید. در این پروژه، شما باید قسمت‌های مشخص شده با `// TODO` را پیاده‌سازی کنید.

## مدل‌های پایگاه داده

ساختار جداول پروژه به شرح زیر باشد:


1- جدول `Programmer`

|    Column Name    |    Type   |   Notes    |
|:------------------:|:------------------:|:-----------:|
| `ID` |  `uint` |`PRIMARY KEY AUTO_INCREMENT`|
| `Name` |`string` | Programmer's name |
| `CreatedAt` |  `time.Time` | Creation timestamp |
| `UpdatedAt` |  `time.Time` | Last update timestamp |




2- جدول `Attendance`

|    Column Name    |    Type   |   Notes    |
|:------------------:|:------------------:|:-----------:|
| `ID` |  `uint` |`PRIMARY KEY AUTO_INCREMENT`|
| `ProgrammerID` |  `uint` | Foreign key to `programmers` table |
| `Date` |  `time.Time` | Attendance date |
| `CheckIn` |  `time.Time` | Check-in time |
| `CheckOut` |  `time.Time` | Check-out time |
| `CreatedAt` |  `time.Time` | Creation timestamp |
| `UpdatedAt` |  `time.Time` | Last update timestamp |


## روابط

روابط بین مدل‌ها باید به شرح زیر پیاده‌سازی شوند:

+ هر برنامه‌نویس می‌تواند دارای چندین رکورد ورود و خروج باشد.
+ هر رکورد ورود و خروج متعلق به یک برنامه‌نویس است.
+ هر برنامه‌نویس در هر روز فقط باید یک رکورد ورود و خروج در سیستم ثبت نماید.


### توابع `Marshal/Unmarshal` جیسون

این توابع باید تبدیل صحیح بین *JSON* و ساختارهای *Go* برای مدل `Attendance` و `AttendanceInput` را مدیریت کنند:


```go
func (a *Attendance) UnmarshalJSON(data []byte) error
func (a Attendance) MarshalJSON() ([]byte, error)
func (ai *AttendanceInput) UnmarshalJSON(data []byte) error
```

<details class="blue">
<summary>
متد `UnmarshalJSON` برای `Attendance`
</summary>

```go
func (a *Attendance) UnmarshalJSON(data []byte) error
```
+ داده‌های JSON را به ساختار `Attendance` تبدیل می‌کند.

+ فیلدهای `date`، `check_in` و `check_out` از نوع رشته (`string`) دریافت می‌شوند و به نوع `time.Time` تبدیل می‌شوند.

+ فرمت‌های مورد انتظار:
	+ `date`: `YYYY-MM-DD`

	+ `check_in`: `YYYY-MM-DD HH:MM:SS`

	+ `check_out`: `YYYY-MM-DD HH:MM:SS`

+ در صورت خطا در تبدیل فرمت‌ها، خطا برگردانده می‌شود.

</details>


<details class="blue">
<summary>
متد `MarshalJSON` برای `Attendance`
</summary>

```go
func (a Attendance) MarshalJSON() ([]byte, error)
```
+ ساختار `Attendance` را به JSON تبدیل می‌کند.

+ فیلدهای `date`، `check_in` و `check_out` از نوع `time.Time` به رشته (`string`) تبدیل می‌شوند.

+ فرمت‌های خروجی:
	+ `date`: `YYYY-MM-DD`

	+ `check_in`: `YYYY-MM-DD HH:MM:SS`

	+ `check_out`: `YYYY-MM-DD HH:MM:SS`

+ در صورت خطا در تبدیل، خطا برگردانده می‌شود.

</details>

<details class="blue">
<summary>
متد `UnmarshalJSON` برای `AttendanceInput`
</summary>

```go
func (i *AttendanceInput) UnmarshalJSON(data []byte) error
```
+ داده‌های JSON را به ساختار `AttendanceInput` تبدیل می‌کند.

+ فیلد `date` اجباری است و باید به فرمت `YYYY-MM-DD` باشد.

+ فیلدهای `check_in` و `check_out` اختیاری هستند و اگر وجود داشته باشند، باید به فرمت `YYYY-MM-DD HH:MM:SS` باشند.

+ در صورت خطا در تبدیل فرمت‌ها یا عدم وجود فیلد `date`، خطا برگردانده می‌شود.

</details>

### هندلرها

شما باید هندلرهای *HTTP* اندپوینت‌های زیر را با مشخصات دقیق آن‌ها که در ادامه بیان خواهد شد پیاده‌سازی کنید:

```go
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
```

<details class="blue">
<summary>
هندلر `CreateAttendanceHandler`
</summary>

+ متد: `POST`
+ اندپوینت یا مسیر: `/attendance`
+ ورودی: بدنه *JSON* شامل:

```json
{
  "programmer_id": uint,
  "date": "YYYY-MM-DD",
  "check_in": "YYYY-MM-DD HH:MM:SS",
  "check_out": "YYYY-MM-DD HH:MM:SS"
}
```
اعتبارسنجی:

+ در صورت دیکد نشدن *JSON* باید `Invalid input format` را با وضعیت 400 برگردانید.
+ در صورت عدم وجود فیلد تاریخ یا `date` باید `Missing required field: date` را با وضعیت 400 برگردانید.
+ در صورتی که `programmer_id` ارسال‌شده در دیتایس موجود نبود باید پیام`Programmer not found! Please hire the programmer first`. را با وضعیت 404 برگردانید.

پاسخ موفقیت‌آمیز:

+ وضعیت: 201 (ایجاد شده)
+ `Content-Type`: `application/json`
+ `Body`: شی `attendance` ایجاد شده

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to create attendance:`  + پیام خطا

</details>


<details class="blue"> <summary> هندلر `GetAttendanceHandler` </summary>

-   متد: `GET`

-   اندپوینت یا مسیر: `/attendance/{programmer_id}`

-   ورودی: `programmer_id` به عنوان پارامتر مسیر

اعتبارسنجی:

-   اگر `programmer_id` نامعتبر باشد (مثلاً عددی نباشد یا کمتر از ۱ باشد)، باید `Invalid Programmer ID` را با وضعیت 400 برگردانید.

-   اگر هیچ رکورد حضور و غیابی برای `programmer_id` داده‌شده وجود نداشته باشد، باید `No attendance records found for the given Programmer ID` را با وضعیت 404 برگردانید.

پاسخ موفقیت‌آمیز:

-   وضعیت: 200 (موفق)

-   `Content-Type`: `application/json`

-   `Body`: لیستی از رکوردهای حضور و غیاب مربوط به `programmer_id`

پاسخ خطا:

-   وضعیت: 500 با پیام `Error retrieving attendance records` + پیام خطا

</details>

<details class="blue">
<summary>
هندلر `UpdateAttendanceHandler`
</summary>

+ متد: `PUT`

+ اندپوینت یا مسیر: `/attendance/{programmer_id}`

+ ورودی: بدنه _JSON_ شامل:

```json
{
  "date": "YYYY-MM-DD",
  "check_in": "YYYY-MM-DD HH:MM:SS",
  "check_out": "YYYY-MM-DD HH:MM:SS"
}
```
اعتبارسنجی:

+ اگر `programmer_id` نامعتبر باشد (مثلاً عددی نباشد یا کمتر از ۱ باشد)، باید `Invalid ID` را با وضعیت 400 برگردانید.

+ اگر `programmer_id` در دیتابیس وجود نداشته باشد، باید `Programmer not found` را با وضعیت 404 برگردانید.

+ اگر `date` ارسال‌شده برای `programmer_id` موجود نباشد، باید `Programmer with the given date not recorded. Please use CreateAttendanceHandler to create it!` را با وضعیت 400 برگردانید.

پاسخ موفقیت‌آمیز:

+ وضعیت: 200 (موفق)

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to update attendance` + پیام خطا

</details>


<details class="blue">
<summary>
هندلر `DeleteAttendanceHandler`
</summary>

+ متد: `DELETE`

+ اندپوینت یا مسیر: `/attendance/{programmer_id}`

+ ورودی: `programmer_id` به عنوان پارامتر مسیر

اعتبارسنجی:

+ اگر `programmer_id` نامعتبر باشد (مثلاً عددی نباشد یا کمتر از ۱ باشد)، باید `Invalid programmer ID` را با وضعیت 400 برگردانید.

پاسخ موفقیت‌آمیز:

+ وضعیت: 202 (پذیرفته‌شده)

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to delete attendance` + پیام خطا

</details>

<details class="blue">
<summary>
هندلر `DeleteOneDayAttendanceHandler`
</summary>

+ متد: `DELETE`

+ اندپوینت یا مسیر: `/attendance/{programmer_id}/{date}`

+ ورودی:
	+ `programmer_id` به عنوان پارامتر مسیر

	+ `date` به عنوان پارامتر مسیر (فرمت: `YYYY-MM-DD`)

اعتبارسنجی:

+ اگر `programmer_id` نامعتبر باشد (مثلاً عددی نباشد یا کمتر از ۱ باشد)، باید `Invalid programmer ID` را با وضعیت 400 برگردانید.

+ اگر `date` فرمت نامعتبر داشته باشد، باید `Invalid Date` را با وضعیت 400 برگردانید.

پاسخ موفقیت‌آمیز:

+ وضعیت: 202 (پذیرفته‌شده)

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to delete attendance` + پیام خطا

</details>

<details class="blue">
<summary>
هندلر `GetAllAttendanceHandler`
</summary>

+ متد: `GET`

+ اندپوینت یا مسیر: `/attendance`

+ ورودی: ندارد

اعتبارسنجی:

+ نیاز به اعتبارسنجی خاصی ندارد (همه رکوردها را برمی‌گرداند)

پاسخ موفقیت‌آمیز:

+ وضعیت: 200 (موفق)

+ `Content-Type`: `application/json`

+ `Body`: لیست کامل تمامی رکوردهای حضور و غیاب

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to fetch attendance records` + پیام خطا

</details>

<details class="blue">
<summary>
هندلر `GetGirinofReportHandler`
</summary>

+ متد: `GET`

+ اندپوینت یا مسیر: `/report/girinof/{programmer_id}/{date}`

+ ورودی:
	+ `programmer_id` به عنوان پارامتر مسیر

	+ `date` به عنوان پارامتر مسیر (فرمت: `YYYY-MM-DD`)

اعتبارسنجی:

+ اگر `programmer_id` یا `date` ارسال‌نشده باشد، باید `Programmer ID is required` یا `Invalid Date` را با وضعیت 400 برگردانید.

پاسخ موفقیت‌آمیز:

+ وضعیت: 200 (موفق)

+ `Content-Type`: `application/json`

+ `Body`: گزارش `GirinofReport`

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to generate report` + پیام خطا

</details>

<details class="blue">
<summary>
هندلر `GetMonthlyReportHandler`
</summary>

+ متد: `GET`

+ اندپوینت یا مسیر: `/report/monthly/{programmer_id}/{start_date}/{end_date}`

+ ورودی:
	+ `programmer_id` به عنوان پارامتر مسیر

	+ `start_date` و `end_date` به عنوان پارامتر مسیر (فرمت: `YYYY-MM-DD`)

اعتبارسنجی:

+ اگر `programmer_id`، `start_date` یا `end_date` ارسال‌نشده باشد، باید `Programmer ID is required` یا `Invalid Date` را با وضعیت 400 برگردانید.

+ اگر `end_date` قبل از `start_date` باشد، باید `End date must be after start date` را با وضعیت 400 برگردانید.

پاسخ موفقیت‌آمیز:

+ وضعیت: 200 (موفق)

+ `Content-Type`: `application/json`

+ `Body`: گزارش ماهانه

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to generate monthly report` + پیام خطا

</details>


<details class="blue">
<summary>
هندلر `GetSalaryHandler`
</summary>

+ متد: `GET`

+ اندپوینت یا مسیر: `/salary/{programmer_id}`

+ ورودی: `programmer_id` به عنوان پارامتر مسیر

اعتبارسنجی:

+ اگر `programmer_id` ارسال‌نشده باشد، باید `Programmer ID is required` را با وضعیت 400 برگردانید.

پاسخ موفقیت‌آمیز:

+ وضعیت: 200 (موفق)

+ `Content-Type`: `application/json`

+ `Body`: حقوق محاسبه‌شده برای ۳۰ روز گذشته

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to calculate salary` + پیام خطا

</details>

<details class="blue">
<summary>
هندلر `CreateProgrammerHandler`
</summary>

+ متد: `POST`

+ اندپوینت یا مسیر: `/programmer`

+ ورودی: بدنه _JSON_ شامل:

```json
{
  "name": "string"
}
```
اعتبارسنجی:

+ اگر `name` ارسال‌نشده باشد، باید `Invalid input format` را با وضعیت 400 برگردانید.

پاسخ موفقیت‌آمیز:

+ وضعیت: 201 (ایجاد شده)

+ `Content-Type`: `application/json`

+ `Body`: شی `programmer` ایجاد شده

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to create programmer` + پیام خطا

</details>


### توابع Repository

متدهای *Repository* زیر را پیاده‌سازی کنید:

```go
func (repo *Repository) CreateAttendance(attendance *Attendance) error
func (repo *Repository) GetAttendancesByProgrammer(pid uint) ([]Attendance, error)
func (repo *Repository) GetAllAttendance() ([]Attendance, error)
func (repo *Repository) UpdateAttendance(attendance *Attendance) error
func (repo *Repository) DeleteAttendance(pid uint) error
func (repo *Repository) DeleteOneDayAttendance(pid uint, date time.Time) error
func (repo *Repository) GetGirinofReport(id string, date time.Time) (*GirinofReport, error)
func (repo *Repository) GetMonthlyReport(id string, startDate, endDate time.Time) (*MonthlySummary, error)
func (repo *Repository) CalculateSalary(id string, startDate, endDate time.Time) (*SalaryReport, error)
func (repo *Repository) CreateProgrammer(name string) (*Programmer, error)
func (repo *Repository) GetProgrammerByID(id uint) (*Programmer, error)
func (repo *Repository) IsExistDateAndProgrammerID(id uint, date time.Time) (*Attendance, bool)
```

<details class="blue">
<summary>
تابع `CreateAttendance`
</summary>

```go
func (repo *Repository) CreateAttendance(attendance *Attendance) error
```

+ رکورد  `Attendance` جدیدی ایجاد می‌کند.
+ اگر عملیات `Create` ناموفق باشد، خطا را برمی‌گرداند.

</details>


<details class="blue">
<summary>
متد `GetAttendancesByProgrammer`
</summary>
‍‍‍
```go
func (repo *Repository) GetAttendancesByProgrammer(pid uint) ([]Attendance, error)
```

+ لیست تمامی رکوردهای حضور و غیاب مربوط به یک برنامه‌نویس خاص را برمی‌گرداند.

+ رکوردها بر اساس تاریخ (`date`) به صورت نزولی مرتب می‌شوند.

+ حداکثر ۱۲ رکورد بازگردانده می‌شود.

+ در صورت خطا، خطا را برمی‌گرداند.

</details>


<details class="blue">
<summary>
متد `GetAllAttendance`
</summary>

```go
func (repo *Repository) GetAllAttendance() ([]Attendance, error)
```
+ لیست تمامی رکوردهای حضور و غیاب موجود در پایگاه داده را برمی‌گرداند.

+ در صورت خطا، خطا را برمی‌گرداند.

</details>

<details class="blue">
<summary>
متد `UpdateAttendance`
</summary>

```go
func (repo *Repository) UpdateAttendance(attendance *Attendance) error
```
+ رکورد حضور و غیاب موجود را به‌روزرسانی می‌کند.

+ فیلدهای `date`، `check_in` و `check_out` به‌روزرسانی می‌شوند.

+ در صورت خطا، خطا را برمی‌گرداند.

</details>

<details class="blue">
<summary>
متد `DeleteAttendance`
</summary>

```go
func (repo *Repository) DeleteAttendance(pid uint) error
```
+ تمامی رکوردهای حضور و غیاب مربوط به یک برنامه‌نویس خاص را حذف می‌کند.

+ در صورت خطا، خطا را برمی‌گرداند.

</details>

<details class="blue">
<summary>
متد `DeleteOneDayAttendance`
</summary>

```go
func (repo *Repository) DeleteOneDayAttendance(pid uint, date time.Time) error
```
+ رکورد حضور و غیاب مربوط به یک برنامه‌نویس خاص در یک تاریخ مشخص را حذف می‌کند.

+ در صورت خطا، خطا را برمی‌گرداند.

</details>

<details class="blue">
<summary>
متد `GetGirinofReport`
</summary>

```go

```
+ گزارش `GirinofReport` را برای یک برنامه‌نویس خاص در یک تاریخ مشخص ایجاد می‌کند.

+ این گزارش شامل مجموع دقیقه‌های تاخیر (`total_delay_minutes`) و مجموع دقیقه‌های خروج زودهنگام (`total_early_departure_minutes`) است.

+ در صورت خطا، خطا را برمی‌گرداند.

</details>

<details class="blue">
<summary>
متد `GetMonthlyReport`
</summary>

```go
func (repo *Repository) GetMonthlyReport(id string, startDate, endDate time.Time) (*MonthlySummary, error)
```
+ گزارش ماهانه (`MonthlySummary`) را برای یک برنامه‌نویس خاص در بازه زمانی مشخص ایجاد می‌کند.

+ این گزارش شامل تعداد روزهای حضور (`total_days_present`)، مجموع دقیقه‌های اضافه‌کاری (`total_overtime_minutes`)، مجموع دقیقه‌های تاخیر (`total_delay_minutes`) و مجموع دقیقه‌های خروج زودهنگام (`total_early_departure_minutes`) است.

+ در صورت خطا، خطا را برمی‌گرداند.

</details>

<details class="blue">
<summary>
متد `CalculateSalary`
</summary>

```go
func (repo *Repository) CalculateSalary(id string, startDate, endDate time.Time) (*SalaryReport, error)
```
+ حقوق (`SalaryReport`) یک برنامه‌نویس خاص را در بازه زمانی مشخص محاسبه می‌کند.

+ حقوق بر اساس تعداد روزهای حضور (`total_days_present`)، دقیقه‌های اضافه‌کاری (`total_overtime_minutes`) و جریمه تاخیر (`total_delay_minutes`) محاسبه می‌شود.

+ در صورت خطا، خطا را برمی‌گرداند.

</details>

<details class="blue">
<summary>
متد `CreateProgrammer`
</summary>

```go
func (repo *Repository) CreateProgrammer(name string) (*Programmer, error)
```
+ یک برنامه‌نویس جدید ایجاد می‌کند.

+ در صورت خطا، خطا را برمی‌گرداند.

</details>

<details class="blue">
<summary>
متد `GetProgrammerByID`
</summary>

```go
func (repo *Repository) GetProgrammerByID(id uint) (*Programmer, error)
```
+ اطلاعات یک برنامه‌نویس خاص را بر اساس `id` برمی‌گرداند.

+ در صورت خطا، خطا را برمی‌گرداند.

</details>

<details class="blue">
<summary>
متد `IsExistDateAndProgrammerID`
</summary>

```go
func (repo *Repository) IsExistDateAndProgrammerID(id uint, date time.Time) (*Attendance, bool)
```
+ بررسی می‌کند که آیا رکورد حضور و غیاب برای یک برنامه‌نویس خاص در تاریخ مشخص وجود دارد یا خیر.

+ اگر رکورد وجود داشته باشد، آن را برمی‌گرداند و `true` را به عنوان نتیجه برمی‌گرداند.

+ اگر رکورد وجود نداشته باشد، `nil` و `false` را برمی‌گرداند.

</details>


**نکات**

+ شما باید در این سوال از *GORM* و ویژگی‌های آن استفاده کنید.

+ فایل `db.go` برای اتصال *GORM* به پایگاه داده با استفاده از الگوی طراحی *singleton* است. شما نیازی به تغییر آن ندارید به جز تغییر ثابت‌ها برای تست روی محیط لوکال خود.

**چه چیزی را آپلود کنید**

پس از پیاده‌سازی ویژگی‌های خواسته‌شده،‌ کل پروژه را به صورت زیپ ارسال کنید.
فایل‌هایی که نیاز به ویرایش دارند:

+ `models.go‍`
+ `handlers.go`
+ `repository.go`

