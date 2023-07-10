package handlers

import (
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"

	"madi_telegram_bot/db"
)

func ShowTasks(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, employeePhoneNumber string) {

	taskInfoList, err := db.ViewTasksForAdmin(dbConnection, employeePhoneNumber)
	if err != nil {
		log.Println("Error retrieving task information:", err)
		response := "Invalid employe phoneNumber. Please provide a valid phoneNumber."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	// Display the task information to the user
	if len(taskInfoList) == 0 {
		// No tasks found for the specified employee
		msg := tgbotapi.NewMessage(message.Chat.ID, "No tasks found for the specified employee.")
		bot.Send(msg)
	} else {
		// Tasks found, send the information as a formatted message

		// taskInfo := taskInfoList[0] // Get the first task info since the fields are the same for all tasks
		// taskInfoMessage := fmt.Sprintf("Elevator: %s\nResidential Complex: %s\nEmployee Phone Number: %s\n\n", taskInfo.ElevatorName, taskInfo.ResidentialComplex, taskInfo.EmployeePhoneNumber)
		var taskInfoMessage string
		for _, task := range taskInfoList {
			doneSymbol := "✅" // Green checkmark
			if !task.IsDone {
				doneSymbol = "❌" // Red cross
			}
			taskInfoMessage += "\nElevator: " + task.ElevatorName + "\nResidential Complex:" + task.ResidentialComplex + "\nEmployee Phone Number:" + task.EmployeePhoneNumber + "\nTask: " + task.TaskName + "\nStart Date: " + task.StartDate + "\nEnd Date: " + task.EndDate + "\nIs Done: " + doneSymbol + "\n\n"
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, taskInfoMessage)
		bot.Send(msg)
	}

}

func CheckUnfinishedForceMajeureReports(dbConnection db.Database, bot *tgbotapi.BotAPI) error {
	// Query the database for unfinished force majeure reports
	query := `
		SELECT id, task_id, residential_complex, elevator_name, employee_phone_number, incident_time, description
		FROM force_majeure
		WHERE is_done = false AND incident_time <= NOW() - INTERVAL 2 DAY;
	`

	rows, err := dbConnection.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Iterate over the rows and process the unfinished reports
	for rows.Next() {
		var (
			id                  int
			description         string
			taskID              int
			residentialComplex  string
			elevatorName        string
			employeePhoneNumber string
			incidentTime        []uint8
		)

		err := rows.Scan(&id, &taskID, &residentialComplex, &elevatorName, &employeePhoneNumber, &incidentTime, &description)
		if err != nil {
			return err
		}
		incidentTime1, _ := time.Parse("2006-01-02", string(incidentTime))
		// Perform the notification for the unfinished report
		HrChatID, err := getChatIDs(dbConnection, "hr_manager") // Replace this with your logic to retrieve the admin chat ID
		if err != nil {
			log.Println("Error getting HrChatID:", err)

			return errors.New("Error getting HrChatID")
		}
		forceMajeureNotification := fmt.Sprintf("Force Majeure Report time is expired:\n\nResident: %s\nLift: %s\nEmployee Phone: %s\nStart Date: %s\n\nDetails:\n%s",
			residentialComplex, elevatorName, employeePhoneNumber, incidentTime1.Format("02/01/2006"), description)

		sendForceMajeureNotifications(bot, HrChatID, fmt.Sprintf(forceMajeureNotification+"\n/completeforcemajor%d", id))

		// SendForceMajeureNotification(id, taskID, residentialComplex, elevatorName, employeePhoneNumber, incidentTime)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

// func SendForceMajeureNotification(id, taskID int, residentialComplex, elevatorName, employeePhoneNumber string, incidentTime time.Time) {
// 	// Implement the logic to send the notification here
// 	log.Printf("Sending force majeure notification for ID: %d\n", id)
// 	log.Printf("Task ID: %d\n", taskID)
// 	log.Printf("Residential Complex: %s\n", residentialComplex)
// 	log.Printf("Elevator Name: %s\n", elevatorName)
// 	log.Printf("Employee Phone Number: %s\n", employeePhoneNumber)
// 	log.Printf("Incident Time: %s\n", incidentTime.String())
// }
