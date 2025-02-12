تیم همکاران‌سیستم تصمیم گرفته سیستم ردیابی ساعات ورود و خروج خود را با قابلیت شخصی‌سازی برای تیم *HR* طراحی کند. همکاران که **به‌شدت** شرکت **منظمی** است و روی نظم کارمندان خود نیز حساس است از شما می‌خواهد برای ورود و استخدام به این شرکت و کمک به تیم منابع انسانی، بک‌اند سیستم ترکینگ ورود و خروج برنامه‌نویسان برای تهیه گزارش‌ها و محاسبه حقوق‌شان را پیاده‌سازی کنید.

# جزئیات پروژه

پروژه اولیه را از %problem.initial_project% دانلود کنید. ساختار پروژه به شکل زیر است: 

```shell
.
├── app
│   ├── db.go
│   ├── handlers.go   # TODO Implement
│   ├── models.go     # TODO Implement
│   ├── repository.go # TODO Implement
│   └── server.go
├── go.mod
├── go.sum
├── main.go
├── test
│   └── sample_test.go
```

در این پروژه، شما باید قسمت‌های مشخص شده با `// TODO` را پیاده‌سازی کنید. در نهایت شما از طریق اجرای فایل باینری `main.go` می‌توانید پروژه خود را اجرا کنید.

## مدل‌های پایگاه‌داده

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
| `ProgrammerID` |  `uint` | Foreign key to `Programmer` table |
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
+ داده‌های *JSON* را به ساختار `Attendance` تبدیل می‌کند.

+ فیلدهای `date`، `check_in` و `check_out` از نوع رشته (`string`) دریافت می‌شوند و به نوع `time.Time` تبدیل می‌شوند.

+ فرمت‌های مورد انتظار:
	+ `date`: `YYYY-MM-DD`

	+ `check_in`: `YYYY-MM-DD HH:MM:SS`

	+ `check_out`: `YYYY-MM-DD HH:MM:SS`

+ در صورت خطا در تبدیل فرمت‌ها، خطا برگردانده می‌شود.

**خطاهای ممکن:**

+ اگر `date` فرمت نامعتبر داشته باشد، خطای `parsing time` با پیام `expected YYYY-MM-DD` برگردانده می‌شود.

+ اگر `check_in` یا `check_out` فرمت نامعتبر داشته باشند، خطای `parsing time` با پیام `expected YYYY-MM-DD HH:MM:SS` برگردانده می‌شود.

**مثال ورودی معتبر:**
```json
{
  "date": "2024-01-01",
  "check_in": "2024-01-01 09:00:00",
  "check_out": "2024-01-01 17:00:00",
  "programmer_id": 1
}
```
**مثال ورودی نامعتبر:**

```json
{
  "date": "01-01-2024",
  "check_in": "09:00:00",
  "check_out": "2024-01-01 17:00:00",
  "programmer_id": 1
}
```

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
	+ `date`: `YYYY-MM-DD` در `UTC`

	+ `check_in`: `YYYY-MM-DD HH:MM:SS` در `UTC`

	+ `check_out`: `YYYY-MM-DD HH:MM:SS` در `UTC`

+ در صورت خطا در تبدیل، خطا برگردانده می‌شود.

**مثال خروجی:**

```json
{
  "id": 1,
  "programmer_id": 1,
  "date": "2024-01-01",
  "check_in": "2024-01-01 09:00:00",
  "check_out": "2024-01-01 17:00:00"
}
```

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

**خطاهای ممکن:**

+ اگر `date` وجود نداشته باشد، خطای `date field is necessary` برگردانده می‌شود.

+ اگر `date` فرمت نامعتبر داشته باشد، خطای `invalid date format, expected YYYY-MM-DD` برگردانده می‌شود.

+ اگر `check_in` یا `check_out` فرمت نامعتبر داشته باشند، خطای `invalid check_in/check_out format, expected YYYY-MM-DD HH:MM:SS` برگردانده می‌شود.

