package main

import (
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type User struct {
	ID             int64
	StartTime      time.Time
	SentAfter1Day  bool
	SentAfter2Days bool
}

var users = make(map[int64]*User)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN не установлен")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	startFollowUpChecker(bot)

	for update := range updates {
		if update.Message != nil {
			userID := update.Message.Chat.ID
			text := update.Message.Text

			if text == "/start" {
				msg := tgbotapi.NewMessage(userID, "Привет! Готов погрузиться в мир Reels?")
				msg.ReplyMarkup = mainMenu()
				bot.Send(msg)
			} else if text == "Хочу обучение" {
				users[userID] = &User{ID: userID, StartTime: time.Now()}

				msg := tgbotapi.NewMessage(userID, "Ты всё ещё думаешь, что крутые Reels — это только для профи с дорогой техникой?")
				msg.ReplyMarkup = tariffKeyboard()
				bot.Send(msg)
			}
		}

		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			switch callback.Data {
			case "select_tariff":
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Выбери подходящий тариф:")
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL("Тариф с поддержкой", "https://example.com/support"),
						tgbotapi.NewInlineKeyboardButtonURL("Классический тариф", "https://example.com/basic"),
					),
				)
				bot.Send(msg)
			}
		}
	}
}

func mainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Хочу обучение"),
		),
	)
}

func tariffKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Выбрать тариф", "select_tariff"),
		),
	)
}

func startFollowUpChecker(bot *tgbotapi.BotAPI) {
	go func() {
		for {
			for _, user := range users {
				since := time.Since(user.StartTime)

				if !user.SentAfter1Day && since > 24*time.Hour {
					msg := tgbotapi.NewMessage(user.ID, "Ты так и не написал… Бронь слетит через сутки 😔")
					bot.Send(msg)
					user.SentAfter1Day = true
				}

				if !user.SentAfter2Days && since > 48*time.Hour {
					msg := tgbotapi.NewMessage(user.ID, "Если остались вопросы — пиши, помогу выбрать 🧠")
					bot.Send(msg)
					user.SentAfter2Days = true
				}
			}
			time.Sleep(10 * time.Minute)
		}
	}()
}
