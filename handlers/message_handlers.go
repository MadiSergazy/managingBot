package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"

	"madi_telegram_bot/db"
)

func HandleUserMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, dbConnection db.Database, updates tgbotapi.UpdatesChannel) {
	switch {

	case message.Text == "/start":

		// Create a reply markup with a request for the user's phone number
		replyMarkup := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.KeyboardButton{
					Text:            "Share Phone Number",
					RequestContact:  true,  // Request contact information from the user
					RequestLocation: false, // Not requesting location
				},
			),
		)

		// Send a message to the user with the phone number request button
		response := "Please click the button below to share your phone number."
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		msg.ReplyMarkup = replyMarkup
		bot.Send(msg)

	case message.Contact != nil:
		// Retrieve the user's phone number from the incoming message
		fmt.Println("Userphone: ", message.Contact.PhoneNumber)
		userPhone := message.Contact.PhoneNumber //todo: store this information for convinient using admin's

		fmt.Println("UserID: ", message.From.ID) //toDO: store thsi info in database for identifying workers and admins

		userName := message.From.UserName
		fmt.Println("userName: ", message.From.UserName)

		// Check if the phone number exists in the admins table
		isAdmin := isWhoUser(userPhone, dbConnection, "admins")
		isHrManager := isWhoUser(userPhone, dbConnection, "hr_manager")
		// If the user is an admin, mark them as such
		if isAdmin {
			fmt.Println("It is admin and i am gonna insert identifier")

			if err := dbConnection.InsertIdentifier(message.From.ID, "admins", userPhone, userName); err != nil {
				log.Println("Error inserting Admin:", err)
			}
			log.Printf("User with phone number %s and userID: %d is an admin", userPhone, message.From.ID)

		}
		if isHrManager {
			fmt.Println("It is HRMANAGER and i am gonna insert identifier")

			if err := dbConnection.InsertIdentifier(message.From.ID, "hr_manager", userPhone, userName); err != nil {
				log.Println("Error inserting hr_manager:", err)
			}
			log.Printf("User with phone number %s and userID: %d is an hr_manager", userPhone, message.From.ID)

		} else {
			// Insert the phone number into the workers table as a regular worker
			err := insertWorker(userPhone, userName, dbConnection)
			if err != nil {
				log.Println("Error inserting worker:", err)
				//return
			}

			if err := dbConnection.InsertIdentifier(message.From.ID, "workers", userPhone, userName); err != nil {
				log.Println("Error inserting Admin:", err)
			}
			log.Printf("User with phone number %s and userID: %d is a worker", userPhone, message.From.ID)

		}

		// Send a response to the user
		response := "Registration successful!"
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		bot.Send(msg)
	default:
		userIdentifier := message.From.ID
		isAdmin := isByIdentifier(userIdentifier, dbConnection, "admins")
		isHrManager := isByIdentifier(userIdentifier, dbConnection, "hr_manager")
		// If the user is an admin, mark them as such
		if isAdmin {
			log.Info(" It is a admin")
			HandleAdminCommand(bot, message, dbConnection, updates) //todo: Implemet me
		} else if isHrManager {
			fmt.Println("IT IS a isHrManager")
			//todo: Implemet worker
			// HandleHrManagerCommand(bot, message, updateBuffer, dbConnection, updates) //todo: Implemet me
		} else {
			fmt.Println("IT IS a worker")
			//todo: Implemet worker
			HandleWorkerCommand(bot, message, dbConnection, updates) //todo: Implemet me
		}

	}

}

func isByIdentifier(userIdentifier int, dbConnection db.Database, tableName string) bool {
	query := "SELECT COUNT(*) FROM " + tableName + " WHERE identifier = ?"
	row := dbConnection.QueryRow(query, userIdentifier)
	if row.Err() != nil {
		log.Println("Error executing query in isByIdentifier:", row.Err())
		return false
	}

	var count int

	if err := row.Scan(&count); err != nil {
		log.Println("Error scanning row:", err)
		return false
	}
	return count > 0
}

// Check if the phone number exists in the admins table
func isWhoUser(phoneNumber string, dbConnection db.Database, role string) bool {
	// Query the admins table to check if the phone number exists
	// Implement your logic to check if the phone number exists in the admins table
	// You can use the db.ExecuteQuery function from the db package

	// Example query:
	query := `SELECT COUNT(*) FROM ` + role + ` WHERE phone_number = ?`
	row := dbConnection.QueryRow(query, phoneNumber)
	if row.Err() != nil {
		log.Println("Error executing query in isWhoUser:", row.Err())
		return false
	}

	var count int

	if err := row.Scan(&count); err != nil {
		log.Println("Error scanning row:", err)
		return false
	}
	return count > 0
}

// Insert the phone number into the workers table
func insertWorker(phoneNumber string, nameWorker string, dbConnection db.Database) error {
	// Insert the phone number into the workers table
	// Implement your logic to insert the phone number into the workers table
	// You can use the db.ExecuteNonQuery function from the db package

	// Example query:
	query := "INSERT INTO workers (phone_number, name) VALUES (?, ?)"
	err := dbConnection.Execute(query, phoneNumber, nameWorker)
	if err != nil {
		//if err == Exists //todo add err data exists in db
		return err
	}
	return nil
}
