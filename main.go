package main

import (
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type UserState struct {
	Stage      int
	LastAction time.Time
}

var userStates = make(map[int64]*UserState)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go reminderScheduler(bot)

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		chatID := update.FromChat().ID
		state, exists := userStates[chatID]
		if !exists {
			state = &UserState{Stage: 0, LastAction: time.Now()}
			userStates[chatID] = state
		}

		if update.Message != nil && update.Message.Text == "Хочу урок" {
			state.Stage = 1
			state.LastAction = time.Now()
			sendStage1(bot, chatID)
			continue
		}

		if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data
			switch data {
			case "checklist":
				state.Stage = 2
				state.LastAction = time.Now()
				sendStage2(bot, chatID)
			case "reels":
				state.Stage = 3
				state.LastAction = time.Now()
				sendStage3(bot, chatID)
			case "course":
				state.Stage = 4
				state.LastAction = time.Now()
				sendStage4(bot, chatID)
			}
		}
	}
}

func sendStage1(bot *tgbotapi.BotAPI, chatID int64) {
	video := tgbotapi.NewVideo(chatID, tgbotapi.FileID("BAACAgIAAxkBAAMPaDQXibPXuGm8U9wG_KjFDwrJ8JkAAlVqAAJOBKFJ0BS7KrQcUS82BA"))
	video.Caption = "Вот бесплатный видеоурок 🎓"
	bot.Send(video)

	msg := tgbotapi.NewMessage(chatID, "Хочешь чек-лист по нейросетям для Reels?")
	btn := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да, хочу", "checklist"),
		),
	)
	msg.ReplyMarkup = btn
	bot.Send(msg)
}

func sendStage2(bot *tgbotapi.BotAPI, chatID int64) {
	doc := tgbotapi.NewDocument(chatID, tgbotapi.FileID("BQACAgIAAxkBAAMUaDQYsojlC47_ygUxnhYkdZGrCEwAAoBqAAJOBKFJGvpBU-vHqYo2BA"))
	doc.Caption = "Вот твой чек-лист ✅"
	bot.Send(doc)

	msg := tgbotapi.NewMessage(chatID, "Хочешь примеры крутых Reels?")
	btn := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да, покажи", "reels"),
		),
	)
	msg.ReplyMarkup = btn
	bot.Send(msg)
}

func sendStage3(bot *tgbotapi.BotAPI, chatID int64) {
	video := tgbotapi.NewVideo(chatID, tgbotapi.FileID("BAACAgIAAxkBAAMWaDVurVxdzzx3_2JS8c0DTXLsrWMAAhxvAAIYyaFJaxtuxnFLdGE2BA"))
	video.Caption = "Ты всё ещё думаешь, что крутые Reels — это только для профи с дорогой техникой?"
	bot.Send(video)

	msg := tgbotapi.NewMessage(chatID, "🔮 Ты всё ещё думаешь, что крутые Reels — это только для профи с дорогой техникой?\n👉 А что если я скажу, что ты можешь создавать WOW-видео с нуля, буквально по тексту или фото — с помощью нейросетей? Даже если ты полный новичок.\n\n🎬 Представь:\n— ты загружаешь фото — и через пару минут у тебя готово стильное видео\n— пишешь пару строк — и нейросеть превращает это в динамичный ролик\n— заменяешь фон, одежду, даже лицо в кадре — и всё это на телефоне\n\n🔥 На курсе «Reels с помощью нейросетей, даже если ты новичок» ты освоишь навыки, которые прокачают тебя до уровня продюсера контента нового поколения:\n\n✅ Создавать видео по тексту и фото — без съёмки\n✅ Работать с 3D-объектами и внедрять их в Reels\n✅ Использовать DeepFake-технологии для эффектной замены лиц\n✅ Генерировать предметы, реквизит, декорации и одежду\n✅ Делать Reels для брендов и продавать свои навыки\n✅ Стилизовать видео под любой жанр или настроение — от кино до fashion\n✅ И главное — освоишь базу, нужную для уверенной работы с нейросетями и мобильными приложениями\n\n📲 Никакой сложной графики. Только конкретные инструменты, готовые сценарии и разбор твоих роликов.\n\n🌟 Это обучение — твой шаг в будущее, где технологии работают на тебя, а ты — создаёшь вирусный, креативный и монетизируемый контент.")
	btn := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Хочу курс", "course"),
		),
	)
	msg.ReplyMarkup = btn
	bot.Send(msg)
}

func sendStage4(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Хочешь на обучение? Напиши в чат «ХОЧУ НА КУРС» и я всё расскажу 🔥")
	bot.Send(msg)
}

// ---------------------------
// Напоминалки через 1 и 2 дня
// ---------------------------

func reminderScheduler(bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(time.Hour * 1)
	defer ticker.Stop()

	for {
		<-ticker.C
		now := time.Now()

		for chatID, state := range userStates {
			if state.Stage < 4 {
				since := now.Sub(state.LastAction)

				if since > 24*time.Hour && since < 25*time.Hour {
					sendReminder(bot, chatID, 1)
				} else if since > 48*time.Hour && since < 49*time.Hour {
					sendReminder(bot, chatID, 2)
				}
			}
		}
	}
}

func sendReminder(bot *tgbotapi.BotAPI, chatID int64, day int) {
	var text string
	if day == 1 {
		text = "Ты ещё думаешь? Урок ждёт тебя 👉"
	} else if day == 2 {
		text = "Если есть вопросы — пиши, я помогу 💬"
	}

	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}