**مثال ورودی معتبر:**

```json
{
  "date": "2024-01-01",
  "check_in": "2024-01-01 09:00:00",
  "check_out": "2024-01-01 17:00:00"
}
```
**مثال ورودی نامعتبر:**

```json
{
  "date": "2024/01/01",
  "check_in": "09:00:00"
}
```

</details>


**نکته:** برای یکتا بودن `ProgrammerID` و `Date` در  `Attendance` می‌توانید از ایندکس کمک بگیرید.

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

توضیح:

این هندلر برای ثبت حضور و غیاب برنامه‌نویسان است. ابتدا داده‌های ورودی را از ریکوئست دریافت کرده و اعتبار آن‌ها را بررسی می‌کند (مانند فرمت درست، وجود فیلد تاریخ، و ترتیب ورود و خروج). در صورت معتبر بودن، اطلاعات حضور را در سیستم ذخیره کرده و پاسخ مناسب ارسال می‌کند.

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
+ در صورتی که ساعت خروج قبل از ساعت ورود بود، باید پیام `CheckOut cannot be before CheckIn` را با وضعیت 400 برگردانید.
+ در صورتی که `programmer_id` ارسال‌شده در پایگاه‌داده موجود نبود باید پیام `Programmer not found! Please hire the programmer first`. را با وضعیت 404 برگردانید.

پاسخ موفقیت‌آمیز:

+ وضعیت: 201 (ایجاد شده)
+ `Content-Type`: `application/json`
+ `Body`: شی `attendance` ایجاد شده

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to create attendance:`  + پیام خطا (اگر رکوردی با همان `programmer_id` و `date` وجود داشته باشد، باید `Duplicate attendance record` را برگرداند که شما این مورد را در لایه‌ی ریپوزیتوری هندل می‌کنید)

</details>


<details class="blue"> <summary> هندلر `GetAttendanceHandler` </summary>

توضیح:

این هندلر برای دریافت سوابق حضور و غیاب یک برنامه‌نویس خاص است. ابتدا شناسه برنامه‌نویس را از مسیر URL استخراج کرده و اعتبارسنجی می‌کند. سپس، اگر داده‌ای موجود باشد، سوابق حضور را در قالب JSON برمی‌گرداند، در غیر این صورت پیام خطای مناسب ارسال می‌شود.

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

توضیح:

این هندلر برای به‌روزرسانی اطلاعات حضور و غیاب یک برنامه‌نویس است. ابتدا شناسه برنامه‌نویس و تاریخ موردنظر را بررسی کرده و در صورت وجود، اطلاعات جدید (مانند ورود و خروج) را ثبت می‌کند. اگر داده‌ای نامعتبر باشد یا رکورد حضور قبلاً ثبت نشده باشد، پیام خطای مناسب ارسال می‌شود.

+ متد: `PUT`

+ اندپوینت یا مسیر: `/attendance/{programmer_id}`

+ ورودی: بدنه _JSON_ دو فیلد ورود و خروج به صورت **اختیاری**:

```json
{
  "date": "YYYY-MM-DD",
  "check_in": "YYYY-MM-DD HH:MM:SS",
  "check_out": "YYYY-MM-DD HH:MM:SS"
}
```
اعتبارسنجی:
 

+ اگر `programmer_id` نامعتبر باشد (مثلاً عددی نباشد یا کمتر از ۱ باشد)، باید `Invalid Programmer ID"` را با وضعیت 400 برگردانید.

+ اگر `programmer_id` در دیتابیس وجود نداشته باشد، باید `Programmer not found` را با وضعیت 404 برگردانید.

+ اگر فیلد date یا programmer_id ارسال نشوند باید پیام `Programmer ID is required` و `Date is required`  با وضعیت 400 برگردانید.

