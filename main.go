package main

import (
	"github.com/joho/godotenv"
	"log"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
)

var (
	waitingUsers   = make(map[int64]time.Time) // Изменено: теперь храним время регистрации
	waitingUsersMu sync.Mutex
)

func main() {
	botToken := os.Getenv("TELEGRAM_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_TOKEN не установлен")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Бот %s запущен", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go reminderLoop(bot)

	for update := range updates {
		if update.Message != nil {
			text := strings.ToLower(update.Message.Text)
			chatID := update.Message.Chat.ID

			switch {
			case update.Message.IsCommand() && update.Message.Command() == "start":
				// Обработка команды /start - просто приветствие без видео
				msg := tgbotapi.NewMessage(chatID, "Привет! Я бот для обучения созданию Reels с помощью нейросетей.")
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки приветственного сообщения:", err)
				}

			case text == "хочу на курс":
				waitingUsersMu.Lock()
				waitingUsers[chatID] = time.Now()
				waitingUsersMu.Unlock()

				msg := tgbotapi.NewMessage(chatID, "Отлично! Ты записан на обучение. Мы напомним тебе через 1 и 2 дня.")
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки подтверждения записи на курс:", err)
				}

			case text == "хочу урок":
				// Отправляем первое видео при получении текста "хочу урок"
				video := tgbotapi.NewVideo(chatID, tgbotapi.FileID("BAACAgIAAxkBAAMPaDQXibPXuGm8U9wG_KjFDwrJ8JkAAlVqAAJOBKFJ0BS7KrQcUS82BA"))
				video.Caption = "Вот бесплатный видеоурок!\n\nХочешь узнать, какие нейросети я использую в Reels?"
				video.ReplyMarkup = secondInlineKeyboard()
				if _, err := bot.Send(video); err != nil {
					log.Println("Ошибка отправки первого видео:", err)
				}

			case text == "урок":
				// Отправляем второе видео при получении текста "урок"
				video := tgbotapi.NewVideo(chatID, tgbotapi.FileID("BAACAgIAAxkBAAOqaEKaNA_S86x5zT0x9wu1Ot75Be8AAqB2AAIuOQhKZ47zBXBvHLU2BA"))
				video.Caption = "Вот еще один видеоурок!\n\nХочешь узнать, какие нейросети я использую в Reels?"
				video.ReplyMarkup = secondInlineKeyboard()
				if _, err := bot.Send(video); err != nil {
					log.Println("Ошибка отправки второго видео:", err)
				}

			default:
				if update.Message.Video != nil {
					log.Printf("Video FileID: %s", update.Message.Video.FileID)
				}
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
				} else {
					// Отправляем второе сообщение с кнопкой выбора тарифов
					tariffMsg := tgbotapi.NewMessage(chatID, "Выбери подходящий тариф обучения:")
					tariffMsg.ReplyMarkup = tariffKeyboard()
					if _, err := bot.Send(tariffMsg); err != nil {
						log.Println("Ошибка отправки выбора тарифа после информации о курсе:", err)
					}
				}

			case "choose_tariff":
				msg := tgbotapi.NewMessage(chatID, "Выбери тариф:")
				msg.ReplyMarkup = tariffKeyboard()
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки выбора тарифа:", err)
				}

			case "tariff_support":
				msg := tgbotapi.NewMessage(chatID, "Ты выбрал тариф с поддержкой от автора. Вот ссылка для оплаты:\nhttps://murat.courstore.com/ru/courses/copy-524189-2")
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки сообщения тарифа поддержки:", err)
				}

			case "tariff_classic":
				msg := tgbotapi.NewMessage(chatID, "Ты выбрал классический тариф. Вот ссылка для оплаты:\nhttps://murat.courstore.com/ru/courses/524189")
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки сообщения классического тарифа:", err)
				}
			}
		}
	}
}

func reminderLoop(bot *tgbotapi.BotAPI) {
	firstReminderSent := make(map[int64]bool)
	secondReminderSent := make(map[int64]bool)

	for {
		time.Sleep(1 * time.Hour) // Проверяем каждый час

		now := time.Now()
		waitingUsersMu.Lock()

		for chatID, registrationTime := range waitingUsers {
			// Проверяем, прошло ли 24 часа с момента регистрации
			if !firstReminderSent[chatID] && now.Sub(registrationTime) >= 24*time.Hour {
				msg := tgbotapi.NewMessage(chatID, "Привет! Напоминаем, что ты записался на курс. Если хочешь, напиши 'ХОЧУ НА КУРС' для подтверждения.")
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки первого напоминания:", err)
				} else {
					firstReminderSent[chatID] = true
					log.Printf("Отправлено первое напоминание пользователю %d", chatID)
				}
			}

			// Проверяем, прошло ли 48 часов с момента регистрации
			if !secondReminderSent[chatID] && now.Sub(registrationTime) >= 48*time.Hour {
				msg := tgbotapi.NewMessage(chatID, "Это финальное напоминание! Если есть вопросы, пиши, я помогу.")
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки второго напоминания:", err)
				} else {
					secondReminderSent[chatID] = true
					log.Printf("Отправлено второе напоминание пользователю %d", chatID)
				}
			}
		}

		waitingUsersMu.Unlock()
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
	buttonSupport := tgbotapi.NewInlineKeyboardButtonData("✅С поддержкой от автора", "tariff_support")
	buttonClassic := tgbotapi.NewInlineKeyboardButtonData("💼Классический тариф", "tariff_classic")
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttonSupport),
		tgbotapi.NewInlineKeyboardRow(buttonClassic),
	)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл:", err)
	}
}
