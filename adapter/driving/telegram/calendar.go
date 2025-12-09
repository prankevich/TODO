package telegram

import (
	"fmt"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildCalendar(year int, month time.Month) tgbot.InlineKeyboardMarkup {
	var rows [][]tgbot.InlineKeyboardButton

	// Заголовок месяца
	monthTitle := fmt.Sprintf("📅 %s %d", month, year)
	rows = append(rows, tgbot.NewInlineKeyboardRow(tgbot.NewInlineKeyboardButtonData(monthTitle, "noop")))

	// Дни недели
	weekdays := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	var weekHeader []tgbot.InlineKeyboardButton
	for _, wd := range weekdays {
		weekHeader = append(weekHeader, tgbot.NewInlineKeyboardButtonData(wd, "noop"))
	}
	rows = append(rows, weekHeader)

	// Первый день месяца
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	weekday := int(firstDay.Weekday())
	if weekday == 0 {
		weekday = 7 // воскресенье → 7
	}

	// Пустые ячейки до первого дня
	var currentRow []tgbot.InlineKeyboardButton
	for i := 1; i < weekday; i++ {
		currentRow = append(currentRow, tgbot.NewInlineKeyboardButtonData(" ", "noop"))
	}

	// Дни месяца
	daysInMonth := firstDay.AddDate(0, 1, -1).Day()
	for d := 1; d <= daysInMonth; d++ {
		dateStr := fmt.Sprintf("calendar:%d-%02d-%02d", year, month, d)
		btn := tgbot.NewInlineKeyboardButtonData(fmt.Sprintf("%02d", d), dateStr)
		currentRow = append(currentRow, btn)

		if len(currentRow) == 7 {
			rows = append(rows, currentRow)
			currentRow = []tgbot.InlineKeyboardButton{}
		}
	}
	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	// Навигация по месяцам
	prevMonth := month - 1
	nextMonth := month + 1
	prevYear := year
	nextYear := year

	if prevMonth < time.January {
		prevMonth = time.December
		prevYear--
	}
	if nextMonth > time.December {
		nextMonth = time.January
		nextYear++
	}

	navRow := tgbot.NewInlineKeyboardRow(
		tgbot.NewInlineKeyboardButtonData("←", fmt.Sprintf("calendar_prev:%d-%02d", prevYear, prevMonth)),
		tgbot.NewInlineKeyboardButtonData(" ", "noop"),
		tgbot.NewInlineKeyboardButtonData("→", fmt.Sprintf("calendar_next:%d-%02d", nextYear, nextMonth)),
	)
	rows = append(rows, navRow)

	return tgbot.NewInlineKeyboardMarkup(rows...)
}

func (r *Router) SendCalendar(chatID int64, year int, month time.Month) {
	keyboard := BuildCalendar(year, month)
	msg := tgbot.NewMessage(chatID, "📅 Выберите дату:")
	msg.ReplyMarkup = keyboard
	r.bot.Send(msg)
}
