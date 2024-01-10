package main

import (
	"log"
	"os"
	"strings"

	//add dotenv
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var membersDetails = make(map[int64]MemberDetails)

type MemberDetails struct {
	GitHubUsername string
	LinkedInURL    string
	PhoneNumber    string
}

func main() {
	//add dotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.Chat.ID

		if update.Message.NewChatMembers != nil {
			greetNewUser(bot, userID, update.Message.NewChatMembers)
		}

		if update.Message.IsCommand() {
			handleCommands(bot, update)
		}
	}
}

func greetNewUser(bot *tgbotapi.BotAPI, userID int64, newUsers []tgbotapi.User) {
	for _, newUser := range newUsers {
		msg := tgbotapi.NewMessage(userID, "Hello "+newUser.FirstName+"! Welcome to the group!")
		bot.Send(msg)
	}
}
func handleCommands(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	command := update.Message.Command()

	switch command {
	case "start":
		sendWelcomeMessage(bot, update.Message.Chat.ID)
	case "resources":
		sendResources(bot, update.Message.Chat.ID)
	case "addinfo":
		handleAddInfoCommand(bot, update.Message.Chat.ID, update.Message.From.ID)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Type /start for a welcome message.")
		bot.Send(msg)
	}
}

func sendWelcomeMessage(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Welcome to the group! Type /resources to get information about Go clean code.")
	bot.Send(msg)
}

func sendResources(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Here are some resources about Go clean code:\n"+
		"- Clean Code in Go: https://github.com/Pungyeon/clean-go-article\n"+
		"- Effective Go: https://golang.org/doc/effective_go.html")
	bot.Send(msg)
}

func handleAddInfoCommand(bot *tgbotapi.BotAPI, chatID, userID int64) {
	msg := tgbotapi.NewMessage(chatID, "Please provide your GitHub username, LinkedIn URL, and phone number separated by commas.")
	bot.Send(msg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.From.ID == userID {
			addInfoMessage := update.Message
			info := strings.Split(addInfoMessage.Text, ",")

			if len(info) != 3 {
				bot.Send(tgbotapi.NewMessage(chatID, "Invalid input. Please provide all three details separated by commas."))
				return
			}

			membersDetails[userID] = MemberDetails{
				GitHubUsername: strings.TrimSpace(info[0]),
				LinkedInURL:    strings.TrimSpace(info[1]),
				PhoneNumber:    strings.TrimSpace(info[2]),
			}

			bot.Send(tgbotapi.NewMessage(chatID, "Information added successfully!"))
		}
	}
}
