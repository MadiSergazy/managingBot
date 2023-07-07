package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"madi_telegram_bot/config"
	"madi_telegram_bot/db"
	"madi_telegram_bot/handlers"
)

// Handles the Telegram bot functionality. The bot.go file encapsulates the creation of the Telegram bot, sets up event handlers for different types of messages, and dispatches them to the appropriate handlers.
func StartBot(cfg config.Config, dbConnection db.Database) error {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	updateBuffer := make([]tgbotapi.Update, 0)
	// Set up event handlers for different types of messages
	for update := range updates {
		if update.Message != nil {
			// Handle text messages
			log.Info(update.Message)
			log.Info(update.Message.From.ID)

			// Add the update to the buffer
			updateBuffer = append(updateBuffer, update)

			handlers.HandleUserMessage(bot, update.Message, updateBuffer, dbConnection, updates)
		} else if update.CallbackQuery != nil {
			// Handle button clicks
			//handlers.HandleButtonCallback(bot, update.CallbackQuery)
		}
	}

	return nil

}
