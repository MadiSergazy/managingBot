package telegram

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"

	"madi_telegram_bot/config"
	"madi_telegram_bot/db"
	"madi_telegram_bot/handlers"
)

// create constuctor for bot
type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{bot: bot}
}

// Handles the Telegram bot functionality. The bot.go file encapsulates the creation of the Telegram bot, sets up event handlers for different types of messages, and dispatches them to the appropriate handlers.
func (b *Bot) StartBot(cfg config.Config, dbConnection db.Database) error {

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	// time.Sleep(time.Millisecond * 500)
	// updates.Clear()

	updateBuffer := make([]tgbotapi.Update, 0)
	StartForceMajeureNotifications(dbConnection, b.bot)
	// Set up event handlers for different types of messages
	b.handleUpdates(updates, updateBuffer, dbConnection)
	return nil

}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel, updateBuffer []tgbotapi.Update, dbConnection db.Database) {
	for update := range updates {
		if update.Message != nil {
			// Handle text messages
			log.Info(update.Message)
			log.Info(update.Message.From.ID)

			// Add the update to the buffer
			updateBuffer = append(updateBuffer, update)

			handlers.HandleUserMessage(b.bot, update.Message, updateBuffer, dbConnection, updates)
		} else if update.CallbackQuery != nil {
			// Handle button clicks
			//handlers.HandleButtonCallback(bot, update.CallbackQuery)
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u) //todo switch to webhook in future
}

func StartForceMajeureNotifications(dbConnection db.Database, bot *tgbotapi.BotAPI) {
	// Create a ticker with the desired duration

	ticker := time.NewTicker(24 * time.Hour) // Check every 24 hours

	// Start a goroutine to perform the notifications
	go func() {
		for {
			select {
			case <-ticker.C:
				// Perform the check for unfinished force majeure reports
				/*	err := handlers.CheckUnfinished(dbConnection, bot, "force_majeure")
					if err != nil {
						log.Println("Error checking unfinished force majeure reports:", err)
					}
					err = handlers.CheckUnfinished(dbConnection, bot, "change_requests")
					if err != nil {
						log.Println("Error checking unfinished change reports:", err)
					}*/
				err := handlers.CheckOverdueTasks(dbConnection, bot)
				if err != nil {
					log.Println("Error checking CheckOverdueTasks:", err)
				}

			}
		}
	}()
}
