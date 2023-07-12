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
	} else if strings.HasPrefix(message.Text, "/changerequest") { //todo in output
		command = "changerequest"
	}

	switch command {
	case "seetodaytasks":
		employeeUsername := message.From.UserName
		handleTodaysTasks(bot, message, dbConnection, updates, employeeUsername) //todo do it
	case "seealltasks":
		employeeUsername := message.From.UserName
		handleAllTasksInfo(bot, message, dbConnection, updates, employeeUsername)
	case "completetask":
		handleCompleteTask(bot, message, dbConnection, updates)
	case "forcemajeure":
		handleForceMajeureOrChangeRequest(bot, message, dbConnection, updates, "forcemajeure")
	case "changerequest":
		handleForceMajeureOrChangeRequest(bot, message, dbConnection, updates, "changerequest")
	case "recomendation":
		handleRecommendation(bot, message, dbConnection, updates)
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
		taskInfoMessage += fmt.Sprintf("Change Request: /changerequest%d\n", task.ID)      // Add button for reporting
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

func handleCompleteTask(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, updates tgbotapi.UpdatesChannel) {
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

	msg := tgbotapi.NewMessage(message.Chat.ID, "Send video or photo for task validation")
	bot.Send(msg)

	file_id, err := collectFile(bot, message.Chat.ID, updates)
	if err != nil {
		log.Println("Error collecting user input:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Wrong data format")
		bot.Send(msg)
		return
	}

	// Save the file ID in the database
	err = SaveFileID(dbConnection, taskID, file_id)
	if err != nil {
		log.Println("Error saving file ID:", err)
		response := "Failed to save the photo or video. Please try again later."
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
	msg = tgbotapi.NewMessage(message.Chat.ID, response)
	bot.Send(msg)
}

func MarkTaskAsCompleted(dbConnection db.Database, taskID int) error {
	// Update the task's status in the database
	err := dbConnection.Execute("UPDATE tasks SET isDone = true, date_requested_to_validate = NOW() WHERE id = ?", taskID)
	if err != nil {
		return err
	}

	return nil
}

func SaveFileID(dbConnection db.Database, taskID int, fileID string) error {
	query := `
		UPDATE tasks
		SET file_id = ?
		WHERE id = ?;
	`

	// Execute the update query
	err := dbConnection.Execute(query, fileID, taskID)
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
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.KeyboardButton{Text: "/recomendation"},
		),
	)

	msg.ReplyMarkup = replyMarkup
	bot.Send(msg)
}