+ اگر `date` ارسال‌شده برای `programmer_id` موجود نباشد، باید `Programmer with the given date not recorded. Please use CreateAttendanceHandler to create it!` را با وضعیت 400 برگردانید.
+ اگر تاریخ خروج قبل از تاریخ ورود باشد، باید `CheckIn cannot be after CheckOut` را با وضعیت 400 برگردانید.

پاسخ موفقیت‌آمیز:

+ وضعیت: 200 (موفق)
+ `Content-Type`: `application/json`
+ `Body`: گزارش `GirinofReport` به شکل زیر:

```json
{
  "programmer_id": "1",
  "total_delays": 2,
  "total_early_departures": 3
}
```

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to update attendance` + پیام خطا

</details>


<details class="blue">
<summary>
هندلر `DeleteAttendanceHandler`
</summary>

توضیح:

این هندلر برای حذف سوابق حضور و غیاب یک برنامه‌نویس است. ابتدا شناسه برنامه‌نویس را بررسی کرده و در صورت معتبربودن، رکورد یا رکوردهای حضور مربوطه را حذف می‌کند. در صورت نامعتبربودن شناسه یا بروز خطا، پیام خطای مناسب ارسال می‌شود.

+ متد: `DELETE`

+ اندپوینت یا مسیر: `/attendance/{programmer_id}`

+ ورودی: `programmer_id` به عنوان پارامتر مسیر

اعتبارسنجی:

+ اگر `programmer_id` نامعتبر باشد (مثلاً عددی نباشد یا کمتر از ۱ باشد)، باید `Invalid programmer ID` را با وضعیت 400 برگردانید.
+ اگر برنامه‌نویسی با آیدی متناظر با `programmer_id` موجود نباشد، باید `Programmer not found` را با وضعیت 404 برگردانید.

+ اگر `date` در قالب `YYYY-MM-DD` نباشد، باید `Invalid date format` را با وضعیت **400** برگرداند. (این مورد را در بخش  توابع `Marshal/Unmarshal` جیسون هندل می‌کنید)

پاسخ موفقیت‌آمیز:

+ وضعیت: 202 (پذیرفته‌شده)

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to delete attendance` + پیام خطا

</details>

<details class="blue">
<summary>
هندلر `DeleteOneDayAttendanceHandler`
</summary>

توضیح:

این هندلر برای حذف سوابق حضور و غیاب یک برنامه‌نویس در یک تاریخ مشخص است. ابتدا شناسه برنامه‌نویس و تاریخ را از مسیر URL دریافت کرده و بررسی می‌کند. در صورت معتبربودن و وجود رکورد، آن را حذف می‌کند؛ در غیر این صورت، پیام خطای مناسب ارسال می‌شود.

+ متد: `DELETE`

+ اندپوینت یا مسیر: `/attendance/{programmer_id}/{date}`

+ ورودی:
	+ `programmer_id` به عنوان پارامتر مسیر

	+ `date` به عنوان پارامتر مسیر (فرمت: `YYYY-MM-DD`)

اعتبارسنجی:

+ اگر `programmer_id` نامعتبر باشد (مثلاً عددی نباشد یا کمتر از ۱ باشد)، باید `Invalid programmer ID` را با وضعیت 400 برگردانید.

+ اگر `date` فرمت نامعتبر داشته باشد، باید `Invalid Date` را با وضعیت 400 برگردانید.
+ اگر `programmer_id` یا `date` مورد نظر وجود نداشت، باید `No attendance record found for the given date` با وضعیت 404 را بازگردانید.

پاسخ موفقیت‌آمیز:

