package main

import (
	"log"
	"os"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// get bot token from env var
	token := os.Getenv("BOT_TOKEN") // <-- changed: use env variable for safety
	if token == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}

	// initialize bot instance
	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		log.Fatal("failed creating bot: ", err)
	}

	log.Printf("Authorized as @%s", bot.Self.UserName)

	// setup update config
	updateConfig := tgbot.NewUpdate(0)
	updateConfig.Timeout = 30

	// start receiving updates
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		// ignore anything that is not a message
		if update.Message == nil {
			continue // <-- changed: skip non-message updates
		}

		text := update.Message.Text

		// create reply
		reply := tgbot.NewMessage(update.Message.Chat.ID, "bot says, "+text)

		// send message
		if _, err := bot.Send(reply); err != nil {
			log.Println("failed sending message:", err) // <-- changed: log errors instead of crash
		}
	}
}
