package handlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"

	"madi_telegram_bot/db"
	"madi_telegram_bot/models"
)

// HandleAdminCommand handles admin-specific commands
func HandleAdminCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updateBuffer []tgbotapi.Update, dbConnection db.Database, updates tgbotapi.UpdatesChannel) {
	response := "Unknown command. Please use valid admin commands."
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	message.CommandArguments()
	command := strings.ToLower(message.Command())

	if strings.HasPrefix(command, "seeproject") {
		fmt.Println("Inside strings.HasPrefix")
		command = "seeproject"
		// commandArg :=
	} else if strings.HasPrefix(message.Text, "/completeforcemajor") {
		command = "completeforcemajor"
	} else if strings.HasPrefix(message.Text, "/completechangerequest") {
		command = "completechangerequest"
	} else if strings.HasPrefix(message.Text, "/CompleteOverdueTaskByAdmin") {
		command = "completeoverduetaskbyadmin"
	}

	switch command {
	case "createproject":
		//msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
		handleCreateProject(bot, message, updateBuffer, dbConnection, updates) //todo add some check fot time's and also check for the if employe phone number is exists
	case "projectlist":
		//msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
		handleProjectList(bot, message, dbConnection, updates) //todo to me
	case "projectinfo":
		//msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true} todo think about it
		handleProjectInfo(bot, message, dbConnection, updates)
	case "seeproject":
		handleProjectByID(bot, message, dbConnection, updates)
	case "completeforcemajor":
		handlerCompleteForceMajorOrChangereques(bot, message, dbConnection, updates, "completeforcemajor")
	case "completechangerequest":
		handlerCompleteForceMajorOrChangereques(bot, message, dbConnection, updates, "completechangerequest")
	case "completeoverduetaskbyadmin":
		handlerCompleteOverdueTaskByAdmin(bot, message, dbConnection, updates)
	default:
		replyMarkup := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.KeyboardButton{Text: "/createproject"},
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.KeyboardButton{Text: "/projectlist"},
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.KeyboardButton{Text: "/projectinfo"},
			),
		)

		msg.ReplyMarkup = replyMarkup
		bot.Send(msg)

	}

}

// Function to send the message with the task buttons
func sendTaskButtons(bot *tgbotapi.BotAPI, chatID int64, messageID int, buttons []tgbotapi.InlineKeyboardButton) {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, button := range buttons {
		row := []tgbotapi.InlineKeyboardButton{button}
		rows = append(rows, row)
	}
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, "Please select tasks:")
	msg.ReplyMarkup = inlineKeyboardMarkup
	msg.ReplyToMessageID = messageID
	bot.Send(msg)
}

