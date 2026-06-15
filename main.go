package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/joho/godotenv/autoload"
)

const telegramMessageLimit = 3900

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}
	authorizedUsers := authorizedTelegramUsers()
	crosspostFlags := crosspostFlags()
	dryRun := envBool("SHITPOST_DRY_RUN")

	// initialize bot instance
	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		log.Fatal("failed creating bot: ", err)
	}

	log.Printf("Authorized as @%s", bot.Self.UserName)
	if dryRun {
		log.Println("Dry-run mode enabled. Messages will not be posted.")
	}

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

		if !isAuthorized(update.Message.From, authorizedUsers) {
			reply := tgbot.NewMessage(update.Message.From.ID, "This bot is private.")
			log.Printf("Intruder: %s", update.Message.From.UserName)
			if _, err := bot.Send(reply); err != nil {
				log.Println("failed sending message:", err)
			}
			continue
		}

		msg := update.Message
		text := msg.Text
		caption := msg.Caption
		log.Printf("\nGot message from %s: with message: %v", update.Message.From.UserName, msg)
		log.Printf("\nGot message text: %s", text)

		// TEXT HANDLING
		if msg.Text != "" {
			go PostViaCrosspost(bot, msg.Chat.ID, crosspostFlags, "", "", text, dryRun)
		}

		// PHOTO HANDLING  (Telegram sends photos as array sorted by size)
		if len(msg.Photo) > 0 {
			photo := msg.Photo[len(msg.Photo)-1] // pick highest resolution
			handleFile(bot, photo.FileID, msg.Chat.ID, caption, "image/jpeg", crosspostFlags, dryRun)
		}

		// VIDEO HANDLING
		if msg.Video != nil {
			handleFile(bot, msg.Video.FileID, msg.Chat.ID, caption, msg.Video.MimeType, crosspostFlags, dryRun)
		}

		// DOCUMENT HANDLING (fallback for files)
		if msg.Document != nil {
			handleFile(bot, msg.Document.FileID, msg.Chat.ID, caption, msg.Document.MimeType, crosspostFlags, dryRun)
		}
	}
}

func envBool(name string) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
	return value == "1" || value == "true" || value == "yes"
}

func authorizedTelegramUsers() map[string]bool {
	configured := strings.TrimSpace(os.Getenv("AUTHORIZED_TELEGRAM_USERS"))
	users := map[string]bool{}

	if configured == "" {
		return users
	}

	for _, user := range strings.Split(configured, ",") {
		user = strings.TrimSpace(strings.TrimPrefix(user, "@"))
		if user != "" {
			users[strings.ToLower(user)] = true
		}
	}

	return users
}

func isAuthorized(user *tgbot.User, authorizedUsers map[string]bool) bool {
	if len(authorizedUsers) == 0 {
		return true
	}

	if authorizedUsers[strconv.FormatInt(user.ID, 10)] {
		return true
	}

	return authorizedUsers[strings.ToLower(user.UserName)]
}

func crosspostFlags() []string {
	configured := strings.TrimSpace(os.Getenv("CROSSPOST_FLAGS"))
	if configured == "" {
		return []string{"-bmt"}
	}

	return strings.Fields(configured)
}

// handleFile downloads + saves + sends back the file
func handleFile(bot *tgbot.BotAPI, fileID string, chatID int64, caption string, mediaType string, crosspostFlags []string, dryRun bool) {
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
	if mediaType == "" {
		mediaType = mime.TypeByExtension(filepath.Ext(filename))
	}

	log.Printf("downloading file: %s, to: %s", redactTelegramToken(url), savePath)

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

	if !strings.HasPrefix(mediaType, "image/") {
		warning := fmt.Sprintf("Downloaded %s, but this bot can only attach image media with the installed crosspost CLI. Posting caption text only.", mediaType)
		log.Println(warning)
		if cleanCaption != "" {
			go PostViaCrosspost(bot, chatID, crosspostFlags, "", "", cleanCaption, dryRun)
		} else {
			sendReply(bot, chatID, warning)
		}
		return
	}

	log.Println("Posting updates via crosspost")

	go PostViaCrosspost(bot, chatID, crosspostFlags, savePath, altText, cleanCaption, dryRun)
}

