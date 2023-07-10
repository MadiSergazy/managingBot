package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"

	"madi_telegram_bot/db"
	"madi_telegram_bot/models"
)

func HandleWorkerCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updateBuffer []tgbotapi.Update, dbConnection db.Database, updates tgbotapi.UpdatesChannel) {
	response := "Unknown command. Please use valid admin commands."
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	message.CommandArguments()
	command := strings.ToLower(message.Command())
	if strings.HasPrefix(command, "completetask") {
		command = "completetask"
		// commandArg :=
	} else if strings.HasPrefix(message.Text, "/forcemajeure") {
		command = "forcemajeure"
	}
	switch command {
	case "seetodaytasks":
		employeeUsername := message.From.UserName
		handleTodaysTasks(bot, message, dbConnection, updates, employeeUsername) //todo do it
	case "seealltasks":
		employeeUsername := message.From.UserName
		handleAllTasksInfo(bot, message, dbConnection, updates, employeeUsername)
	case "completetask":
		handleCompleteTask(bot, message, dbConnection)
	case "forcemajeure":
		handleForceMajeure(bot, message, dbConnection, updates)
	case "returnToMenu":
		returnToMenu(bot, msg)
	default:
		returnToMenu(bot, msg)

	}

}

func handleAllTasksInfo(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, updates tgbotapi.UpdatesChannel, employeeUsername string) {
	employeePhoneNumber, err := getEmployeePhoneNumber(dbConnection, employeeUsername)
	if err != nil {
		log.Println("Error retrieving employee phone number:", err)
		response := "Failed to retrieve employee phone number."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}
	ShowTasks(bot, message, dbConnection, employeePhoneNumber)
}

func getEmployeePhoneNumber(dbConnection db.Database, username string) (string, error) {
	query := "SELECT phone_number FROM workers WHERE name = ?"
	row := dbConnection.QueryRow(query, username)

	var phoneNumber string
	err := row.Scan(&phoneNumber)
	if err != nil {
		return "", err
	}

	return phoneNumber, nil
}