func handleCreateProject(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updateBuffer []tgbotapi.Update, dbConnection db.Database, updates tgbotapi.UpdatesChannel) {
	log.Info("Inide of the handleCreateProject")
	// Create a channel to receive user messages
	//userInputChan := make(chan string)

	// Prompt the user to enter the name of the residential complex
	msg := tgbotapi.NewMessage(message.Chat.ID, "Please enter the name of the residential complex:")
	bot.Send(msg)

	// Start a goroutine to collect user input
	residentialComplex, err := collectUserInput(bot, message.Chat.ID, updates, "complexName", dbConnection)
	if err != nil {
		log.Println("Error collecting user input:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Wrong data format")
		bot.Send(msg)
		return
	}

	fmt.Println("Residential Complex:", residentialComplex)

	msg = tgbotapi.NewMessage(message.Chat.ID, "Please enter the name of the elevator:")
	bot.Send(msg)
	elevatorName, err := collectUserInput(bot, message.Chat.ID, updates, "elevatorName", dbConnection)
	if err != nil {
		log.Println("Error collecting user input:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Wrong data format")
		bot.Send(msg)
		return
	}
	fmt.Println("elevatorName: ", elevatorName)

	msg = tgbotapi.NewMessage(message.Chat.ID, "Please enter the phone number of the responsible employee(example: 77078566392):")
	bot.Send(msg)
	employeePhoneNumber, err := collectUserInput(bot, message.Chat.ID, updates, "employeePhoneNumber", dbConnection)
	if err != nil {
		log.Println("Error collecting user input:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Wrong data format")
		bot.Send(msg)
		return
	}
	fmt.Println("employeePhoneNumber: ", employeePhoneNumber)

	//todo: Insert the project into the database

	// Query the database to get the tasks from "task_of_lifts" table
	tasks, err := GetTasksFromLifts(dbConnection)
	if err != nil {
		log.Println("Error retrieving tasks from the database:", err)
		// Handle the error
		return
	}

	// Create a map to store the task IDs and their corresponding names
	taskIDToName := make(map[int]string)
	for _, task := range tasks {
		taskIDToName[task.ID] = task.TaskName
	}

	// Create a slice to store the task inline buttons
	// Create inline buttons for each task and add them to the slice
	taskButtons := make([]tgbotapi.InlineKeyboardButton, 0)
	for _, task := range tasks {
		taskButton := tgbotapi.NewInlineKeyboardButtonData(task.TaskName, strconv.Itoa(task.ID))
		taskButtons = append(taskButtons, taskButton)
	}

	// Create an inline keyboard markup with the task buttons
	var rows [][]tgbotapi.InlineKeyboardButton //todo rewtire it to another function
	for _, button := range taskButtons {
		row := []tgbotapi.InlineKeyboardButton{button}
		rows = append(rows, row)
	}
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// Create a map to store the selected tasks and their corresponding dates
	selectedTasks := make(map[int]models.TaskDates)
	// Send a message to the user with the task buttons
	msg = tgbotapi.NewMessage(message.Chat.ID, "Please select tasks:")
	msg.ReplyMarkup = inlineKeyboardMarkup
	bot.Send(msg)

	// Wait for the user to select tasks
	for len(taskButtons) > 0 {
		select {
		case update := <-updates:
			if update.CallbackQuery != nil {
				selectedTaskID, err := strconv.Atoi(update.CallbackQuery.Data)
				if err != nil {
					log.Println("Invalid task ID:", err)
					// Handle the error
					return
				}

				selectedTaskName := taskIDToName[selectedTaskID]

				// Remove the selected task from the list of buttons
				taskButtons = removeTaskButton(taskButtons, selectedTaskID)

				// Send the updated list of buttons to the user
				sendTaskButtons(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, taskButtons)

				// Create an inline keyboard markup with the updated task buttons
				rows := make([][]tgbotapi.InlineKeyboardButton, 0)
				for _, button := range taskButtons {
					row := []tgbotapi.InlineKeyboardButton{button}
					rows = append(rows, row)
				}
				inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup()
				for _, row := range rows {
					inlineKeyboardMarkup.InlineKeyboard = append(inlineKeyboardMarkup.InlineKeyboard, row)
				}

				// Send a message to the user to set the start date
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("You have selected the task: %s\nPlease enter the start date (dd/mm/yyyy):", selectedTaskName))
				bot.Send(msg)

				// Wait for the user to enter the start date
				startDate, err := collectUserInput(bot, update.CallbackQuery.Message.Chat.ID, updates, "startDate", dbConnection)
				if err != nil {
					log.Println("Error collecting user input:", err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong data format")
					bot.Send(msg)
					return
				}
				fmt.Println("startDate: ", startDate)

				// Send a message to the user to set the end date
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Please enter the end date (dd/mm/yyyy):")
				bot.Send(msg)

				// Wait for the user to enter the end date
				endDate, err := collectUserInput(bot, update.CallbackQuery.Message.Chat.ID, updates, "endDate", dbConnection)
				if err != nil {
					log.Println("Error collecting user input:", err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong data format")
					bot.Send(msg)
					return
				}

				fmt.Println("endDate: ", endDate)

				// Store the selected task and its dates in the map
				selectedTasks[selectedTaskID] = models.TaskDates{
					StartDate: startDate,
					EndDate:   endDate,
				}

				// Send a confirmation message to the user
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Task dates have been set.")
				bot.Send(msg)

				// Update the message with the updated task buttons
				editMsg := tgbotapi.EditMessageReplyMarkupConfig{
					BaseEdit: tgbotapi.BaseEdit{
						ChatID:      update.CallbackQuery.Message.Chat.ID,
						MessageID:   update.CallbackQuery.Message.MessageID,
						ReplyMarkup: &inlineKeyboardMarkup,
					},
				}
				bot.Send(editMsg)

			} else if update.Message != nil && update.Message.Text == "cancel" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Process canceled.")
				bot.Send(msg)
				return
			}
		}
	}

	// TODO: Insert the task, start date, and end date into the database

	fmt.Println(selectedTasks)

	// Step 1: Get employeeID from the database using employeePhoneNumber
	employeeID, err := dbConnection.GetEmployeeID(employeePhoneNumber)
	if err != nil {
		log.Println("Error getting employee ID:", err)
		// Handle the error
		return
	}

	// Step 2: Insert the lift information into the lifts table
	liftID, err := dbConnection.InsertLiftInfo(employeeID, elevatorName)
	if err != nil {
		log.Println("Error inserting lift information:", err)
		// Handle the error
		return
	}

	// Step 2b: Insert the lift details into the lift_details table
	// liftDetailsID, err := dbConnection.InsertLiftDetails(residentialComplex, elevatorName)//todo del me
	// if err != nil {
	// 	log.Println("Error inserting lift details:", err)
	// 	// Handle the error
	// 	return
	// }

	// Step 3: Insert the tasks into the tasks table
	err = dbConnection.InsertTasks(taskIDToName, employeeID, selectedTasks, liftID)
	if err != nil {
		log.Println("Error inserting tasks:", err)
		// Handle the error
		return
	}

	// Step 2a: Insert the project information into the projects table
	_, err = dbConnection.InsertProjectInfo(residentialComplex, employeeID, liftID)
	if err != nil {
		log.Println("Error inserting project information:", err)
		// Handle the error
		return
	}

	// Step 4: Insert a record into the lift_tasks table for each task
	// err = dbConnection.InsertLiftTasks(liftID, selectedTasks)
	// if err != nil {
	// 	log.Println("Error inserting lift tasks:", err)
	// 	// Handle the error
	// 	return
	// }
	msg = tgbotapi.NewMessage(message.Chat.ID, "Project created successfully!")
	bot.Send(msg)

}