+ وضعیت: 202 (پذیرفته‌شده)

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to delete attendance` + پیام خطا

</details>

<details class="blue">
<summary>
هندلر `GetAllAttendanceHandler`
</summary>

توضیح:

این هندلر برای دریافت تمامی سوابق حضور و غیاب است. ابتدا تمام داده‌های حضور را از پایگاه‌داده واکشی کرده و در قالب JSON بازمی‌گرداند. در صورت بروز خطا در دریافت اطلاعات، پیام خطای مناسب ارسال می‌شود.

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

توضیح:

این هندلر برای دریافت گزارش حضور و غیاب یک برنامه‌نویس در یک تاریخ مشخص است. ابتدا شناسه برنامه‌نویس و تاریخ را بررسی کرده و در صورت وجود اطلاعات، گزارش را تولید و در قالب *JSON* بازمی‌گرداند. در صورت عدم وجود داده یا خطای پردازش، پیام خطای مناسب ارسال می‌شود.

+ متد: `GET`

+ اندپوینت یا مسیر: `/girinof/{programmer_id}/{date}`

+ ورودی:
	+ `programmer_id` به عنوان پارامتر مسیر

	+ `date` به عنوان پارامتر مسیر (فرمت: `YYYY-MM-DD`)

اعتبارسنجی:

+ اگر `programmer_id` ارسال‌ نشده باشد، باید `Programmer ID is required` را با وضعیت **400** برگرداند.
+ اگر `date` ارسال‌ نشده باشد، باید `Date is required` را با وضعیت **400** برگرداند.
+ اگر `date` فرمت نامعتبر داشته باشد، باید `Invalid Date` را با وضعیت **400** برگرداند.
+ اگر `programmer_id` عدد نباشد یا کمتر از ۱ باشد، باید `Invalid Programmer ID` را با وضعیت **400** برگرداند.
+ اگر برنامه‌نویس با `programmer_id` داده‌شده وجود نداشته باشد، باید `Programmer not found` را با وضعیت **404** برگرداند.
+ اگر رکورد حضور برای `programmer_id` و `date` داده‌شده وجود نداشته باشد، باید `Programmer with the given date not recorded. Please use CreateAttendanceHandler to create it!` را با وضعیت **400** برگرداند.


پاسخ موفقیت‌آمیز:

+ وضعیت: 200 (موفق)

+ `Content-Type`: `application/json`

+ `Body`: گزارش `GirinofReport`

+ محاسبات باید بر اساس اختلاف زمانی با ساعت استاندارد (۹:۰۰-۱۷:۰۰) انجام شود.

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to generate report` + پیام خطا

</details>

<details class="blue">
<summary>
هندلر `GetMonthlyReportHandler`
</summary>

توضیح:

این هندلر برای دریافت گزارش حضور و غیاب یک برنامه‌نویس در یک بازه زمانی مشخص است. ابتدا شناسه برنامه‌نویس، تاریخ شروع و پایان را بررسی کرده و در صورت معتبر بودن، گزارش را تولید و در قالب *JSON* بازمی‌گرداند. در صورت نامعتبر بودن ورودی‌ها یا بروز خطا، پیام خطای مناسب ارسال می‌شود.

+ متد: `GET`

+ اندپوینت یا مسیر: `/report/monthly/{programmer_id}/{start_date}/{end_date}`

+ ورودی:
	+ `programmer_id` به عنوان پارامتر مسیر

	+ `start_date` و `end_date` به عنوان پارامتر مسیر (فرمت: `YYYY-MM-DD`)

اعتبارسنجی:

+ اگر `programmer_id` ارسال‌ نشده باشد، باید `Programmer ID is required` را با وضعیت **400** برگرداند.
+ اگر `start_date` یا `end_date` ارسال‌ نشده باشد، باید `Start date is required` یا `End date is required` را با وضعیت **400** برگرداند.
+ اگر `programmer_id` عدد نباشد یا کمتر از ۱ باشد، باید `Invalid programmer ID` را با وضعیت **400** برگرداند.
+ اگر `start_date` یا `end_date` فرمت نامعتبر داشته باشند، باید `Invalid start date format` یا `Invalid end date format` را با وضعیت **400** برگرداند.
+ اگر `end_date` قبل از `start_date` باشد، باید `End date cannot be before start date` را با وضعیت **400** برگرداند.
+ اگر برنامه‌نویس با `programmer_id` داده‌شده وجود نداشته باشد، باید `Programmer not found` را با وضعیت **404** برگرداند.


