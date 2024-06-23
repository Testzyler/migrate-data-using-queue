package main

// This script was created with the assistance of Chat-GPT

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Employee struct to hold employee data
type Employee struct {
	EmpNo     string    `gorm:"primaryKey"`
	BirthDate time.Time `gorm:"type:date"`
	FirstName string    `gorm:"type:varchar(14)"`
	LastName  string    `gorm:"type:varchar(16)"`
	Gender    string    `gorm:"type:char(1)"`
	HireDate  time.Time `gorm:"type:date"`
}

// Department struct to hold department data
type Department struct {
	DeptNo   string `gorm:"primaryKey;type:char(4)"`
	DeptName string `gorm:"type:varchar(40);unique"`
}

// DeptManager struct to hold department manager data
type DeptManager struct {
	DeptNo   string    `gorm:"primaryKey"`
	EmpNo    string    `gorm:"primaryKey"`
	FromDate time.Time `gorm:"primaryKey"`
	ToDate   time.Time `gorm:"not null"`
}

// DeptEmp struct to hold department employee data
type DeptEmp struct {
	EmpNo    string    `gorm:"primaryKey"`
	DeptNo   string    `gorm:"primaryKey;type:char(4)"`
	FromDate time.Time `gorm:"primaryKey"`
	ToDate   time.Time `gorm:"not null"`
}

// Title struct to hold employee title data
type Title struct {
	EmpNo    string    `gorm:"primaryKey"`
	Title    string    `gorm:"primaryKey"`
	FromDate time.Time `gorm:"primaryKey"`
	ToDate   time.Time
}

// Salary struct to hold employee salary data
type Salary struct {
	EmpNo    string    `gorm:"primaryKey"`
	Salary   int       `gorm:"not null"`
	FromDate time.Time `gorm:"primaryKey"`
	ToDate   time.Time `gorm:"not null"`
}

// Seed data generation functions
func generateRandomEmployees(n int64) []Employee {
	var employees []Employee
	empNoCounter := make(map[string]int)

	for i := int64(1); i <= n; i++ {
		gofakeit.Seed(time.Now().UnixNano() + i) // Use a unique seed for each employee

		hireDate := gofakeit.DateRange(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), time.Now())
		baseEmpNo := fmt.Sprintf("MQ%s", hireDate.Format("060102")) // Format: MQyymmdd
		// Check if this base EmpNo has been used before
		count, exists := empNoCounter[baseEmpNo]
		if !exists {
			empNoCounter[baseEmpNo] = 1
		} else {
			empNoCounter[baseEmpNo] = count + 1
		}

		// Append the counter to make the EmpNo unique
		empNo := fmt.Sprintf("%s-%03d", baseEmpNo, empNoCounter[baseEmpNo])
		emp := Employee{
			EmpNo:     empNo,
			BirthDate: gofakeit.DateRange(time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC)),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Gender:    gofakeit.RandomString([]string{"M", "F"}),
			HireDate:  hireDate,
		}
		employees = append(employees, emp)
	}
	return employees
}

func generateDepartments() []Department {
	return []Department{
		{DeptNo: "d001", DeptName: "Marketing"},
		{DeptNo: "d002", DeptName: "Finance"},
		{DeptNo: "d003", DeptName: "Human Resources"},
		{DeptNo: "d004", DeptName: "Engineering"},
		{DeptNo: "d005", DeptName: "Sales"},
	}
}

func generateDeptManagers(employees []Employee) []DeptManager {
	var managers []DeptManager
	for _, emp := range employees {
		lastChar := emp.EmpNo[len(emp.EmpNo)-1]
		deptNum := int(lastChar-'0')%5 + 1
		manager := DeptManager{
			DeptNo:   fmt.Sprintf("d00%d", deptNum),
			EmpNo:    emp.EmpNo,
			FromDate: gofakeit.DateRange(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Now()),
			ToDate:   gofakeit.DateRange(time.Now(), time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)),
		}
		managers = append(managers, manager)
	}
	return managers
}