// Function to remove the task button with the specified task ID from the button list
func removeTaskButton(taskButtons []tgbotapi.InlineKeyboardButton, taskID int) []tgbotapi.InlineKeyboardButton {
	for i, button := range taskButtons {
		if *button.CallbackData == strconv.Itoa(taskID) {
			// Remove the task button from the list
			taskButtons = append(taskButtons[:i], taskButtons[i+1:]...)
			break
		}
	}
	return taskButtons
}

func GetTasksFromLifts(dbConnection db.Database) ([]models.Task, error) {
	query := "SELECT id, task_name FROM task_of_lifts LIMIT 3;" //todo change it

	rows, err := dbConnection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)

	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(&task.ID, &task.TaskName)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

/*
	// Insert the project into the database
	//err = insertProject(projectName, dbConnection)
	//if err != nil {
	//	log.Println("Error creating project:", err)
	//	response := "Failed to create the project. Please try again later."
	//	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	//	bot.Send(msg)
	//	return
	//}

	msg = tgbotapi.NewMessage(message.Chat.ID, "Project created successfully!")
	bot.Send(msg)
*/

func collectUserInput(bot *tgbotapi.BotAPI, chatID int64, updates tgbotapi.UpdatesChannel, field string, dbConnection db.Database) (string, error) {

	// Create a channel to receive user input
	userInputChan := make(chan string)

	// Start a goroutine to collect user input
	go func() {
		for update := range updates {
			if update.Message != nil && update.Message.Chat.ID == chatID {
				userInput := update.Message.Text
				if userInput != "" {
					userInputChan <- userInput // Send the user input to the channel
					return
				} else {
					msg := tgbotapi.NewMessage(chatID, "Input cannot be empty. Please try again:")
					bot.Send(msg)
				}
			}
		}
		close(userInputChan) // Close the channel when the updates channel is closed
	}()

	// Receive user input from the channel
	userInput, ok := <-userInputChan
	if !ok {
		return "", fmt.Errorf("user input channel closed unexpectedly")
	}

	switch field {
	case "startDate", "endDate":
		if !isValidDate(userInput) {
			return "", fmt.Errorf("invalid date format. Please enter a valid date (dd/mm/yyyy)")
		}
	case "employeePhoneNumber":

		if !isValidPhoneNumber(userInput) || !isPhoneNumberExists(dbConnection, userInput) {
			return "", fmt.Errorf("invalid phone number format. Please enter a valid phone number")
		}
	case "complexName", "elevatorName", "description":
		// Additional validation or checks specific to complexName can be added here
	}

	return userInput, nil
}