پاسخ موفقیت‌آمیز:

+ وضعیت: 200 (موفق)

+ `Content-Type`: `application/json`

+ `Body`: گزارش `MonthlySummary` به شکل زیر:

```json
{
  "programmer_id": "1",
  "total_days_present": 4,
  "total_overtime_minutes": 3,
  "total_delay_minutes": 2,
  "total_early_departure_minutes": 1
}
```
+ محاسبات باید بر اساس اختلاف زمانی با ساعت استاندارد (۹:۰۰-۱۷:۰۰) انجام شود.

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to generate monthly report` + پیام خطا

</details>


<details class="blue">
<summary>
هندلر `GetSalaryHandler`
</summary>

توضیح:

این هندلر برای محاسبه حقوق یک برنامه‌نویس در ۳۰ روز گذشته است. ابتدا شناسه برنامه‌نویس را بررسی کرده و در صورت معتبربودن، حقوق او را بر اساس سوابق حضور و غیاب محاسبه و در قالب *JSON* بازمی‌گرداند. در صورت عدم وجود اطلاعات یا بروز خطا، پیام خطای مناسب ارسال می‌شود.

+ متد: `GET`

+ اندپوینت یا مسیر: `/salary/{programmer_id}`

+ ورودی: `programmer_id` به عنوان پارامتر مسیر

اعتبارسنجی:

+ اگر `programmer_id` ارسال‌ نشده باشد، باید `Programmer ID is required` را با وضعیت **400** برگرداند.
+ اگر `programmer_id` عدد نباشد یا کمتر از ۱ باشد، باید `Invalid programmer ID` را با وضعیت **400** برگرداند.
+ اگر برنامه‌نویس با `programmer_id` داده‌شده وجود نداشته باشد، باید `Programmer not found` را با وضعیت **404** برگرداند.


پاسخ موفقیت‌آمیز:

+ وضعیت: 200 (موفق)

+ `Content-Type`: `application/json`

+ `Body`: گزارش `SalaryReport` به شکل زیر:

```json
{
  "name": "Younes Sinwar",
  "total_days_present": 30,
  "total_overtime_minutes": 1,
  "total_delay_minutes": 2,
  "total_early_departure_minutes": 3,
  "total_salary": 3000.0
}
```
+ محاسبات باید بر اساس اختلاف زمانی با ساعت استاندارد (۹:۰۰-۱۷:۰۰) انجام شود.

محاسبه حقوق: حقوق بر اساس فرمول زیر محاسبه می‌شود:

    حقوق پایه = تعداد روزهای حضور × DailyRate (۱۰۰)

    اضافه‌کاری = مجموع دقیقه‌های اضافه‌کاری × OvertimeRate (۲۰)

    جریمه تاخیر = مجموع دقیقه‌های تاخیر × DelayPenalty (۱۰)

    حقوق نهایی = حقوق پایه + اضافه‌کاری - جریمه تاخیر

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to calculate salary` + پیام خطا

</details>

<details class="blue">
<summary>
هندلر `CreateProgrammerHandler`
</summary>

توضیح:

این هندلر برای ایجاد یک برنامه‌نویس جدید در سیستم است. ابتدا داده‌های ورودی را بررسی کرده و در صورت معتبربودن، برنامه‌نویس را در پایگاه‌داده ثبت می‌کند. در صورت بروز خطا یا نامعتبربودن اطلاعات، پیام خطای مناسب ارسال می‌شود.

+ متد: `POST`

+ اندپوینت یا مسیر: `/programmer`

+ ورودی: بدنه _JSON_ شامل:

```json
{
  "name": "string"
}
```
اعتبارسنجی:

+ اگر `name` ارسال‌ نشده باشد، باید `Invalid input format` را با وضعیت 400 برگردانید.
+ + اگر فیلد `name` خالی باشد، باید `Name is required` را با وضعیت 400 برگردانید.

