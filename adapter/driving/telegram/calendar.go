package telegram

import (
	"fmt"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildCalendar(year int, month time.Month) tgbot.InlineKeyboardMarkup {
	var rows [][]tgbot.InlineKeyboardButton

	// Заголовок месяца
	rows = append(rows,
		tgbot.NewInlineKeyboardRow(
			tgbot.NewInlineKeyboardButtonData(fmt.Sprintf("📅 %s %d", month, year), "noop"),
		),
	)
	// Дни недели
	weekdays := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	var weekHeader []tgbot.InlineKeyboardButton
	for _, wd := range weekdays {
		weekHeader = append(weekHeader, tgbot.NewInlineKeyboardButtonData(wd, "noop"))
	}
	rows = append(rows, weekHeader)

	// Первый день месяца и смещение
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	weekday := int(firstDay.Weekday())
	if weekday == 0 {
		weekday = 7 // воскресенье → 7
	}

	// Дни месяца
	daysInMonth := firstDay.AddDate(0, 1, -1).Day()
	currentRow := make([]tgbot.InlineKeyboardButton, weekday-1) // сразу заполняем пустыми
	for i := range currentRow {
		currentRow[i] = tgbot.NewInlineKeyboardButtonData(" ", "noop")
	}

	for d := 1; d <= daysInMonth; d++ {
		dateStr := fmt.Sprintf("calendar:%d-%02d-%02d", year, month, d)
		currentRow = append(currentRow, tgbot.NewInlineKeyboardButtonData(fmt.Sprintf("%02d", d), dateStr))

		if len(currentRow) == 7 {
			rows = append(rows, currentRow)
			currentRow = nil
		}
	}
	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	// Навигация по месяцам
	prevYear, prevMonth := year, month-1
	if prevMonth < time.January {
		prevMonth, prevYear = time.December, year-1
	}
	nextYear, nextMonth := year, month+1
	if nextMonth > time.December {
		nextMonth, nextYear = time.January, year+1
	}

	rows = append(rows,
		tgbot.NewInlineKeyboardRow(
			tgbot.NewInlineKeyboardButtonData("←", fmt.Sprintf("calendar_prev:%d-%02d", prevYear, prevMonth)),
			tgbot.NewInlineKeyboardButtonData(" ", "noop"),
			tgbot.NewInlineKeyboardButtonData("→", fmt.Sprintf("calendar_next:%d-%02d", nextYear, nextMonth)),
		),
	)

	return tgbot.NewInlineKeyboardMarkup(rows...)
}
func (r *Router) SendCalendar(chatID int64, year int, month time.Month) {
	keyboard := BuildCalendar(year, month)
	msg := tgbot.NewMessage(chatID, "📅 Выберите дату:")
	msg.ReplyMarkup = keyboard
	r.bot.Send(msg)
}
