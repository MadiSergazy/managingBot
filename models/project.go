package models

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
	NameLift     string
	WorkerID     int
	LiftID       int
}