پاسخ موفقیت‌آمیز:

+ وضعیت: 201 (ایجاد شده)

+ `Content-Type`: `application/json`

+ `Body`: شی `programmer` ایجادشده

پاسخ خطا:

+ وضعیت: 500 با پیام `Failed to create programmer` + پیام خطا

</details>

+ در تمام هندلرها اگر خطای داخلی سرور رخ دهد، باید وضعیت **500** برگردانده شود.

### توابع Repository

متدهای *Repository* زیر را پیاده‌سازی کنید:

```go
func (repo *Repository) CreateAttendance(attendance *Attendance) error
func (repo *Repository) GetAttendancesByProgrammer(pid uint) ([]Attendance, error)
func (repo *Repository) GetAllAttendance() ([]Attendance, error)
func (repo *Repository) UpdateAttendance(attendance *Attendance) error
func (repo *Repository) DeleteAttendance(pid uint) error
func (repo *Repository) DeleteOneDayAttendance(pid uint, date time.Time) error
func (repo *Repository) GetGirinofReport(pid uint, date time.Time) (*GirinofReport, error)
func (repo *Repository) GetMonthlyReport(pid uint, startDate, endDate time.Time) (*MonthlySummary, error)
func (repo *Repository) CalculateSalary(pid uint, startDate, endDate time.Time) (*SalaryReport, error)
func (repo *Repository) CreateProgrammer(name string) (*Programmer, error)
func (repo *Repository) GetProgrammerByID(id uint) (*Programmer, error)
func (repo *Repository) IsExistDateAndProgrammerID(id uint, date time.Time) (*Attendance, bool)
```

<details class="blue">
<summary>
متد `CreateAttendance`
</summary>

```go
func (repo *Repository) CreateAttendance(attendance *Attendance) error
```

+ رکورد `Attendance` جدیدی با تمام فیلدهای الزامی ایجاد می‌کند

+ در صورت وجود `attendance` تکراری (با همان `ProgrammerID` و `Date`) خطا برمی‌گرداند

+ آی‌دی تولید شده برای رکورد جدید باید غیرصفر باشد

+ در صورت خطا در عملیات پایگاه‌داده، خطا را برمی‌گرداند


</details>


<details class="blue">
<summary>
متد `GetAttendancesByProgrammer`
</summary>
‍‍‍
```go
func (repo *Repository) GetAttendancesByProgrammer(pid uint) ([]Attendance, error)
```

+ لیست تمام رکوردهای حضور و غیاب یک برنامه‌نویس را برمی‌گرداند

+ نتایج باید بر اساس تاریخ (`date`) به صورت نزولی مرتب شوند

+ حداکثر ۱۲ رکورد بازگردانده شود

+ در صورت عدم وجود رکورد، لیست خالی برمی‌گرداند

+ در صورت خطا در عملیات پایگاه‌داده، خطا را برمی‌گرداند


</details>


<details class="blue">
<summary>
متد `GetAllAttendance`
</summary>

```go
func (repo *Repository) GetAllAttendance() ([]Attendance, error)
```
+ لیست تمامی رکوردهای حضور و غیاب موجود در پایگاه‌داده را برمی‌گرداند.

+  در صورت عدم وجود رکورد، لیست خالی برمی‌گرداند.

+ در صورت خطا، خطا را برمی‌گرداند.

</details>

<details class="blue">
<summary>
متد `UpdateAttendance`
</summary>

```go
func (repo *Repository) UpdateAttendance(attendance *Attendance) error
```
+ رکورد موجود را با مقادیر جدید به‌روزرسانی می‌کند.

+ فیلدهای قابل به‌روزرسانی: `date`، `check_in` و `check_out`

+ در صورت عدم وجود رکورد مرتبط، خطا برمی‌گرداند.

+ در صورت خطا در عملیات پایگاه‌داده، خطا را برمی‌گرداند.