func handleTodaysTasks(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, updates tgbotapi.UpdatesChannel, employeeUsername string) {
	// Get the employee's phone number based on the provided username
	employeePhoneNumber, err := getEmployeePhoneNumber(dbConnection, employeeUsername)
	if err != nil {
		log.Println("Error retrieving employee phone number:", err)
		response := "Failed to retrieve employee information. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	// Get today's tasks for the employee
	todayTasks, err := GetTasksForToday(dbConnection, employeePhoneNumber)
	if err != nil {
		log.Println("Error retrieving today's tasks:", err)
		response := "Failed to retrieve today's tasks. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	// Display the tasks to the user
	if len(todayTasks) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "No tasks for today.")
		bot.Send(msg)
		return
	}

	taskInfoMessage := "Today's Tasks:\n\n"
	for _, task := range todayTasks {
		doneSymbol := "✅" // Green checkmark
		if !task.IsDone {
			doneSymbol = "❌" // Red cross
		}
		taskInfoMessage += "Resident: " + task.Name_resident + "\nLift: " + task.Name_lift + "\nTask: " + task.TaskName + "\nStart Date: " + task.StartDate + "\nEnd Date: " + task.EndDate + "\nIs Done: " + doneSymbol + "\n"

		// Add inline keyboard buttons to mark the task as completed
		if !task.IsDone {
			taskInfoMessage += fmt.Sprintf("Mark as Completed: /completetask%d\n", task.ID)
		}

		taskInfoMessage += fmt.Sprintf("Report Force Majeure: /forcemajeure%d\n", task.ID) // Add button for reporting force majeure
		taskInfoMessage += "\n"

	}

	msg := tgbotapi.NewMessage(message.Chat.ID, taskInfoMessage)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func GetTasksForToday(dbConnection db.Database, employeePhoneNumber string) ([]models.Task, error) {
	// Get the current date
	currentDate := time.Now().Format("2006-01-02")
	// Retrieve the worker_id for the given employeePhoneNumber
	// var workerID int
	// err := dbConnection.QueryRow("SELECT id FROM workers WHERE phone_number = ?", employeePhoneNumber).Scan(&workerID)
	// if err != nil {
	// 	return nil, err
	// }
	//todo check me
	rows, err := dbConnection.Query(
		`SELECT
		t.id,
		p.name_resident,
		l.name_lift,
		t.nameOfTask,
		t.dateStart,
		t.dateEnd,
		t.isDone
	  FROM projects p
	  JOIN lifts l ON p.lift_id = l.id
	  JOIN workers w ON p.worker_id = w.id
	  JOIN tasks t ON t.lift_id = l.id
	  WHERE w.phone_number = ?  AND ? BETWEEN t.dateStart AND t.dateEnd
	  ORDER BY p.id, t.id;`, employeePhoneNumber, currentDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and populate the task information
	tasks := []models.Task{}
	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(&task.ID, &task.Name_resident, &task.Name_lift, &task.TaskName, &task.StartDate, &task.EndDate, &task.IsDone)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func handleCompleteTask(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database) {
	fmt.Println("Inside handleCompleteTask")
	// Extract the task ID from the command
	command := strings.TrimPrefix(message.Text, "/completetask")
	taskID, err := strconv.Atoi(command)
	fmt.Println("TAskID: ", taskID)

	if err != nil {
		log.Println("Invalid task ID:", err)
		response := "Invalid task ID. Please try again."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	// Update the task as completed in the database
	err = MarkTaskAsCompleted(dbConnection, taskID)
	if err != nil {
		log.Println("Error marking task as completed:", err)
		response := "Failed to mark task as completed. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	response := "Task marked as completed successfully. return to menu: /returnToMenu"
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	bot.Send(msg)
}

func MarkTaskAsCompleted(dbConnection db.Database, taskID int) error {
	// Update the task's status in the database
	err := dbConnection.Execute("UPDATE tasks SET isDone = true WHERE id = ?", taskID)
	if err != nil {
		return err
	}

	return nil
}

func returnToMenu(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	replyMarkup := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{Text: "/seetodaytasks"},
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{Text: "/seealltasks"},
		),
	)

	msg.ReplyMarkup = replyMarkup
	bot.Send(msg)
}

func handleForceMajeure(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, updates tgbotapi.UpdatesChannel) {
	command := strings.TrimPrefix(message.Text, "/forcemajeure")
	taskID, err := strconv.Atoi(command)
	fmt.Println("TAskID: ", taskID)

	// Get the task information based on the task ID
	taskInfo, err := GetTaskByID(dbConnection, taskID)
	if err != nil {
		log.Println("Error retrieving task information:", err)
		response := "Failed to retrieve task information. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	// Prompt the employee to enter the force majeure details
	msg := tgbotapi.NewMessage(message.Chat.ID, "Please provide the details of the force majeure situation:")
	bot.Send(msg)

	// Collect the force majeure details from the user
	forceMajeureDetails, err := collectUserInput(bot, message.Chat.ID, updates, "complexName", dbConnection)
	if err != nil {
		log.Println("Error collecting user input:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Wrong data format")
		bot.Send(msg)
		return
	}
	fmt.Println("forceMajeureDetails:", forceMajeureDetails)

	forceMajeureNotification := fmt.Sprintf("Force Majeure Report:\n\nResident: %s\nLift: %s\nEmployee Phone: %s\nStart Date: %s\n\nDetails:\n%s",
		taskInfo.ResidentialComplex, taskInfo.ElevatorName, taskInfo.EmployeePhoneNumber, taskInfo.StartDate, forceMajeureDetails)

	ok, err := checkForceMajor(dbConnection, taskID)
	if err != nil {
		log.Println("Error checking checkForceMajor:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "We had some issues, please try again later")
		bot.Send(msg)
		return
	}
	// Save the force majeure report to the database
	lastInsertID, err := SaveForceMajeureReport(dbConnection, taskID, taskInfo.ResidentialComplex, taskInfo.ElevatorName, taskInfo.EmployeePhoneNumber, forceMajeureDetails, message.From.ID)
	if err != nil {
		log.Println("Error saving force majeure report:", err)
		response := "Failed to save force majeure report. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}
	if !ok {

		// Notify the admin about the force majeure report
		// adminUsername := getAdminUsername(dbConnection) // Replace this with your logic to retrieve the admin username
		adminChatID, err := getChatIDs(dbConnection, "admins") // Replace this with your logic to retrieve the admin chat ID
		if err != nil {
			log.Println("Error getting admins chatID:", err)
			response := "Failed to save force majeure report. Please try again later."
			msg := tgbotapi.NewMessage(message.Chat.ID, response)
			bot.Send(msg)
			return
		}

		// send notification for admins
		sendForceMajeureNotifications(bot, adminChatID, fmt.Sprintf(forceMajeureNotification+"\n/completeforcemajor%d", lastInsertID))

		// Notify the employee about the successful report submission
		response := "Force majeure report submitted successfully. The admin will be notified."
		msg = tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)

	} else {
		HrChatID, err := getChatIDs(dbConnection, "hr_manager") // Replace this with your logic to retrieve the admin chat ID
		if err != nil {
			log.Println("Error getting HrChatID:", err)
			response := "Failed to save force majeure report. Please try again later."
			msg := tgbotapi.NewMessage(message.Chat.ID, response)
			bot.Send(msg)
			return
		}
		response := "your request  will pass to hr_manager's"
		msg = tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		sendForceMajeureNotifications(bot, HrChatID, fmt.Sprintf(forceMajeureNotification+"\n/completeforcemajor%d", lastInsertID))

	}

}

func GetTaskByID(dbConnection db.Database, taskID int) (models.TaskInfo, error) {
	query := `
		SELECT
			p.name_resident,
			l.name_lift,
			w.phone_number,
			t.dateStart
		FROM projects p
		JOIN lifts l ON p.lift_id = l.id
		JOIN workers w ON p.worker_id = w.id
		JOIN tasks t ON t.lift_id = l.id
		WHERE t.id = ?;
	`

	row := dbConnection.QueryRow(query, taskID)

	taskInfo := models.TaskInfo{}
	err := row.Scan(
		&taskInfo.ResidentialComplex,
		&taskInfo.ElevatorName,
		&taskInfo.EmployeePhoneNumber,
		&taskInfo.StartDate,
	)
	if err != nil {
		return taskInfo, err
	}

	return taskInfo, nil
}

func SaveForceMajeureReport(dbConnection db.Database, taskID int, nameResident, nameLift, employeePhoneNumber, details string, identifier int) (int, error) {
	query := `
		INSERT INTO force_majeure (task_id, residential_complex, elevator_name, employee_phone_number, description, employee_identifier)
		VALUES (?, ?, ?, ?, ?, ?);
	`

	lastInsertID, err := dbConnection.ExecuteWithLastInsertID(query, taskID, nameResident, nameLift, employeePhoneNumber, details, identifier)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func getChatIDs(dbConnection db.Database, tableName string) ([]int64, error) {
	query := `select identifier from ` + tableName + `;`

	rows, err := dbConnection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chatIDs := make([]int64, 0)
	for rows.Next() {
		var chatID int64
		err := rows.Scan(&chatID)
		if err != nil {
			return nil, err
		}
		chatIDs = append(chatIDs, chatID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chatIDs, nil
}

func sendForceMajeureNotifications(bot *tgbotapi.BotAPI, chatIDs []int64, forceMajeureNotification string) {
	var wg sync.WaitGroup
	wg.Add(len(chatIDs))

	for _, chatID := range chatIDs {
		go func(id int64) {
			defer wg.Done()

			msg := tgbotapi.NewMessage(id, forceMajeureNotification)
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("Failed to send notification to chat ID:", id)
			}
		}(chatID)
	}

	wg.Wait()
}

func checkForceMajor(dbConnection db.Database, taskID int) (bool, error) {
	query := "SELECT COUNT(*) FROM force_majeure WHERE task_id = ?"
	row := dbConnection.QueryRow(query, taskID)
	if row.Err() != nil {
		log.Println("Error executing query in checkForceMajor:", row.Err())
		return false, row.Err()
	}

	var count int

	if err := row.Scan(&count); err != nil {
		log.Println("Error scanning row:", err)
		return false, err
	}
	return count >= 2, nil
} //"your request  will pass to hr_manager's"
