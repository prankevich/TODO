package telegram

import (
	"TODO/adapter/driven/models"
	"fmt"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (r *Router) SendTask(chatID int64, task models.Task) {
	status := "⏳"
	if task.Done {
		status = "✅"
	}
	text := fmt.Sprintf("%s %s — %s : %s Дата исполнения %s", status, task.ID, task.Title, task.Notes, task.DueAt)
	msg := tgbot.NewMessage(chatID, text)

	completeBtn := tgbot.NewInlineKeyboardButtonData("Завершить", "complete:"+task.ID)
	deleteBtn := tgbot.NewInlineKeyboardButtonData("🗑 Удалить", "delete:"+task.ID)
	updateBtn := tgbot.NewInlineKeyboardButtonData("✏️ Изменить", "update:"+task.ID)

	keyboard := tgbot.NewInlineKeyboardMarkup(
		tgbot.NewInlineKeyboardRow(completeBtn, deleteBtn, updateBtn),
	)

	msg.ReplyMarkup = keyboard
	r.bot.Send(msg)
}
func (r *Router) SendMainMenu(chatID int64) {
	addBtn := tgbot.NewInlineKeyboardButtonData("➕ Добавить", "menu:add")
	listBtn := tgbot.NewInlineKeyboardButtonData("📋Список", "menu:list")
	storyBtn := tgbot.NewInlineKeyboardButtonData("📋Список выполненных", "menu:listStory")
	keyboard := tgbot.NewInlineKeyboardMarkup(
		tgbot.NewInlineKeyboardRow(addBtn, listBtn, storyBtn),
	)

	msg := tgbot.NewMessage(chatID, "Выберите действие:")
	msg.ReplyMarkup = keyboard
	r.bot.Send(msg)
}
