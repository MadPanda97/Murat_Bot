package main

import (
	"log"
	"os"
	"strings"

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
		if update.Message != nil {
			text := strings.ToLower(update.Message.Text)

			if strings.Contains(text, "хочу урок") {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Вот бесплатный видеоурок 👇")
				msg.ReplyMarkup = firstInlineKeyboard()
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки по ключевому слову:", err)
				}
			}

			if update.Message.Video != nil {
				log.Printf("Video FileID: %s", update.Message.Video.FileID)
			}
		}

		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				log.Println("Ошибка подтверждения callback:", err)
			}

			data := update.CallbackQuery.Data
			chatID := update.CallbackQuery.Message.Chat.ID

			switch data {
			case "free_lesson":
				video := tgbotapi.NewVideo(chatID, tgbotapi.FileID("BAACAgIAAxkBAAMPaDQXibPXuGm8U9wG_KjFDwrJ8JkAAlVqAAJOBKFJ0BS7KrQcUS82BA"))
				video.Caption = "Вот бесплатный видеоурок!\n\nХочешь узнать, какие нейросети я использую в Reels?"
				video.ReplyMarkup = secondInlineKeyboard()
				if _, err := bot.Send(video); err != nil {
					log.Println("Ошибка отправки видео:", err)
				}

			case "ai_tools":
				doc := tgbotapi.NewDocument(chatID, tgbotapi.FileID("BQACAgIAAxkBAAMUaDQYsojlC47_ygUxnhYkdZGrCEwAAoBqAAJOBKFJGvpBU-vHqYo2BA"))
				doc.Caption = "Вот чек-лист нейросетей, которые я использую 💡"
				doc.ReplyMarkup = thirdInlineKeyboard()
				if _, err := bot.Send(doc); err != nil {
					log.Printf("Ошибка отправки PDF: %v", err)
				}

			case "watch_reels":
				video := tgbotapi.NewVideo(chatID, tgbotapi.FileID("BAACAgIAAxkBAAMWaDVurVxdzzx3_2JS8c0DTXLsrWMAAhxvAAIYyaFJaxtuxnFLdGE2BA"))
				video.Caption = "Вот примеры Reels, которые я создаю с помощью нейросетей 🔥"
				video.ReplyMarkup = finalInlineKeyboard()
				if _, err := bot.Send(video); err != nil {
					log.Printf("Ошибка отправки Reels-видео: %v", err)
				}

			case "course_learn":
				msg := tgbotapi.NewMessage(chatID, `🚨 Закрываем набор через 2 дня — потом будет только следующий поток.
💰 Сейчас цена — 50 000 тг. После запуска поднимаем до 75 000 тг.

💰 Сейчас можно выбрать один из 2 тарифов:
✅ *С поддержкой автора* — обратная связь, помощь, ответы на вопросы
💼 *Классический* — просто уроки, без сопровождения

Выбери свой тариф:`)
				msg.ParseMode = "Markdown"
				msg.ReplyMarkup = tariffKeyboard()
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки тарифов:", err)
				}

			case "tariff_support":
				msg := tgbotapi.NewMessage(chatID, "🔗 Ссылка на тариф с поддержкой: https://your-link.com/support")
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки ссылки с поддержкой:", err)
				}

			case "tariff_classic":
				msg := tgbotapi.NewMessage(chatID, "🔗 Ссылка на классический тариф: https://your-link.com/classic")
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки ссылки на классический тариф:", err)
				}
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
	button := tgbotapi.NewInlineKeyboardButtonData("📚 Хочу обучение", "course_learn")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
}

func tariffKeyboard() tgbotapi.InlineKeyboardMarkup {
	btn1 := tgbotapi.NewInlineKeyboardButtonData("✅ С поддержкой", "tariff_support")
	btn2 := tgbotapi.NewInlineKeyboardButtonData("💼 Классический", "tariff_classic")
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(btn1),
		tgbotapi.NewInlineKeyboardRow(btn2),
	)
}
