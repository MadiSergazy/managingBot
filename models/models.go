package models

import "errors"

var (
	ErrNoOneExec      = errors.New("no one exec")
	ErrValidateTask   = errors.New("task cannot be validated due to incomplete requirements")
	ErrGettingAdminID = errors.New("Error getting AdminChatID")
	ErrGettingHRID    = errors.New("Error getting HRManagerChatID")
)

type Worker struct {
	ID          int
	PhoneNumber string
	Name        string
}

type Lift struct {
	ID       int
	Name     string
	TaskID   int
	WorkerID int
}

type Project struct {
	ID           int
	NameResident string
	// NameLift     string
	WorkerID int
	LiftID   int
}

type TaskDates struct {
	StartDate string
	EndDate   string
}

type Task struct {
	ID            int
	Name_resident string
	Name_lift     string
	TaskName      string
	StartDate     string
	EndDate       string
	IsDone        bool
}

// TaskInfo represents the information of a task
type TaskInfo struct {
	TaskID              int
	ElevatorName        string
	ResidentialComplex  string
	EmployeePhoneNumber string
	TaskName            string
	StartDate           string
	EndDate             string
	IsDone              bool
	Is_validate         bool
}
