package handlers

import (
	"errors"
	"fmt"
	"strconv"
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
			validateSymbol := "✅" // Green checkmark
			if !task.Is_validate {
				validateSymbol = "❌" // Red cross
			}

			taskID := strconv.Itoa(task.TaskID)
			taskInfoMessage += "\nElevator: " + task.ElevatorName + "\nResidential Complex:" + task.ResidentialComplex + "\nEmployee Phone Number:" + task.EmployeePhoneNumber + "\nTask: " + task.TaskName + "\nStart Date: " + task.StartDate + "\nEnd Date: " + task.EndDate + "\nIs Done: " + doneSymbol + "\nIs Validate: " + validateSymbol + "\n/getfile" + taskID + "\n/reject" + taskID + "\n/validate" + taskID + "\n\n"
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, taskInfoMessage)
		bot.Send(msg)
	}

}

func CheckUnfinished(dbConnection db.Database, bot *tgbotapi.BotAPI, nameTable string) error {
	var query string

	switch nameTable {
	case "force_majeure":
		query = `
		SELECT id, task_id, residential_complex, elevator_name, employee_phone_number, incident_time, description
		FROM force_majeure
		WHERE is_done = false AND incident_time <= NOW() - INTERVAL 2 DAY;
	`
	case "change_requests":
		query = `
		SELECT id, task_id, residential_complex, elevator_name, employee_phone_number, incident_time, description
		FROM change_requests
		WHERE is_done = false AND incident_time <= NOW() - INTERVAL 2 DAY;
	`
	}
	// Query the database for unfinished force majeure reports

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

func CheckOverdueTasks(dbConnection db.Database, bot *tgbotapi.BotAPI) error {
	query := `SELECT
	t.id,
	p.name_resident,
	l.name_lift,
	w.phone_number,
	t.nameOfTask,
	t.dateStart,
	t.dateEnd,
	t.isDone
  FROM projects p
  JOIN lifts l ON p.lift_id = l.id
  JOIN workers w ON p.worker_id = w.id
  JOIN tasks t ON t.lift_id = l.id
  WHERE t.dateEnd < (NOW() - INTERVAL 1 HOUR)  and t.isDone = false
  ORDER BY p.id, t.id;`

	rows, err := dbConnection.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Iterate over the rows and process the overdue tasks
	for rows.Next() {
		var (
			taskID              int
			nameResident        string
			nameLift            string
			employeePhoneNumber string
			taskName            string
			startDate           []uint8
			endDate             []uint8
			isDone              bool
		)

		err := rows.Scan(&taskID, &nameResident, &nameLift, &employeePhoneNumber, &taskName, &startDate, &endDate, &isDone)
		if err != nil {
			return err
		}

		// Convert the start and end date strings to time.Time objects
		startTime, _ := time.Parse("2006-01-02", string(startDate))
		endTime, _ := time.Parse("2006-01-02", string(endDate))

		overdueId, err := InsertOverdueTask(dbConnection, taskID, nameResident, nameLift, employeePhoneNumber, taskName, startTime, endTime)
		if err != nil {
			return err
		}

		// Perform the notification for the overdue task
		AdminChatID, err := getChatIDs(dbConnection, "admins") // Replace this with your logic to retrieve the admin chat ID
		if err != nil {
			log.Println("Error getting AdminChatID:", err)
			return errors.New("Error getting AdminChatID")
		}

		overdueTaskNotification := fmt.Sprintf("Overdue Task:\n\nResident: %s\nLift: %s\nEmployee Phone: %s\nTask: %s\nStart Date: %s\nEnd Date: %s",
			nameResident, nameLift, employeePhoneNumber, taskName, startTime.Format("02/01/2006"), endTime.Format("02/01/2006"))

		sendForceMajeureNotifications(bot, AdminChatID, fmt.Sprintf(overdueTaskNotification+"\n/CompleteOverdueTaskByAdmin%d", overdueId))

		// check if admin has solved the task within 24 hours
		checkForAdminsCompleteExpiredTasks(dbConnection, bot, overdueId, overdueTaskNotification)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func InsertOverdueTask(dbConnection db.Database, taskID int, nameResident, nameLift, phoneNumber, nameOfTask string, dateStart, dateEnd time.Time) (int, error) {
	query := `
		INSERT INTO overdue_task (task_id, name_resident, name_lift, phone_number, name_of_task, date_start, date_end, is_done_by_worker)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`

	overdueID, err := dbConnection.ExecuteWithLastInsertID(query, taskID, nameResident, nameLift, phoneNumber, nameOfTask, dateStart.Format("2006-01-02"), dateEnd.Format("2006-01-02"), false)
	if err != nil {
		return 0, err
	}

	return overdueID, nil
}

func IsTaskDoneByAdmin(dbConnection db.Database, taskID int) (bool, error) {
	query := "SELECT is_done_by_admin FROM overdue_task WHERE id = ?"
	row := dbConnection.QueryRow(query, taskID)

	var isDoneByAdmin bool
	err := row.Scan(&isDoneByAdmin)
	if err != nil {
		return false, err
	}

	return isDoneByAdmin, nil
}

// Start a goroutine to check if admin has solved the task within 24 hours
func checkForAdminsCompleteExpiredTasks(dbConnection db.Database, bot *tgbotapi.BotAPI, overdueId int, overdueTaskNotification string) {
	timer := delay(24 * time.Hour) // Check after 24 hours

	// Perform other tasks concurrently

	// Wait for the timer to expire or receive a value from the channel
	select {
	case <-timer:
		// Check if the task is still not marked as done by the admin
		isDoneByAdmin, err := IsTaskDoneByAdmin(dbConnection, overdueId)
		if err != nil {
			log.Println("Error checking task status by admin:", err)
			return
		}

		// If the task is still not done by the admin, notify the HR manager
		if !isDoneByAdmin {
			hrManagerChatID, err := getChatIDs(dbConnection, "hr_manager") // Replace this with your logic to retrieve the HR manager chat ID
			if err != nil {
				log.Println("Error getting HRManagerChatID:", err)
				return
			}
			//todo implement command handler for /CompleteOverdueTaskByHrManager
			sendForceMajeureNotifications(bot, hrManagerChatID, fmt.Sprintf(overdueTaskNotification+"\n/CompleteOverdueTaskByHrManager%d", overdueId))

		}
		break
	}
}

func delay(duration time.Duration) <-chan time.Time {
	ch := make(chan time.Time, 1)
	go func() {
		time.Sleep(duration)
		ch <- time.Now()
		close(ch)
	}()
	return ch
}

func CheckPendingValidationTasks(dbConnection db.Database, bot *tgbotapi.BotAPI) error {
	query := `
		SELECT t.id, p.name_resident, l.name_lift, w.phone_number, t.nameOfTask, t.dateStart, t.dateEnd, t.isDone, t.date_requested_to_validate
		FROM tasks t
		JOIN projects p ON t.lift_id = p.lift_id
		JOIN lifts l ON t.lift_id = l.id
		JOIN workers w ON p.worker_id = w.id
		WHERE t.isDone = true AND t.is_validate = false AND t.date_requested_to_validate <= NOW() - INTERVAL 48 HOUR;
	`

	rows, err := dbConnection.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Iterate over the rows and process the pending validation tasks
	for rows.Next() {
		var (
			id            int
			nameResident  string
			nameLift      string
			phone         string
			nameOfTask    string
			dateStart     []uint8
			dateEnd       []uint8
			isDone        bool
			dateRequested []uint8
		)

		err := rows.Scan(&id, &nameResident, &nameLift, &phone, &nameOfTask, &dateStart, &dateEnd, &isDone, &dateRequested)
		if err != nil {
			return err
		}

		dateStartTime, _ := time.Parse("2006-01-02", string(dateStart))
		dateEndTime, _ := time.Parse("2006-01-02", string(dateEnd))
		dateRequestedTime, _ := time.Parse("2006-01-02", string(dateRequested))

		hrManagerChatID, err := getChatIDs(dbConnection, "hr_manager")
		if err != nil {
			log.Println("Error getting HRManagerChatID:", err)
			return errors.New("Error getting HRManagerChatID")
		}

		notification := fmt.Sprintf("Task pending validation:\nTask ID: %d\nResident: %s\nLift: %s\nEmployee Phone: %s\nTask: %s\nStart Date: %s\nEnd Date: %s\nDate Requested: %s",
			id, nameResident, nameLift, phone, nameOfTask, dateStartTime.Format("02/01/2006"), dateEndTime.Format("02/01/2006"), dateRequestedTime.Format("02/01/2006"))
		sendForceMajeureNotifications(bot, hrManagerChatID, notification)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}
