package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/joho/godotenv/autoload"
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

		msg := update.Message
		text := msg.Text
		log.Printf("\nGot message from %s: with message: %s", update.Message.From.UserName, msg)
		log.Printf("\nGot message text: %s", text)

		// TEXT HANDLING
		if msg.Text != "" {
			reply := tgbot.NewMessage(msg.Chat.ID, "bot says that "+msg.From.UserName+"said, "+msg.Text)
			// send message
			if _, err := bot.Send(reply); err != nil {
				log.Println("failed sending message:", err)
			}
		}

		// PHOTO HANDLING  (Telegram sends photos as array sorted by size)
		if len(msg.Photo) > 0 {
			photo := msg.Photo[len(msg.Photo)-1] // pick highest resolution
			handleFile(bot, photo.FileID, msg.Chat.ID)
		}

		// VIDEO HANDLING
		if msg.Video != nil {
			handleFile(bot, msg.Video.FileID, msg.Chat.ID)
		}

		// DOCUMENT HANDLING (fallback for files)
		if msg.Document != nil {
			handleFile(bot, msg.Document.FileID, msg.Chat.ID)
		}

	}
}

// handleFile downloads + saves + sends back the file
func handleFile(bot *tgbot.BotAPI, fileID string, chatID int64) {
	file, err := bot.GetFile(tgbot.FileConfig{FileID: fileID})
	if err != nil {
		log.Println("error getting file:", err)
		return
	}

	// construct direct download URL
	url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)

	// create storage folder if needed
	os.MkdirAll("downloads", 0755)

	// determine local save path
	filename := filepath.Base(file.FilePath)
	savePath := filepath.Join("downloads", filename)

	log.Printf("downloading file: %s, to: %s", url, savePath)

	// download
	resp, err := http.Get(url)
	if err != nil {
		log.Println("download err:", err)
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(savePath)
	if err != nil {
		log.Println("file create err:", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("write err:", err)
		return
	}

	log.Println("saved file:", savePath)

	// SEND BACK THE SAME FILE
	send := tgbot.NewDocument(chatID, tgbot.FilePath(savePath))
	send.Caption = "here’s your file back"
	// bot.Send(send)

	// send message
	if _, err := bot.Send(send); err != nil {
		log.Println("failed sending message:", err)
	}
}
