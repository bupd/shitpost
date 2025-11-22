package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// get bot token from env var
	token := os.Getenv("BOT_TOKEN")
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
		caption := msg.Caption
		log.Printf("\nGot message from %s: with message: %v", update.Message.From.UserName, msg)
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
			handleFile(bot, photo.FileID, msg.Chat.ID, caption)
		}

		// VIDEO HANDLING
		if msg.Video != nil {
			handleFile(bot, msg.Video.FileID, msg.Chat.ID, caption)
		}

		// DOCUMENT HANDLING (fallback for files)
		if msg.Document != nil {
			handleFile(bot, msg.Document.FileID, msg.Chat.ID, caption)
		}
	}
}

// handleFile downloads + saves + sends back the file
func handleFile(bot *tgbot.BotAPI, fileID string, chatID int64, caption string) {
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

	cleanCaption, altText := ParseCaptionAlt(caption)

	log.Println("Posting Updates via crosspost")
	log.Printf("\n run cmd: crosspost -bmtl --image %s --image-alt '%s' '%s'", savePath, altText, cleanCaption)

	go PostViaCrosspost(bot, chatID, savePath, altText, caption)
}

// PostViaCrosspost runs crosspost, captures logs, then sends back the file
func PostViaCrosspost(bot *tgbot.BotAPI, chatID int64, savePath, altText, caption string) {
	// build command
	cmd := exec.Command("crosspost", "-bmt", "--image", savePath, "--image-alt", altText, caption)
	// cmd := exec.Command("crosspost", "-b", "--image", savePath, caption)

	// capture stdout & stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("[XPOST] failed to get stdout:", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println("[XPOST] failed to get stderr:", err)
		return
	}

	// start command
	log.Printf("[XPOST] running crosspost for %s", savePath)
	if err := cmd.Start(); err != nil {
		log.Println("[XPOST] failed starting:", err)
		return
	}

	// read output streams
	outBytes, _ := io.ReadAll(stdout)
	errBytes, _ := io.ReadAll(stderr)

	// wait until exit
	if err := cmd.Wait(); err != nil {
		log.Println("[XPOST] cmd finished with error:", err)
	}

	// format response caption
	newCaption := fmt.Sprintf(
		"%s\n\n=== crosspost logs ===\nstdout:\n%s\n\nstderr:\n%s",
		caption,
		string(outBytes),
		string(errBytes),
	)

	// send same file back
	send := tgbot.NewDocument(chatID, tgbot.FilePath(savePath))
	send.Caption = newCaption

	if _, err := bot.Send(send); err != nil {
		log.Println("failed sending message:", err)
	}

	log.Println("[XPOST] file sent back successfully")
}

// ParseCaptionAlt extracts the alt text from a caption.
// Format expected: "... text ... \nalt: something here"
func ParseCaptionAlt(caption string) (cleanCaption, altText string) {
	// Split by newlines
	lines := strings.Split(strings.TrimSpace(caption), "\n")

	// last line
	last := strings.TrimSpace(lines[len(lines)-1])

	// Check if alt exists
	if strings.HasPrefix(strings.ToLower(last), "alt:") {
		altText = strings.TrimSpace(strings.TrimPrefix(last, "alt:"))

		// remove last line for clean caption
		cleanCaption = strings.TrimSpace(strings.Join(lines[:len(lines)-1], "\n"))
		return cleanCaption, altText
	}

	// no alt found → return original
	return caption, ""
}
