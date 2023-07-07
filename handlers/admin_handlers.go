package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"madi_telegram_bot/db"
	"strconv"
	"strings"
)

// HandleAdminCommand handles admin-specific commands
func HandleAdminCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updateBuffer []tgbotapi.Update, dbConnection db.Database, updates tgbotapi.UpdatesChannel) {
	response := "Unknown command. Please use valid admin commands."
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	message.CommandArguments()
	command := strings.ToLower(message.Command())
	switch command {
	case "createproject":
		//msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
		handleCreateProject(bot, message, updateBuffer, dbConnection, updates)
	case "projectlist":
		//msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
		handleProjectList(bot, message, dbConnection) //todo to me
	case "projectinfo":
		//msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true} todo think about it
		handleProjectInfo(bot, message, dbConnection) //todo do me
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

type TaskDates struct {
	StartDate string
	EndDate   string
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
	residentialComplex, err := collectUserInput(bot, message.Chat.ID, updates)
	if err != nil {
		log.Println("Error collecting user input:", err)
		// Handle the error
		return
	}

	fmt.Println("Residential Complex:", residentialComplex)

	msg = tgbotapi.NewMessage(message.Chat.ID, "Please enter the name of the elevator:")
	bot.Send(msg)
	elevatorName, err := collectUserInput(bot, message.Chat.ID, updates)
	if err != nil {
		log.Println("Error collecting user input:", err)
		// Handle the error
		return
	}
	fmt.Println("elevatorName: ", elevatorName)

	msg = tgbotapi.NewMessage(message.Chat.ID, "Please enter the phone number of the responsible employee(example: 7078566392):")
	bot.Send(msg)
	employeePhoneNumber, err := collectUserInput(bot, message.Chat.ID, updates)
	if err != nil {
		log.Println("Error collecting user input:", err)
		// Handle the error
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
	selectedTasks := make(map[int]TaskDates)
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
				startDate, err := collectUserInput(bot, update.CallbackQuery.Message.Chat.ID, updates)
				if err != nil {
					log.Println("Error collecting user input:", err)
					// Handle the error
					return
				}
				fmt.Println("startDate: ", startDate)

				// Send a message to the user to set the end date
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Please enter the end date (dd/mm/yyyy):")
				bot.Send(msg)

				// Wait for the user to enter the end date
				endDate, err := collectUserInput(bot, update.CallbackQuery.Message.Chat.ID, updates)
				if err != nil {
					log.Println("Error collecting user input:", err)
					// Handle the error
					return
				}
				fmt.Println("endDate: ", endDate)

				// Store the selected task and its dates in the map
				selectedTasks[selectedTaskID] = TaskDates{
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

	// Insert the tasks, start dates, and end dates into the database
	/*for taskID, dates := range selectedTasks {
		taskName := taskIDToName[taskID]
		startDate := dates.StartDate
		endDate := dates.EndDate

		// TODO: Insert the task, start date, and end date into the database
	}*/
	fmt.Println(selectedTasks)

	/*
		// Wait for the user to select a task
		for update := range updates {
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
				// Create an inline keyboard markup with the task buttons
				rows := make([][]tgbotapi.InlineKeyboardButton, 0) //todo rewtire it to another function
				for _, button := range taskButtons {
					row := []tgbotapi.InlineKeyboardButton{button}
					rows = append(rows, row)
				}
				inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)

				// Send a message to the user to set the start date
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("You have selected the task: %s\nPlease enter the start date (dd/mm/yyyy):", selectedTaskName))
				bot.Send(msg)

				// Wait for the user to enter the start date
				startDate, err := collectUserInput(bot, update.CallbackQuery.Message.Chat.ID, updates)
				if err != nil {
					log.Println("Error collecting user input:", err)
					// Handle the error
					return
				}
				fmt.Println("startDate: ", startDate)

				// Send a message to the user to set the end date
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Please enter the end date (dd/mm/yyyy):")
				bot.Send(msg)

				// Wait for the user to enter the end date
				endDate, err := collectUserInput(bot, update.CallbackQuery.Message.Chat.ID, updates)
				if err != nil {
					log.Println("Error collecting user input:", err)
					// Handle the error
					return
				}
				fmt.Println("endDate: ", endDate)
				// TODO: Insert the task, start date, and end date into the database

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

				break
			}
		}
	*/
	// Send the success message to the user
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

type Task struct {
	ID       int
	TaskName string
}

func GetTasksFromLifts(dbConnection db.Database) ([]Task, error) {
	query := "SELECT id, task_name FROM task_of_lifts;"

	rows, err := dbConnection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)

	for rows.Next() {
		task := Task{}
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

func collectUserInput(bot *tgbotapi.BotAPI, chatID int64, updates tgbotapi.UpdatesChannel) (string, error) {

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
	residentialComplex, ok := <-userInputChan
	if !ok {
		return "", fmt.Errorf("user input channel closed unexpectedly")
	}

	return residentialComplex, nil
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

func handleProjectList(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database) {
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
		response += "- " + project.Name + "\n"
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	bot.Send(msg)
}

func handleProjectInfo(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database) {
	// Extract the project ID from the command arguments
	projectID, err := strconv.Atoi(strings.TrimSpace(message.CommandArguments()))
	if err != nil {
		response := "Invalid project ID. Please provide a valid project ID."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	// Query the database to fetch project information
	project, err := getProjectInfo(projectID, dbConnection)
	if err != nil {
		log.Println("Error retrieving project information:", err)
		response := "Failed to fetch project information. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
		return
	}

	response := "Project Information:\n"
	response += "ID: " + strconv.Itoa(project.ID) + "\n"
	response += "Name: " + project.Name + "\n"
	// Add more project details as needed

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	bot.Send(msg)
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
func getProjectList(dbConnection db.Database) ([]Project, error) {
	// Implement your logic to fetch the list of projects from the projects table
	// You can use the db.ExecuteQuery function from the db package

	// Example query:
	query := "SELECT id, name FROM projects"
	rows, err := dbConnection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectList []Project
	for rows.Next() {
		var project Project
		err := rows.Scan(&project.ID, &project.Name)
		if err != nil {
			return nil, err
		}
		projectList = append(projectList, project)
	}

	return projectList, nil
}

// Retrieve project information from the database
func getProjectInfo(projectID int, dbConnection db.Database) (*Project, error) {
	// Implement your logic to fetch project information from the projects table
	// You can use the db.QueryRow function from the db package

	// Example query:
	query := "SELECT id, name FROM projects WHERE id = ?"
	row := dbConnection.QueryRow(query, projectID)

	var project Project
	err := row.Scan(&project.ID, &project.Name)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

// Define a struct to represent a project
type Project struct {
	ID   int
	Name string
	// Add more fields as needed
}
