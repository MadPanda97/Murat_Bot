package main

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
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
	StateStart      = "start"
	StatePDF        = "pdf"
	StateCourseInfo = "course_info"
	StatePayment    = "payment"
	StatePaid       = "paid"
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

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			log.Printf("Получено сообщение от пользователя %d: %s", chatID, update.Message.Text)

			// Проверяем команду /start или ключевое слово "хочу"
			if (update.Message.IsCommand() && update.Message.Command() == "start") ||
				(update.Message.Text != "" && strings.ToLower(update.Message.Text) == "хочу") {

				usersMutex.Lock()
				users[chatID] = &UserInfo{
					State:        StateStart,
					RegisteredAt: time.Now(),
					Paid:         false,
				}
				usersMutex.Unlock()

				log.Printf("Отправляю приветственное сообщение пользователю %d", chatID)
				msg := tgbotapi.NewMessage(chatID, "Привет! Рад видеть тебя здесь.\n🎯 Хочешь узнать, как делать Reels с помощью нейросетей — быстро, стильно?")
				msg.ReplyMarkup = welcomeKeyboard()
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Ошибка отправки приветственного сообщения пользователю %d: %v", chatID, err)
				} else {
					log.Printf("Приветственное сообщение успешно отправлено пользователю %d", chatID)
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
			log.Printf("Получен callback: %s от пользователя %d", data, chatID)

			switch data {
			case "get_pdf":
				usersMutex.Lock()
				if userInfo, exists := users[chatID]; exists {
					userInfo.State = StatePDF
				}
				usersMutex.Unlock()

				// Отправляем PDF
				log.Printf("Отправляю PDF пользователю %d", chatID)
				doc := tgbotapi.NewDocument(chatID, tgbotapi.FileID("BQACAgIAAxkBAAMUaDQYsojlC47_ygUxnhYkdZGrCEwAAoBqAAJOBKFJGvpBU-vHqYo2BA"))
				doc.Caption = "📘 Вот твой PDF-гайд:\n5 нейросетей, которые делают Reels за тебя"
				if _, err := bot.Send(doc); err != nil {
					log.Printf("Ошибка отправки PDF: %v", err)
					return // Если PDF не отправился, не продолжаем
				}
				log.Printf("PDF успешно отправлен пользователю %d", chatID)

				// Небольшая задержка между отправкой PDF и видео
				time.Sleep(1 * time.Second)

				// Отправляем видео с кнопкой
				log.Printf("Отправляю видео пользователю %d", chatID)
				video := tgbotapi.NewVideo(chatID, tgbotapi.FileID("BAACAgIAAxkBAAIBrWjSY0Ey6QwZ3wI_GjPkhXowkRF8AALZhAACmzxRSqv6032r5Wm3NgQ"))
				video.Caption = "🎬 Урок второй воронка - пример создания видео с ИИ\n\n❗️ А если хочешь освоить создание коммерческих видео с ИИ за 4 недели, нажми кнопку ниже!"
				video.ReplyMarkup = courseInfoKeyboard()
				
				// Попробуем отправить видео с повторными попытками
				maxRetries := 3
				for i := 0; i < maxRetries; i++ {
					if _, err := bot.Send(video); err != nil {
						log.Printf("Попытка %d/%d - Ошибка отправки видео: %v", i+1, maxRetries, err)
						if i == maxRetries-1 {
							log.Printf("Не удалось отправить видео пользователю %d после %d попыток", chatID, maxRetries)
						} else {
							time.Sleep(2 * time.Second)
						}
					} else {
						log.Printf("Видео успешно отправлено пользователю %d (попытка %d)", chatID, i+1)
						break
					}
				}

			case "course_info":
				usersMutex.Lock()
				if userInfo, exists := users[chatID]; exists {
					userInfo.State = StateCourseInfo
				}
				usersMutex.Unlock()

				msg := tgbotapi.NewMessage(chatID, `🎬 Создавайте коммерческие видео с ИИ за 4 недели

🚀 Без студии, не нужен мощный телефон или компьютер.

Освойте инструменты AI, соберите портфолио и подготовьтесь к первым заказам. Формат: записи 24/7, доступ 6 месяцев.

💡 За месяц вы научитесь создавать контент с помощью AI.

✨ Вы будете выгодно отличаться от множества одинаковых работ: вместо типовых роликов — заметные видео под задачи клиентов.

🎯 Дам рекомендации, где и как искать клиентов. Вы получите доступ к чатам, где публикуются заказы.

📋 Главное — вы перестанете действовать бессистемно: начнете работать по четкому пайплайну, соберите портфолио и подготовитесь к первым заказам.

🏠 Работать можно из дома и часто без съемок: достаточно телефона или ноутбука.

👥 Для кого:
• Новички в видео и SMM, которым нужна быстрая, прикладная база по AI-видео
• Создатели контента, желающие повысить чек и ускорить продакшн
• Фрилансеры, которые хотят находить клиентов и работать из дома

🎁 Что вы получите за месяц:
• Видео для клиентов с помощью ИИ для портфолио: из текста, из фото, для брендов
• Понимание пайплайна: от идеи и референсов до финального видео
• Рекомендации по поиску клиентов и доступ в чаты с заказами`)

				// Отправляем первое сообщение
				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки информации о курсе:", err)
				}

				// Отправляем второе сообщение с программой курса
				msg2 := tgbotapi.NewMessage(chatID, `📚 Программа курса (модули):
• Работа с нейросетями: Runway, Pika, CapCut, Midjourney, Kling, Nano Banana и др
• Подготовительный модуль: основы ИИ для видео
• Видео из текста и фото: генерация, стили, референсы, монтаж и звук
• 3D и интеграции: создание/импорт 3D‑объектов, внедрение в кадр
• Бренд‑видео с нейросетями: работа под задачу клиента, адаптация под форматы Reels/Shorts/TikTok
• Генерация предметов и сцен: декорации, реквизит, смена пространства и одежды
• Коммерция: упаковка портфолио, чек‑лист качества, поиск клиентов, коммуникация и чек

📖 Формат и поддержка:
• Уроки в записи, доступ 24/7 на 6 месяцев
• Разбитие на модули с практикой
• Комьюнити и чаты с заказами
• Обратная связь по домашним заданиям
• В ближайшие дни добавим новые актуальные уроки

🏆 Результат:
• Собранное мини‑портфолио
• Понимание, где и как искать первые заказы
• Готовые шаблоны и чек‑листы, чтобы работать быстрее и дороже`)

				msg2.ReplyMarkup = paymentKeyboard()
				if _, err := bot.Send(msg2); err != nil {
					log.Println("Ошибка отправки программы курса:", err)
				}

			case "payment":
				usersMutex.Lock()
				if userInfo, exists := users[chatID]; exists {
					userInfo.State = StatePayment
				}
				usersMutex.Unlock()

				msg := tgbotapi.NewMessage(chatID, `💳 Записаться на курс "AI-видео за 4 недели"

📞 Для записи свяжитесь с менеджером:
@Murat_76video

💰 Стоимость курса и условия оплаты уточняйте у менеджера.

🎁 При записи через бота — скидка 10%!

✅ После оплаты вы получите:
• Доступ ко всем урокам на 6 месяцев
• Вход в закрытое комьюнити
• Доступ к чатам с заказами
• Обратную связь по домашним заданиям`)

				if _, err := bot.Send(msg); err != nil {
					log.Println("Ошибка отправки информации об оплате:", err)
				}
			}
		}
	}
}

func welcomeKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("📘 Получить PDF + Подробнее", "get_pdf")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
}

func courseInfoKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("Подробнее", "course_info")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
}

func paymentKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("Записаться на курс", "payment")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
}