func handleForceMajeureOrChangeRequest(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, updates tgbotapi.UpdatesChannel, issueName string) {

	var command, respond string
	switch issueName {
	case "forcemajeure":
		command = strings.TrimPrefix(message.Text, "/forcemajeure")
		respond = "Please provide the details of the force majeure situation:"
	case "changerequest":
		command = strings.TrimPrefix(message.Text, "/changerequest")
		respond = "Please provide the details of the change request:"
	}

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
	msg := tgbotapi.NewMessage(message.Chat.ID, respond)
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

	ok, err := checkForceMajorOrChangeRequest(dbConnection, taskID, issueName)
	if err != nil {
		log.Println("Error checking checkForceMajor:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "We had some issues, please try again later")
		bot.Send(msg)
		return
	}
	// Save the force majeure report to the database
	lastInsertID, err := SaveForceMajeureReportOrChangeRequest(dbConnection, taskID, taskInfo.ResidentialComplex, taskInfo.ElevatorName, taskInfo.EmployeePhoneNumber, forceMajeureDetails, message.From.ID, issueName)
	if err != nil {
		log.Println("Error saving force majeure report OrChangeRequest:", err)
		response := "Failed to save force majeure report OrChangeRequest. Please try again later."
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
		var response string
		switch issueName {
		case "forcemajeure":
			sendForceMajeureNotifications(bot, adminChatID, fmt.Sprintf(forceMajeureNotification+"\n/completeforcemajor%d", lastInsertID))
			response = "Force majeure report submitted successfully. The admin will be notified."
		case "changerequest":
			//todo fo  /completechangerequest
			sendForceMajeureNotifications(bot, adminChatID, fmt.Sprintf(forceMajeureNotification+"\n/completechangerequest%d", lastInsertID))
			response = "Change request report submitted successfully. The admin will be notified."
		}

		// Notify the employee about the successful report submission
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

		switch issueName {
		case "forcemajeure":
			sendForceMajeureNotifications(bot, HrChatID, fmt.Sprintf(forceMajeureNotification+"\n/completeforcemajor%d", lastInsertID))
		case "changerequest":
			sendForceMajeureNotifications(bot, HrChatID, fmt.Sprintf(forceMajeureNotification+"\n/completechangerequest%d", lastInsertID))
		}

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

func SaveForceMajeureReportOrChangeRequest(dbConnection db.Database, taskID int, nameResident, nameLift, employeePhoneNumber, details string, identifier int, issueName string) (int, error) {
	var query string
	switch issueName {
	case "forcemajeure":
		query = `
		INSERT INTO force_majeure (task_id, residential_complex, elevator_name, employee_phone_number, description, employee_identifier)
		VALUES (?, ?, ?, ?, ?, ?);
	`
	case "changerequest":
		query = `
		INSERT INTO change_requests (task_id, residential_complex, elevator_name, employee_phone_number, description, employee_identifier)
		VALUES (?, ?, ?, ?, ?, ?);
	`
	}

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

func checkForceMajorOrChangeRequest(dbConnection db.Database, taskID int, issueName string) (bool, error) {
	var query string
	switch issueName {
	case "forcemajeure":
		query = "SELECT COUNT(*) FROM force_majeure WHERE task_id = ?"
	case "changerequest":
		query = "SELECT COUNT(*) FROM change_requests WHERE task_id = ?"

	}

	row := dbConnection.QueryRow(query, taskID)
	if row.Err() != nil {
		log.Println("Error executing query in checkForceMajorOrChangeRequest:", row.Err())
		return false, row.Err()
	}

	var count int

	if err := row.Scan(&count); err != nil {
		log.Println("Error scanning row  in checkForceMajorOrChangeRequest:", err)
		return false, err
	}
	return count >= 2, nil
} //"your request  will pass to hr_manager's"

func handleGetFile(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, updates tgbotapi.UpdatesChannel) {
	fmt.Println("inside handleGetFile: ", message.Text)
	command := strings.TrimPrefix(message.Text, "/getfile")
	taskID, err := strconv.Atoi(command)
	fmt.Println("TaskID: ", taskID)

	if err != nil {
		log.Println("Invalid task ID:", err)
		response := "Invalid task ID. Please try again."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	// Retrieve the file ID from the database
	var fileID string
	query := `SELECT file_id FROM tasks WHERE id = ?`
	row := dbConnection.QueryRow(query, taskID)
	err = row.Scan(&fileID)
	if err != nil {
		log.Println("Error retrieving file ID:", err)
		response := "Failed to retrieve the file. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	// Retrieve the file information using Telegram Bot API's getFile method
	file, err := bot.GetFile(tgbotapi.FileConfig{
		FileID: fileID,
	})
	if err != nil {
		log.Println("Error retrieving file:", err)
		response := "Failed to retrieve the file. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	// Get the file URL and send it back to the user
	fileURL := file.Link(bot.Token)
	response := fmt.Sprintf("Here's the file you requested: %s", fileURL)
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	bot.Send(msg)
}

func handleRecommendation(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, updates tgbotapi.UpdatesChannel) {
	chatID := message.Chat.ID
	phoneNumber := message.From.UserName // Assuming the phone number is retrieved from the user's username field

	// Prompt the user to enter the recommendation description
	msg := tgbotapi.NewMessage(chatID, "Please enter your recommendation:")
	bot.Send(msg)

	// Wait for the user's input
	description, err := collectUserInput(bot, chatID, updates, "description", dbConnection)
	if err != nil {
		log.Println("Error collecting recommendation description:", err)
		response := "Failed to collect recommendation description. Please try again later."
		msg := tgbotapi.NewMessage(chatID, response)
		bot.Send(msg)
		return
	}

	// Save the recommendation to the database
	err = SaveRecommendation(dbConnection, phoneNumber, description)
	if err != nil {
		log.Println("Error saving recommendation:", err)
		response := "Failed to save the recommendation. Please try again later."
		msg := tgbotapi.NewMessage(chatID, response)
		bot.Send(msg)
		return
	}

	response := "Recommendation saved successfully. Thank you for your feedback!"
	msg = tgbotapi.NewMessage(chatID, response)
	bot.Send(msg)
}

func SaveRecommendation(dbConnection db.Database, phoneNumber, description string) error {
	query := "INSERT INTO recommendations (phone_number, description) VALUES (?, ?)"

	err := dbConnection.Execute(query, phoneNumber, description)
	if err != nil {
		return err
	}

	return nil
}