func generateDeptEmps(employees []Employee) []DeptEmp {
	var deptEmps []DeptEmp
	for _, emp := range employees {
		lastChar := emp.EmpNo[len(emp.EmpNo)-1]
		deptNum := int(lastChar-'0')%5 + 1
		deptEmp := DeptEmp{
			EmpNo:    emp.EmpNo,
			DeptNo:   fmt.Sprintf("d00%d", deptNum),
			FromDate: gofakeit.DateRange(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Now()),
			ToDate:   gofakeit.DateRange(time.Now(), time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)),
		}
		deptEmps = append(deptEmps, deptEmp)
	}
	return deptEmps
}

func generateTitles(employees []Employee) []Title {
	var titles []Title
	titleNames := []string{"Manager", "Software Engineer", "Sales Associate", "HR Specialist"}

	for _, emp := range employees {
		lastChar := emp.EmpNo[len(emp.EmpNo)-1]
		titleIndex := int(lastChar-'0') % 4
		title := Title{
			EmpNo:    emp.EmpNo,
			Title:    titleNames[titleIndex],
			FromDate: gofakeit.DateRange(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Now()),
			ToDate:   gofakeit.DateRange(time.Now(), time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)),
		}
		titles = append(titles, title)
	}
	return titles
}

func generateSalaries(employees []Employee) []Salary {
	var salaries []Salary
	for _, emp := range employees {
		salary := Salary{
			EmpNo:    emp.EmpNo,
			Salary:   gofakeit.Number(50000, 150000), // Random salary between 50,000 and 150,000
			FromDate: gofakeit.DateRange(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Now()),
			ToDate:   gofakeit.DateRange(time.Now(), time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)),
		}
		salaries = append(salaries, salary)
	}
	return salaries
}

func main() {
	host := "localhost"
	port := "5430"
	user := "primary_user"
	password := "primary_password"
	dbname := "employees"
	schema := "employees"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s TimeZone=Asia/Bangkok",
		host, port, user, password, dbname, schema)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&Employee{}, &Department{}, &DeptManager{}, &DeptEmp{}, &Title{}, &Salary{})
	if err != nil {
		log.Fatalf("Failed to auto migrate schema: %v", err)
	}

	// Generate departments and insert into departments table
	departments := generateDepartments()
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "dept_no"}},
		DoUpdates: clause.AssignmentColumns([]string{"dept_name"}),
	}).Create(&departments).Error; err != nil {
		log.Fatalf("Failed to insert departments: %v", err)
	}
	fmt.Println("Departments inserted successfully")

	// Generate random employees and insert into employees table
	employees := generateRandomEmployees(1000000)
	insertInBatches(db, employees, 1000, "employees")
	fmt.Println("Employees inserted successfully")

	// Generate and insert dept managers
	deptManagers := generateDeptManagers(employees)
	insertInBatches(db, deptManagers, 1000, "department managers")

	// Generate and insert dept emps
	deptEmps := generateDeptEmps(employees)
	insertInBatches(db, deptEmps, 1000, "department employees")

	// Generate and insert titles
	titles := generateTitles(employees)
	insertInBatches(db, titles, 1000, "titles")

	// Generate and insert salaries
	salaries := generateSalaries(employees)
	insertInBatches(db, salaries, 1000, "salaries")
}

func insertInBatches(db *gorm.DB, data interface{}, batchSize int, name string) {
	val := reflect.ValueOf(data)
	for i := 0; i < val.Len(); i += batchSize {
		end := i + batchSize
		if end > val.Len() {
			end = val.Len()
		}
		if err := db.Create(val.Slice(i, end).Interface()).Error; err != nil {
			log.Fatalf("Failed to insert %s batch %d-%d: %v", name, i, end, err)
		}
		fmt.Printf("Inserted %s batch %d-%d\n", name, i, end)
	}
	fmt.Printf("All %s inserted successfully\n", name)
}