func collectFile(bot *tgbotapi.BotAPI, chatID int64, updates tgbotapi.UpdatesChannel) (string, error) {
	// Create a channel to receive user input
	userInputChan := make(chan string)

	// Start a goroutine to collect user input
	go func() {
		for update := range updates {
			// Check if the message contains a photo or video
			if (update.Message.Photo != nil || update.Message.Video != nil) && update.Message.Chat.ID == chatID {
				// The message contains a photo or video, handle it accordingly
				fileID := ""
				if update.Message.Photo != nil {
					// Get the file ID of the photo
					fileID = (*update.Message.Photo)[0].FileID
					userInputChan <- fileID // Send the user input to the channel
					return
				} else if update.Message.Video != nil {
					// Get the file ID of the video
					fileID = update.Message.Video.FileID
					userInputChan <- fileID // Send the user input to the channel
					return
				}
			}

		}
		close(userInputChan) // Close the channel when the updates channel is closed
	}()
	// Receive user input from the channel
	userInput, ok := <-userInputChan
	if !ok {
		return "", fmt.Errorf("user input channel closed unexpectedly")
	}

	return userInput, nil
}

func isValidDate(date string) bool {
	_, err := time.Parse("02/01/2006", date)
	return err == nil
}

func isValidPhoneNumber(phoneNumber string) bool {
	// Remove any non-digit characters from the phone number
	phoneNumber = regexp.MustCompile(`\D`).ReplaceAllString(phoneNumber, "")

	// Check if the phone number has 11 digits
	if len(phoneNumber) != 11 {
		return false
	}

	// Check if the phone number starts with "7"
	if phoneNumber[:1] != "7" {
		return false
	}

	return true
}

func isPhoneNumberExists(dbConnection db.Database, phoneNumber string) bool {

	var count int
	rows, err := dbConnection.Query("select count(*) from workers where phone_number = ?", phoneNumber)
	if err != nil {

		return false
	}
	// fmt.Println("count: ", count)
	rows.Scan(&count)
	return count > 0
}

func createTaskInlineKeyboard(dbConnection db.Database) (tgbotapi.InlineKeyboardMarkup, error) {
	var inlineKeyboardRows []tgbotapi.InlineKeyboardButton

	rows, err := dbConnection.Query("select * from task_of_lifts;")
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	defer rows.Close()
	for rows.Next() {
		var taskID int
		var taskName string
		err := rows.Scan(&taskID, &taskName)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}

		callbackData := fmt.Sprintf("task:%d", taskID)
		inlineBtn := tgbotapi.NewInlineKeyboardButtonData(taskName, callbackData)
		inlineKeyboardRows = append(inlineKeyboardRows, inlineBtn)
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(inlineKeyboardRows)
	return inlineKeyboard, nil
}

