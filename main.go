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
				msg := tgbotapi.NewMessage(chatID, `🔮 Ты всё ещё думаешь, что крутые Reels — это только для профи с дорогой техникой?
👉 А что если я скажу, что ты можешь создавать WOW-видео с нуля, буквально по тексту или фото — с помощью нейросетей? Даже если ты полный новичок.

🎬 Представь:
— ты загружаешь фото — и через пару минут у тебя готово стильное видео
— пишешь пару строк — и нейросеть превращает это в динамичный ролик
— заменяешь фон, одежду, даже лицо в кадре — и всё это на телефоне

🔥 На курсе «Reels с помощью нейросетей, даже если ты новичок» ты освоишь навыки, которые прокачают тебя до уровня продюсера контента нового поколения:

✅ Создавать видео по тексту и фото — без съёмки
✅ Работать с 3D-объектами и внедрять их в Reels
✅ Использовать DeepFake-технологии для эффектной замены лиц
✅ Генерировать предметы, реквизит, декорации и одежду
✅ Делать Reels для брендов и продавать свои навыки
✅ Стилизовать видео под любой жанр или настроение — от кино до fashion
✅ И главное — освоишь базу, нужную для уверенной работы с нейросетями и мобильными приложениями

📲 Никакой сложной графики. Только конкретные инструменты, готовые сценарии и разбор твоих роликов.

🌟 Это обучение — твой шаг в будущее, где технологии работают на тебя, а ты — создаёшь вирусный, креативный и монетизируемый контент.`)

				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки текста об обучении:", err)
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
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(button),
	)
}
