package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
	"time"
)

type UserInfo struct {
	State        string
	RegisteredAt time.Time
	Paid         bool
	Name         string
}

var (
	users      = make(map[int64]*UserInfo)
	usersMutex sync.Mutex
)

const (
	StateStart             = "start"
	StatePDF               = "pdf"
	StateBreakfastInfo     = "breakfast_info"
	StatePayment           = "payment"
	StateWaitingScreenshot = "waiting_screenshot"
	StatePaid              = "paid"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл:", err)
	}

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

	go sendReminders(bot)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID

			if update.Message.IsCommand() && update.Message.Command() == "start" {
				usersMutex.Lock()
				users[chatID] = &UserInfo{
					State:        StateStart,
					RegisteredAt: time.Now(),
					Paid:         false,
				}
				usersMutex.Unlock()

				msg := tgbotapi.NewMessage(chatID, "Привет! Рад видеть тебя здесь.\n🎯 Хочешь узнать, как делать Reels с помощью нейросетей — быстро, стильно?")
				msg.ReplyMarkup = welcomeKeyboard()
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки приветственного сообщения:", err)
				}
			} else if update.Message.Photo != nil {
				usersMutex.Lock()
				userInfo, exists := users[chatID]
				usersMutex.Unlock()

				if exists && userInfo.State == StateWaitingScreenshot {
					usersMutex.Lock()
					userInfo.State = StatePaid
					userInfo.Paid = true
					usersMutex.Unlock()

					msg := tgbotapi.NewMessage(chatID, "✅ Спасибо за скриншот! Я проверю оплату и подтвержу твоё участие.\n\n📍 Адрес проведения бизнес-завтрака будет отправлен за день до мероприятия.\n\nДо встречи 28 июня! ☕️")
					if _, err := bot.Send(msg); err != nil {
						log.Println("Ошибка отправки подтверждения получения скриншота:", err)
					}
				}
			}
		}

		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				log.Println("Ошибка подтверждения callback:", err)
			}

			chatID := update.CallbackQuery.Message.Chat.ID
			data := update.CallbackQuery.Data

			switch data {
			case "get_pdf":
				usersMutex.Lock()
				if userInfo, exists := users[chatID]; exists {
					userInfo.State = StatePDF
				}
				usersMutex.Unlock()

				doc := tgbotapi.NewDocument(chatID, tgbotapi.FileID("BQACAgIAAxkBAAMUaDQYsojlC47_ygUxnhYkdZGrCEwAAoBqAAJOBKFJGvpBU-vHqYo2BA"))
				doc.Caption = "📘 Вот твой PDF-гайд:\n5 нейросетей, которые делают Reels за тебя\n\n❗️ А если хочешь вживую увидеть, как они работают, приходи на бизнес-завтрак!"
				doc.ReplyMarkup = breakfastInfoKeyboard()
				if _, err := bot.Send(doc); err != nil {
					log.Printf("Ошибка отправки PDF: %v", err)
				}

			case "breakfast_info":
				usersMutex.Lock()
				if userInfo, exists := users[chatID]; exists {
					userInfo.State = StateBreakfastInfo
				}
				usersMutex.Unlock()

				msg := tgbotapi.NewMessage(chatID, `📍 Бизнес-завтрак в Семее
🗓 28 июня | 🕐 13:00
📍 Уютная кофейня (точный адрес после оплаты)
💸 Участие: 4 900 ₸
👥 Только 12 мест

🎁 Что ты получишь:
– PDF + демонстрация нейросетей
– Reels, которые продают
– Нетворкинг и мои фишки, которые никто не показывает онлайн

👇 Чтобы зафиксировать место, внеси предоплату.`)
				msg.ReplyMarkup = paymentKeyboard()
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки информации о бизнес-завтраке:", err)
				}

			case "payment":
				usersMutex.Lock()
				if userInfo, exists := users[chatID]; exists {
					userInfo.State = StatePayment
				}
				usersMutex.Unlock()

				photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileID("AgACAgIAAxkBAAIBdGXLXKNXJAABdGXLXKNXJAABdGXLXKNXJ")) // Замените на реальный FileID QR-кода
				photo.Caption = "📱 Отсканируй QR-код для быстрой оплаты"
				if _, err := bot.Send(photo); err != nil {
					log.Println("Ошибка отправки QR-кода:", err)
				}

				msg := tgbotapi.NewMessage(chatID, `💳 Оплата:
Переведи 4 900 ₸ на Kaspi:
https://pay.kaspi.kz/pay/s2ompo9l

📝 В комментарии укажи:
"ЗАВТРАК + Имя"

После оплаты — нажми кнопку ниже и пришли скрин перевода.`)
				msg.ReplyMarkup = sendScreenshotKeyboard()
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки информации об оплате:", err)
				}

			case "send_screenshot":
				// Обновляем состояние пользователя
				usersMutex.Lock()
				if userInfo, exists := users[chatID]; exists {
					userInfo.State = StateWaitingScreenshot
				}
				usersMutex.Unlock()

				// Просим отправить скриншот оплаты (ЭКРАН 5)
				msg := tgbotapi.NewMessage(chatID, `📥 Пришли, пожалуйста, скрин оплаты (фото Kaspi перевода)

✅ Как только я увижу оплату — ты в списке участников!
📍 И вышлю адрес за 1 день до мероприятия.

Спасибо за доверие! До встречи 28 июня ☕️`)
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки запроса на скриншот:", err)
				}
			}
		}
	}
}

// Функция для отправки напоминаний за день до мероприятия
func sendReminders(bot *tgbotapi.BotAPI) {
	for {
		// Проверяем каждый час
		time.Sleep(1 * time.Hour)

		// Дата мероприятия
		eventDate := time.Date(2023, time.June, 28, 13, 0, 0, 0, time.Local)
		now := time.Now()

		// Если до мероприятия осталось 1 день
		if eventDate.Sub(now) <= 24*time.Hour && eventDate.Sub(now) > 0 {
			usersMutex.Lock()
			for chatID, userInfo := range users {
				if userInfo.Paid {
					msg := tgbotapi.NewMessage(chatID, "Привет! Напоминаю: завтра встречаемся в 13:00, Семей. Адрес: [тут]")
					if _, err := bot.Send(msg); err != nil {
						log.Println("Ошибка отправки напоминания:", err)
					}
				}
			}
			usersMutex.Unlock()
		}
	}
}

func welcomeKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("📘 Получить PDF + Подробнее", "get_pdf")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
}

func breakfastInfoKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("Узнать о завтраке", "breakfast_info")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
}

func paymentKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("Оплатить участие", "payment")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
}

func sendScreenshotKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("Отправить скриншот", "send_screenshot")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
}