</details>

<details class="blue">
<summary>
متد `DeleteAttendance`
</summary>

```go
func (repo *Repository) DeleteAttendance(pid uint) error
```
+ تمام رکوردهای حضور و غیاب مربوط به یک برنامه‌نویس را حذف می‌کند.

+ پس از حذف، هیچ رکوردی با ProgrammerID مربوطه نباید وجود داشته باشد.

+ در صورت خطا در عملیات پایگاه‌داده، خطا را برمی‌گرداند.


</details>

<details class="blue">
<summary>
متد `DeleteOneDayAttendance`
</summary>

```go
func (repo *Repository) DeleteOneDayAttendance(pid uint, date time.Time) error
```
+ رکورد حضور و غیاب مربوط به یک روز خاص را حذف می‌کند.

+ تاریخ باید دقیقاً مطابقت داشته باشد (برای انجام این کار تمامی رکوردهای تاریخ آن روز را از ساعت `00:00:00` تا ساعت `24:00:00` حذف کنید. می‌توانید از `date.Truncate` استفاده کنید)

+ در صورت عدم وجود رکورد، خطا برمی‌گرداند.

+ در صورت خطا در عملیات پایگاه‌داده، خطا را برمی‌گرداند.

</details>

<details class="blue">
<summary>
متد `GetGirinofReport`
</summary>

```go
func (repo *Repository) GetGirinofReport(pid uint, date time.Time) (*GirinofReport, error)
```

+ گزارش روزانه شامل مجموع تاخیرها (`total_delay_minutes`) و خروج‌های زودهنگام (`total_early_departure_minutes`) را همراه با `programmer_id` ایجاد می‌کند.

+ محاسبات باید بر اساس اختلاف زمانی با ساعت استاندارد (۹:۰۰-۱۷:۰۰) انجام شود.

+ در صورت عدم وجود رکورد در تاریخ مشخص، خطا برمی‌گرداند.

+ در صورت خطا در محاسبات یا عملیات پایگاه‌داده، خطا را برمی‌گرداند.


</details>

<details class="blue">
<summary>
متد `GetMonthlyReport`
</summary>

```go
func (repo *Repository) GetMonthlyReport(id string, startDate, endDate time.Time) (*MonthlySummary, error)
```
+ گزارش ماهانه شامل آمار حضور، اضافه‌کاری و تاخیرها را در بازه زمانی داده شده ایجاد می‌کند

+ محاسبات باید بر اساس اختلاف زمانی با ساعت استاندارد (۹:۰۰-۱۷:۰۰) انجام شود.

+ این گزارش شامل آی‌دی برنامه‌نویس (`programmer_id`) تعداد روزهای حضور (`total_days_present`)، مجموع دقیقه‌های اضافه‌کاری (`total_overtime_minutes`)، مجموع دقیقه‌های تاخیر (`total_delay_minutes`) و مجموع دقیقه‌های خروج زودهنگام (`total_early_departure_minutes`) است.

+ در صورت خطا در محاسبات یا عملیات پایگاه‌داده، خطا را برمی‌گرداند.


</details>

<details class="blue">
<summary>
متد `CalculateSalary`
</summary>

```go
func (repo *Repository) CalculateSalary(id string, startDate, endDate time.Time) (*SalaryReport, error)
```
+ حقوق (`SalaryReport`) یک برنامه‌نویس خاص را در سی روز گذشته محاسبه می‌کند.

+ محاسبات باید بر اساس اختلاف زمانی با ساعت استاندارد (۹:۰۰-۱۷:۰۰) انجام شود.

+  حقوق خالص را بر اساس فرمول زیر محاسبه می‌کند:

	+ حقوق پایه = تعداد روزهای حضور × نرخ روزانه (`DailyRate`)

	+ اضافه‌کاری = مجموع دقیقه‌های اضافه‌کاری × نرخ اضافه‌کاری (`OvertimeRate`)

	+ جریمه تاخیر = مجموع دقیقه‌های تاخیر × نرخ جریمه (`DelayPenalty`)

	+ حقوق نهایی = حقوق پایه + اضافه‌کاری - جریمه تاخیر