// PostViaCrosspost runs crosspost, captures logs, then sends back the file
func PostViaCrosspost(bot *tgbot.BotAPI, chatID int64, crosspostFlags []string, savePath, altText, caption string, dryRun bool) {
	args := append([]string{}, crosspostFlags...)
	if savePath == "" {
		log.Println("PostViaCrosspost: img not found, sending as text tweet")
	} else {
		args = append(args, "--image", savePath, "--image-alt", altText)
	}
	args = append(args, caption)

	if dryRun {
		preview := fmt.Sprintf("DRY RUN: would run `%s`", commandPreview("crosspost", args))
		log.Println(preview)
		sendReply(bot, chatID, preview)
		return
	}

	cmd := exec.Command("crosspost", args...)
	cmd.Env = normalizedCrosspostEnv()

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
	log.Printf("[XPOST] running crosspost for img: %s, alt: %s, caption: %s", savePath, altText, caption)
	if err := cmd.Start(); err != nil {
		log.Println("[XPOST] failed starting:", err)
		return
	}

	// read output streams
	outBytes, _ := io.ReadAll(stdout)
	errBytes, _ := io.ReadAll(stderr)

	status := "crosspost completed"

	// wait until exit
	if err := cmd.Wait(); err != nil {
		log.Println("[XPOST] cmd finished with error:", err)
		status = fmt.Sprintf("crosspost failed: %v", err)
	}

	logMessage := fmt.Sprintf(
		"%s\n\n%s\n\n=== crosspost logs ===\nstdout:\n%s\n\nstderr:\n%s",
		caption,
		status,
		string(outBytes),
		string(errBytes),
	)

	if savePath != "" {
		send := tgbot.NewDocument(chatID, tgbot.FilePath(savePath))
		send.Caption = caption
		if _, err := bot.Send(send); err != nil {
			log.Println("failed sending message:", err)
		}
		log.Println("[XPOST] file sent back successfully")
	}

	sendLongReply(bot, chatID, logMessage)

}

func commandPreview(name string, args []string) string {
	quoted := []string{name}
	for _, arg := range args {
		quoted = append(quoted, strconv.Quote(arg))
	}
	return strings.Join(quoted, " ")
}

func normalizedCrosspostEnv() []string {
	env := envMap(os.Environ())
	setDefaultEnv(env, "TWITTER_AUTH_TOKEN", "AUTH_TOKEN")
	setDefaultEnv(env, "TWITTER_API_CONSUMER_KEY", "consumer_key", "TWITTER_CONSUMER_KEY")
	setDefaultEnv(env, "TWITTER_API_CONSUMER_SECRET", "consumer_key_secret", "TWITTER_CONSUMER_SECRET")
	setDefaultEnv(env, "TWITTER_ACCESS_TOKEN_KEY", "access_token", "access_token_key", "TWITTER_ACCESS_TOKEN")
	setDefaultEnv(env, "TWITTER_ACCESS_TOKEN_SECRET", "access_token_secret", "TWITTER_ACCESS_SECRET")

	result := make([]string, 0, len(env))
	for key, value := range env {
		result = append(result, key+"="+value)
	}

	return result
}

func envMap(values []string) map[string]string {
	env := map[string]string{}
	for _, value := range values {
		key, val, ok := strings.Cut(value, "=")
		if ok {
			env[key] = val
		}
	}
	return env
}

func setDefaultEnv(env map[string]string, target string, aliases ...string) {
	if strings.TrimSpace(env[target]) != "" {
		return
	}

	for _, alias := range aliases {
		if value := strings.TrimSpace(env[alias]); value != "" {
			env[target] = value
			return
		}
	}
}

func sendReply(bot *tgbot.BotAPI, chatID int64, text string) {
	reply := tgbot.NewMessage(chatID, text)
	if _, err := bot.Send(reply); err != nil {
		log.Println("failed sending message:", err)
	}
}

func sendLongReply(bot *tgbot.BotAPI, chatID int64, text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}

	for len(text) > telegramMessageLimit {
		chunk := text[:telegramMessageLimit]
		if index := strings.LastIndex(chunk, "\n"); index > 0 {
			chunk = chunk[:index]
		}
		sendReply(bot, chatID, chunk)
		text = strings.TrimSpace(strings.TrimPrefix(text, chunk))
	}

	sendReply(bot, chatID, text)
}

func redactTelegramToken(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "<telegram-file-url>"
	}

	parts := strings.Split(parsed.Path, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, "bot") {
			parts[i] = "bot<redacted>"
		}
	}
	parsed.Path = strings.Join(parts, "/")
	return parsed.String()
}

// ParseCaptionAlt extracts the alt text from a caption.
// Format expected: "... text ... \nalt: something here"
func ParseCaptionAlt(caption string) (cleanCaption, altText string) {
	caption = strings.TrimSpace(caption)
	if caption == "" {
		return "", ""
	}

	// Split by newlines
	lines := strings.Split(caption, "\n")

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
