package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"madi_telegram_bot/bot"
	"madi_telegram_bot/config"
	"madi_telegram_bot/db"
)

// The entry point of your application. It initializes the necessary components, such as the configuration, database, and Telegram bot. It sets up the necessary event loop to handle incoming messages.
func main() {
	cfg := config.LoadConfig()

	dbConnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	fmt.Println("DB Connection:", dbConnection)

	botToken := cfg.BotToken
	fmt.Println("Bot Token:", botToken)

	dbConn, err := db.NewDatabase(cfg)
	if err != nil {
		log.Error("error getting connect to db: ", err)
	}

	// Rest of your code...

	err = bot.StartBot(cfg, *dbConn)
	if err != nil {
		log.Error("Error when starting tbot: ", err)
	}
	log.Info("Running bot with username: %s\n", "https://t.me/lift_kz_bot")
}
