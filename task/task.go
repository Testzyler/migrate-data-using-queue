package task

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// A list of task types.
const (
	TypeMigrateEmployee = "export:employee"
	TypeRemoveEmployee  = "purge:employee"
)

type EmployeeProcessor struct {
	// fields for struct
	DB *gorm.DB
}

type EmployeeTaskPayload struct {
	EmpID string
}

func (p *EmployeeProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload EmployeeTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	db := p.DB

	db = db.WithContext(ctx)

	log.Println(" [*] Processing task:", t.Type, "with payload:", payload)

	// Start a transaction
	tx := db.Begin()
	err := tx.Raw("DELETE FROM employees.salaries WHERE emp_no = ?", &payload.EmpID).Error
	if err != nil {
		tx.Rollback()
		log.Println(" [*] Failed to delete from salaries:", err)
		return fmt.Errorf("failed to delete from salaries: %v", err)
	}

	// Delete from titles table
	err = tx.Exec("DELETE FROM employees.titles WHERE emp_no = $1", payload.EmpID).Error
	if err != nil {
		return fmt.Errorf("failed to delete from titles: %v", err)
	}

	// Delete from dept_emp table
	err = tx.Exec("DELETE FROM employees.dept_emps WHERE emp_no = $1", payload.EmpID).Error
	if err != nil {
		return fmt.Errorf("failed to delete from dept_emp: %v", err)
	}

	// Delete from dept_manager table
	err = tx.Exec("DELETE FROM employees.dept_managers WHERE emp_no = $1", payload.EmpID).Error
	if err != nil {
		return fmt.Errorf("failed to delete from dept_manager: %v", err)
	}

	// Finally, delete from employees table
	err = tx.Exec("DELETE FROM employees.employees WHERE emp_no = $1", payload.EmpID).Error
	if err != nil {
		return fmt.Errorf("failed to delete from employees: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf(" [*] Remove Employee %d", payload.EmpID)
	return nil
}

type EmployeeMigrator struct {
	Ctx *context.Context
	DB  *gorm.DB
}

func NewEmployeeProcessor(db *gorm.DB) *EmployeeProcessor {
	return &EmployeeProcessor{
		DB: db,
	}
}

func (migrator *EmployeeMigrator) NewMigrateEmployeeTasks(EmpID []string) ([]*asynq.Task, error) {
	var tasks []*asynq.Task
	// ctx := *migrator.ctx

	for _, empID := range EmpID {
		payload := EmployeeTaskPayload{
			EmpID: empID,
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		task := asynq.NewTask(TypeMigrateEmployee, payloadBytes)
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (migrator *EmployeeMigrator) HandleMigrateEmployeeTask(ctx context.Context, t *asynq.Task) error {
	var p EmployeeTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	log.Printf(" [*] Migrate Employee %d", p.EmpID)
	return nil

}

func (migrator *EmployeeMigrator) NewRemoveEmployeeTasks(EmpID []string) ([]*asynq.Task, error) {
	var tasks []*asynq.Task
	// ctx := *migrator.ctx

	for _, empID := range EmpID {
		payload := EmployeeTaskPayload{
			EmpID: empID,
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		task := asynq.NewTask(TypeRemoveEmployee, payloadBytes)
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (migrator *EmployeeMigrator) HandleRemoveEmployeeTask(ctx context.Context, t *asynq.Task) error {
	var p EmployeeTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	db := migrator.DB

	db = db.WithContext(ctx)

	// Start a transaction
	tx := db.Begin()
	// Delete from salaries table
	err := tx.Exec("DELETE FROM employees.salaries WHERE emp_no = $1", p.EmpID).Error
	if err != nil {
		return fmt.Errorf("failed to delete from salaries: %v", err)
	}

	// Delete from titles table
	err = tx.Exec("DELETE FROM employees.titles WHERE emp_no = $1", p.EmpID).Error
	if err != nil {
		return fmt.Errorf("failed to delete from titles: %v", err)
	}

	// Delete from dept_emp table
	err = tx.Exec("DELETE FROM employees.dept_emps WHERE emp_no = $1", p.EmpID).Error
	if err != nil {
		return fmt.Errorf("failed to delete from dept_emps: %v", err)
	}

	// Delete from dept_manager table
	err = tx.Exec("DELETE FROM employees.dept_managers WHERE emp_no = $1", p.EmpID).Error
	if err != nil {
		return fmt.Errorf("failed to delete from dept_managers: %v", err)
	}

	// Finally, delete from employees table
	err = tx.Exec("DELETE FROM employees.employees WHERE emp_no = $1", p.EmpID).Error
	if err != nil {
		return fmt.Errorf("failed to delete from employees: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("[*] Remove Employee %v", p.EmpID)

	return nil
}

func (migrator *EmployeeMigrator) NewRemoveAllEmployeeTasks() ([]*asynq.Task, error) {
	var tasks []*asynq.Task
	db := migrator.DB
	// Get all employees
	var empIDs []string
	err := db.Raw("SELECT emp_no FROM employees.employees").Scan(&empIDs).Error
	if err != nil {
		return nil, err
	}

	for _, empID := range empIDs {
		payload := EmployeeTaskPayload{
			EmpID: empID,
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		task := asynq.NewTask(TypeRemoveEmployee, payloadBytes)
		tasks = append(tasks, task)
	}
	return tasks, nil
}