func handleProjectList(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, updates tgbotapi.UpdatesChannel) {

	// msg := tgbotapi.NewMessage(message.Chat.ID, "Please enter the phone number of the responsible employee(example: 77078566392):")
	// bot.Send(msg)
	// employeePhoneNumber, err := collectUserInput(bot, message.Chat.ID, updates, "employeePhoneNumber", dbConnection)
	// if err != nil {
	// 	log.Println("Error collecting user input:", err)
	// 	msg := tgbotapi.NewMessage(message.Chat.ID, "Wrong data format")
	// 	bot.Send(msg)
	// 	return
	// }
	// fmt.Println("employeePhoneNumber: ", employeePhoneNumber)

	// Query the database to fetch the list of projects
	projectList, err := getProjectList(dbConnection)
	if err != nil {
		log.Println("Error retrieving project list:", err)
		response := "Failed to fetch the project list. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	response := "Project List:\n"
	for _, project := range projectList {
		response += strconv.Itoa(project.ID) + "- " + project.NameResident + " /seeProject" + strconv.Itoa(project.ID) + "\n"
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	bot.Send(msg)
}

func handleProjectInfo(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, updates tgbotapi.UpdatesChannel) {
	// Extract the project ID from the command arguments
	// phoneNumber := strings.TrimSpace(message.CommandArguments())
	msg := tgbotapi.NewMessage(message.Chat.ID, "Please enter the phone number of the responsible employee(example: 77078566392):")
	bot.Send(msg)
	employeePhoneNumber, err := collectUserInput(bot, message.Chat.ID, updates, "employeePhoneNumber", dbConnection)
	if err != nil {
		log.Println("Error collecting user input:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Wrong data format")
		bot.Send(msg)
		return
	}
	fmt.Println("employeePhoneNumber: ", employeePhoneNumber)
	ShowTasks(bot, message, dbConnection, employeePhoneNumber)
}

// Insert the project into the database
func insertProject(projectName string, dbConnection db.Database) error {
	// Implement your logic to insert the project into the projects table
	// You can use the db.ExecuteNonQuery function from the db package

	// Example query:
	query := "INSERT INTO projects (name) VALUES (?)"
	err := dbConnection.Execute(query, projectName)
	if err != nil {
		return err
	}
	return nil
}

// Retrieve the list of projects from the database
func getProjectList(dbConnection db.Database) ([]models.Project, error) {
	// Implement your logic to fetch the list of projects from the projects table
	// You can use the db.ExecuteQuery function from the db package

	// Example query:
	query := "SELECT id, name_resident FROM projects;"
	rows, err := dbConnection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectList []models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(&project.ID, &project.NameResident)
		if err != nil {
			return nil, err
		}
		projectList = append(projectList, project)
	}

	return projectList, nil
}

func handleProjectByID(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, updates tgbotapi.UpdatesChannel) {
	fmt.Println("Inside handleCompleteTask")

	// Extract the task ID from the command
	command := strings.TrimPrefix(strings.ToLower(message.Text), "/seeproject")
	projecID, err := strconv.Atoi(command)
	fmt.Println("ProjecID: ", projecID)

	if err != nil {
		log.Println("Invalid seeproject ID:", err)
		response := "Invalid projec ID. Please try again."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	// Update the task as completed in the database
	taskList, err := getProjectByID(dbConnection, projecID)
	if err != nil {
		log.Println("Error getting projec:", err)
		response := "Failed to getting projec. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}
	// Display the task information to the user
	if len(taskList) == 0 {
		// No tasks found for the specified employee
		msg := tgbotapi.NewMessage(message.Chat.ID, "No project found for the specified id.")
		bot.Send(msg)
	} else {
		// Tasks found, send the information as a formatted message

		// taskInfo := taskInfoList[0] // Get the first task info since the fields are the same for all tasks
		// taskInfoMessage := fmt.Sprintf("Elevator: %s\nResidential Complex: %s\nEmployee Phone Number: %s\n\n", taskInfo.ElevatorName, taskInfo.ResidentialComplex, taskInfo.EmployeePhoneNumber)
		var taskInfoMessage string
		for _, task := range taskList {
			doneSymbol := "✅" // Green checkmark
			if !task.IsDone {
				doneSymbol = "❌" // Red cross
			}
			taskInfoMessage += "\nElevator: " + task.ElevatorName + "\nResidential Complex:" + task.ResidentialComplex + "\nEmployee Phone Number:" + task.EmployeePhoneNumber + "\nTask: " + task.TaskName + "\nStart Date: " + task.StartDate + "\nEnd Date: " + task.EndDate + "\nIs Done: " + doneSymbol + "\n\n"
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, taskInfoMessage)
		bot.Send(msg)
	}

	response := ""
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	bot.Send(msg)
}

// Retrieve project information from the database
func getProjectByID(dbConnection db.Database, projectID int) ([]models.TaskInfo, error) {
	query := `
		SELECT
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
		WHERE p.id = ?
		ORDER BY p.id, t.id;
	`

	rows, err := dbConnection.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	taskInfoList := make([]models.TaskInfo, 0)

	for rows.Next() {
		var (
			elevatorName        string
			residentialComplex  string
			employeePhoneNumber string
			taskName            string
			startDateStr        []uint8
			endDateStr          []uint8
			isDone              bool
		)

		err := rows.Scan(
			&residentialComplex,
			&elevatorName,
			&employeePhoneNumber,
			&taskName,
			&startDateStr,
			&endDateStr,
			&isDone,
		)
		if err != nil {
			return nil, err
		}

		startDate, _ := time.Parse("2006-01-02", string(startDateStr))
		endDate, _ := time.Parse("2006-01-02", string(endDateStr))

		taskInfo := models.TaskInfo{
			ElevatorName:        elevatorName,
			ResidentialComplex:  residentialComplex,
			EmployeePhoneNumber: employeePhoneNumber,
			TaskName:            taskName,
			StartDate:           startDate.Format("02/01/2006"),
			EndDate:             endDate.Format("02/01/2006"),
			IsDone:              isDone,
		}

		taskInfoList = append(taskInfoList, taskInfo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return taskInfoList, nil
}

func handlerCompleteForceMajorOrChangereques(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, updates tgbotapi.UpdatesChannel, issueName string) {
	fmt.Println("Inside handlerCompleteForceMajor")
	// Extract the task ID from the command
	var command string
	switch issueName {
	case "completeforcemajor":
		command = strings.TrimPrefix(message.Text, "/completeforcemajor")
	case "completechangerequest":
		command = strings.TrimPrefix(message.Text, "/completechangerequest")
	}

	force_majeureID, err := strconv.Atoi(command)
	fmt.Println("forcemajor: ", force_majeureID)

	if err != nil {
		log.Println("Invalid forcemajor ID:", err)
		response := "Invalid forcemajor ID. Please try again."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "What did you do for solving this problem:")
	bot.Send(msg)
	descriptionOfWhatDid, err := collectUserInput(bot, message.Chat.ID, updates, "description", dbConnection)
	if err != nil {
		log.Println("Error collecting user input:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Wrong data format")
		bot.Send(msg)
		return
	}
	var response string
	switch issueName {
	case "completeforcemajor":
		response = "forcemajor marked as completed successfully."
		err = dbConnection.Execute("UPDATE force_majeure SET is_done = true, description_of_what_done = ? WHERE id = ?", descriptionOfWhatDid, force_majeureID)

	case "completechangerequest":
		response = "changerequest marked as completed successfully."
		err = dbConnection.Execute("UPDATE change_requests SET is_done = true, description_of_what_done = ? WHERE id = ?", descriptionOfWhatDid, force_majeureID)

	}

	if err != nil {
		log.Println("Error marking forcemajor as completed:", err)
		response := "Failed to mark forcemajor as completed. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, response)
	bot.Send(msg)
}

func handlerCompleteOverdueTaskByAdmin(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, updates tgbotapi.UpdatesChannel) {
	command := strings.TrimPrefix(message.Text, "/CompleteOverdueTaskByAdmin")

	overdueTaskID, err := strconv.Atoi(command)
	fmt.Println("overdueTasID: ", overdueTaskID)

	if err != nil {
		log.Println("Invalid overdueTasID:", err)
		response := "Invalid overdueTasID. Please try again."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "What did you do for solving this problem:")
	bot.Send(msg)
	descriptionOfWhatDid, err := collectUserInput(bot, message.Chat.ID, updates, "description", dbConnection)
	if err != nil {
		log.Println("Error collecting user input:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Wrong data format")
		bot.Send(msg)
		return
	}

	// Mark the overdue task as done by the admin
	err = MarkOverdueTaskDoneByAdmin(dbConnection, descriptionOfWhatDid, overdueTaskID)
	if err != nil {
		log.Println("Error marking overdue task as done by admin:", err)
		// Handle the error
		return
	}

	// Send a confirmation message to the admin
	response := fmt.Sprintf("Overdue task with ID %d has been marked as done by the admin.", overdueTaskID)
	msg = tgbotapi.NewMessage(message.Chat.ID, response)
	bot.Send(msg)

}

func MarkOverdueTaskDoneByAdmin(dbConnection db.Database, descriptionOfWhatDid string, overdueTaskID int) error {
	query := "UPDATE overdue_task SET is_done_by_admin = true, description_by_admin = ? WHERE id = ?"
	err := dbConnection.Execute(query, descriptionOfWhatDid, overdueTaskID)
	if err != nil {
		return err
	}

	query = `SELECT task_id FROM overdue_task WHERE id = ?`
	row := dbConnection.QueryRow(query, overdueTaskID)
	if err != nil {
		return err
	}
	var id int
	err = row.Scan(&id)
	if err != nil {
		return err
	}

	query = "UPDATE tasks SET isDone = true WHERE id = ?"
	err = dbConnection.Execute(query, id)
	if err != nil {
		return err
	}

	return nil
}
