package main

import (
	"fmt"
	"time"
)

type TaskDates struct {
	StartDate string
	EndDate   string
}
type TaskInfo struct {
	ElevatorName        string
	ResidentialComplex  string
	EmployeePhoneNumber string
	TaskName            string
	StartDate           string
	EndDate             string
	IsDone              bool
}

func main() {
	taskDatesMap := make(map[int]TaskDates)

	taskInfoMap := map[int]TaskInfo{
		1: {ElevatorName: "Elevator 1", ResidentialComplex: "Residential Complex 1", EmployeePhoneNumber: "1234567890", TaskName: "Task 1", StartDate: "01/07/2023", EndDate: "01/07/2023", IsDone: true},
		2: {ElevatorName: "Elevator 2", ResidentialComplex: "Residential Complex 2", EmployeePhoneNumber: "2345678901", TaskName: "Task 2", StartDate: "01/07/2023", EndDate: "01/07/2023", IsDone: false},
		3: {ElevatorName: "Elevator 3", ResidentialComplex: "Residential Complex 3", EmployeePhoneNumber: "3456789012", TaskName: "Task 3", StartDate: "01/07/2023", EndDate: "01/07/2023", IsDone: true},
		4: {ElevatorName: "Elevator 2", ResidentialComplex: "Residential Complex 2", EmployeePhoneNumber: "2345678901", TaskName: "Task 2", StartDate: "1/7/2023", EndDate: "01/07/2023", IsDone: false},
		5: {ElevatorName: "Elevator 3", ResidentialComplex: "Residential Complex 3", EmployeePhoneNumber: "3456789012", TaskName: "Task 3", StartDate: "1/7/2023", EndDate: "01/07/2023", IsDone: true},
		// Add more task information as needed
	}

	for taskID, taskInfo := range taskInfoMap {
		taskDatesMap[taskID] = TaskDates{
			StartDate: taskInfo.StartDate,
			EndDate:   taskInfo.EndDate,
		}
	}

	for taskID, dates := range taskDatesMap {
		// taskName := taskIDToName[taskID]
		startDate, _ := time.Parse("02/01/2006", dates.StartDate)
		endDate, _ := time.Parse("02/01/2006", dates.EndDate)
		formattedStartDate := startDate.Format("2006-01-02")
		formattedEndDate := endDate.Format("2006-01-02")
		fmt.Println("taskID: ", taskID)
		fmt.Println("Start Date: ", formattedStartDate)
		fmt.Println("Finish date: ", formattedEndDate)

		// _, err := db.connection.Exec(query, taskName, startDate, endDate, employeeID)
		// if err != nil {
		// 	fmt.Println("Start Date: ", startDate)
		// 	fmt.Println("Finish date: ", endDate)
		// 	fmt.Println("taskIDToName : ", taskIDToName)
		// 	return err
		// }
	}
}
