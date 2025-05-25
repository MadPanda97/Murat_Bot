package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	botToken := os.Getenv("TELEGRAM_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Бот %s запущен", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Video != nil {
			log.Printf("Video FileID: %s", update.Message.Video.FileID)
		}

		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Начнем?")
				msg.ReplyMarkup = firstInlineKeyboard()
				bot.Send(msg)
			}
		}

		if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data
			chatID := update.CallbackQuery.Message.Chat.ID

			switch data {
			case "free_lesson":
				video := tgbotapi.NewVideo(chatID, tgbotapi.FileID("ВСТАВЬ_СЮДА_file_id"))
				video.Caption = "Вот бесплатный видеоурок!\n\nХочешь узнать, какие нейросети я использую в Reels?"
				video.ReplyMarkup = secondInlineKeyboard()
				bot.Send(video)

			case "ai_tools":
				msg := tgbotapi.NewMessage(chatID, "Вот чек-лист нейросетей, которые я использую.")
				msg.ReplyMarkup = thirdInlineKeyboard()
				bot.Send(msg)

			case "watch_reels":
				msg := tgbotapi.NewMessage(chatID, "Посмотри, какие Reels я создаю с этими нейросетями.")
				msg.ReplyMarkup = finalInlineKeyboard()
				bot.Send(msg)

			case "course_offer":
				msg := tgbotapi.NewMessage(chatID, "Ты сможешь научиться создавать такие видео за пару недель и зарабатывать на этом. Курс включает:\n\n• Видео из текста\n• 3D-объекты\n• DeepFake\n• Видео для брендов\n• Генерация предметов\n\nУроки доступны 24/7 на 6 месяцев.")
				bot.Send(msg)
			}
		}
	}
}

func firstInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("🎥 Бесплатный видеоурок", "free_lesson")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
}

func secondInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("🤖 Нейросети для контента", "ai_tools")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
}

func thirdInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("📱 Смотреть Reels", "watch_reels")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
}

func finalInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("🚀 Перейти к обучению", "course_offer")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
}