+ این گزارش شامل نام برنامه‌نویس (`name`)، تعداد روزهای حضور (`total_days_present`)، مجموع دقیقه‌های اضافه‌کاری (`total_overtime_minutes`)، مجموع دقیقه‌های تاخیر (`total_delay_minutes`)، مجموع دقیقه‌های خروج زودهنگام (`total_early_departure_minutes`) و مجموع حقوق (`total_salary`) است.

+  در صورت خطا در محاسبات یا عملیات پایگاه‌داده، خطا را برمی‌گرداند

</details>

<details class="blue">
<summary>
متد `CreateProgrammer`
</summary>

```go
func (repo *Repository) CreateProgrammer(name string) (*Programmer, error)
```
+ یک برنامه‌نویس جدید ایجاد می‌کند.
+  آی‌دی تولید شده باید غیرصفر باشد.
+  در صورت خطا در عملیات پایگاه‌داده، خطا را برمی‌گرداند.

</details>

<details class="blue">
<summary>
متد `GetProgrammerByID`
</summary>

```go
func (repo *Repository) GetProgrammerByID(id uint) (*Programmer, error)
```
+ اطلاعات کامل برنامه‌نویس را بر اساس ID برمی‌گرداند.

+ در صورت عدم وجود برنامه‌نویس با ID داده شده، خطا برمی‌گرداند.

+ در صورت خطا در عملیات پایگاه‌داده، خطا را برمی‌گرداند.


</details>

<details class="blue">
<summary>
متد `IsExistDateAndProgrammerID`
</summary>

```go
func (repo *Repository) IsExistDateAndProgrammerID(id uint, date time.Time) (*Attendance, bool)
```
+ بررسی می‌کند که آیا رکورد حضور و غیاب برای یک برنامه‌نویس خاص در تاریخ مشخص وجود دارد یا خیر.

+ تاریخ باید دقیقاً مطابقت داشته باشد (برای انجام این کار تمامی رکوردهای تاریخ آن در آن روز یعنی از ساعت `00:00:00` تا ساعت `24:00:00` را بررسی کنید)

+  در صورت وجود، رکورد کامل را همراه با `true` برمی‌گرداند.

+  در صورت عدم وجود، `nil` و `false` برمی‌گرداند.

+  تاریخ باید دقیقاً مطابقت داشته باشد (با در نظر گرفتن ساعت ۰۰:۰۰:۰۰)


</details>


**نکات**

+ شما باید در این سوال از *GORM* و ویژگی‌های آن استفاده کنید و به جز **پکیج‌های موجود در پروژه اولیه** اجازه استفاده از **پکیج دیگری** را **ندارید**.

+ فایل `db.go` برای اتصال *GORM* به پایگاه‌داده با استفاده از الگوی طراحی *singleton* است. شما نیازی به تغییر آن ندارید به جز تغییر ثابت‌ها برای تست روی محیط لوکال خود.

+ فایل‌ها و کدهای پروژه‌اولیه را نباید تغییر بدهید و صرفا باید مکان‌هایی که با کامنت `// TODO Implement` برای شما مشخص شده را پیاده‌سازی کنید.

+ نیازی به ارسال فایل‌های `go.mod` و `go.sum` ندارید! 

**چه چیزی را آپلود کنید**

پس از پیاده‌سازی ویژگی‌های خواسته‌شده،‌ **فقط** دایرکتوری `app` که شامل سه فایل `handlers.go` و `models.go` و `repository.go` است را فشرده کرده و صورت زیپ ارسال کنید.
فایل‌هایی که نیاز به ویرایش دارند:

+ `app/models.go‍`
+ `app/handlers.go`
+ `app/repository.go`

