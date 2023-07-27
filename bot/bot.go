package telegram

import (
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"

	"madi_telegram_bot/config"
	"madi_telegram_bot/db"
	"madi_telegram_bot/handlers"
)

// create constuctor for bot
type Bot struct {
	bot         *tgbotapi.BotAPI
	updateMutex sync.Mutex
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{bot: bot,
		updateMutex: sync.Mutex{}}
}

// Handles the Telegram bot functionality. The bot.go file encapsulates the creation of the Telegram bot, sets up event handlers for different types of messages, and dispatches them to the appropriate handlers.
func (b *Bot) StartBot(cfg config.Config, dbConnection db.Database) error {

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 500) //for cleaning updates that was send when bot was inactive
	updates.Clear()

	//search change_requests, force_majeure, expired tasks and expired validations
	StartSearching(dbConnection, b.bot)
	// Set up event handlers for different types of messages
	b.handleUpdates(updates, dbConnection)
	return nil

}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel, dbConnection db.Database) {
	for update := range updates {

		// Lock the mutex to ensure only one instance processes updates at a time
		b.updateMutex.Lock()
		if update.Message != nil {
			// Handle text messages
			log.Info(update.Message)
			log.Info(update.Message.From.ID)

			handlers.HandleUserMessage(b.bot, update.Message, dbConnection, updates)
		} else if update.CallbackQuery != nil {
			// Handle button clicks
			//handlers.HandleButtonCallback(bot, update.CallbackQuery)
		}
		// Unlock the mutex after processing the update
		b.updateMutex.Unlock()
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u) //todo switch to webhook in future
}

func StartSearching(dbConnection db.Database, bot *tgbotapi.BotAPI) {
	if bot == nil {
		log.Println("Database connection is null")
		return
	}
	// Create a ticker with the desired duration
	ticker := time.NewTicker(24 * time.Hour) // Check every 24 hours
	// ticker := time.NewTicker(20 * time.Second)
	// Start a goroutine to perform the notifications
	go func() {
		for {
			select {
			case <-ticker.C:
				// Perform the check for unfinished force majeure reports
				err := handlers.CheckUnfinished(dbConnection, bot, "force_majeure")
				if err != nil {
					log.Println("Error checking unfinished force majeure reports:", err)
				}
				err = handlers.CheckUnfinished(dbConnection, bot, "change_requests")
				if err != nil {
					log.Println("Error checking unfinished change reports:", err)
				}
				err = handlers.CheckOverdueTasks(dbConnection, bot)
				if err != nil {
					log.Println("Error checking CheckOverdueTasks:", err)
				}
				err = handlers.CheckPendingValidationTasks(dbConnection, bot)
				if err != nil {
					log.Println("Error checking pending validation tasks:", err)
				}
			}
		}
	}()

}
