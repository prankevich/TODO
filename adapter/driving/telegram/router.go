package telegram

import (
	"TODO/adapter/driven/models"
	"TODO/services/contracts"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userStates = make(map[int64]*UserState)

type Router struct {
	bot  *tgbot.BotAPI
	todo contracts.TodoRepository
}

type UserState struct {
	Step     int
	TempTask models.Task
}

func NewRouter(bot *tgbot.BotAPI, todo contracts.TodoRepository) *Router {
	return &Router{
		bot:  bot,
		todo: todo,
	}
}

func (r *Router) Run(ctx context.Context) error {
	u := tgbot.NewUpdate(0)
	u.Timeout = 30
	updates := r.bot.GetUpdatesChan(u)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case upd := <-updates:
			if upd.Message != nil {
				r.HandleMessage(ctx, upd.Message)
			}
			if upd.CallbackQuery != nil {
				r.HandleCallback(ctx, upd.CallbackQuery)
			}
		}
	}
}

func (r *Router) HandleMessage(ctx context.Context, msg *tgbot.Message) {
	if msg == nil || msg.From == nil {
		log.Println("Получено пустое сообщение")
		return
	}
	telegramID := msg.From.ID
	text := strings.TrimSpace(msg.Text)
	log.Printf("User %s (tg_%d) написал: %q", msg.From.UserName, telegramID, text)

	state, ok := userStates[telegramID]
	if !ok {
		r.SendMainMenu(msg.Chat.ID)
		return
	}

	switch state.Step {
	case 1, 11:
		date, err := time.Parse("02-01-2006", text)
		if err != nil {
			r.Reply(msg.Chat.ID, "❌ Неверный формат даты. Используйте DD-MM-YYYY или выберите из календаря ниже:")
			r.SendCalendar(msg.Chat.ID, time.Now().Year(), time.Now().Month())
			return
		}
		state.TempTask.DueAt = &date
		if state.Step == 1 {
			state.Step = 2
			r.Reply(msg.Chat.ID, fmt.Sprintf("Вы выбрали дату: %s", date.Format("02-01-2006")))
			r.Reply(msg.Chat.ID, "Введите заголовок задачи:")
		} else {
			r.finishUpdate(ctx, msg.Chat.ID, telegramID, state.TempTask)
		}

	case 2:
		if text == "" {
			r.Reply(msg.Chat.ID, "⚠️ Заголовок не может быть пустым.")
			return
		}
		state.TempTask.Title = text
		state.Step = 3
		r.Reply(msg.Chat.ID, "Введите описание задачи:")

	case 3:
		if text == "" {
			r.Reply(msg.Chat.ID, "⚠️ Описание не может быть пустым.")
			return
		}
		state.TempTask.Notes = text
		state.TempTask.CreatedAt = time.Now().UTC()

		task, err := r.todo.CreateTask(ctx, state.TempTask)
		if err != nil {
			log.Printf("Ошибка при добавлении задачи: %v", err)
			r.Reply(msg.Chat.ID, "❌ Ошибка при добавлении задачи")
		} else {
			r.Reply(msg.Chat.ID, fmt.Sprintf("✅ Задача \"%s\" добавлена на %s",
				task.Title, task.DueAt.Format("02-01-2006")))
		}

		delete(userStates, telegramID)
		r.SendMainMenu(msg.Chat.ID)

	case 12:
		if text == "" {
			r.Reply(msg.Chat.ID, "⚠️ Заголовок не может быть пустым.")
			return
		}
		state.TempTask.Title = text
		r.finishUpdate(ctx, msg.Chat.ID, telegramID, state.TempTask)

	case 13:
		if text == "" {
			r.Reply(msg.Chat.ID, "⚠️ Описание не может быть пустым.")
			return
		}
		state.TempTask.Notes = text
		r.finishUpdate(ctx, msg.Chat.ID, telegramID, state.TempTask)
	}
}
func (r *Router) HandleCallback(ctx context.Context, cb *tgbot.CallbackQuery) {
	data := cb.Data
	chatID := cb.Message.Chat.ID
	msgID := cb.Message.MessageID
	log.Printf("User %s (tg_%d) нажал: %q", cb.From.UserName, chatID, data)

	switch {
	case data == "menu:add":
		userID := strconv.Itoa(int(cb.From.ID))
		userStates[cb.From.ID] = &UserState{
			Step:     1,
			TempTask: models.Task{UserID: userID},
		}
		_, _ = r.bot.Request(tgbot.NewCallback(cb.ID, "Добавление задачи"))
		r.SendCalendar(chatID, time.Now().Year(), time.Now().Month())

	case strings.HasPrefix(data, "calendar_prev:"),
		strings.HasPrefix(data, "calendar_next:"):
		parts := strings.Split(strings.Split(data, ":")[1], "-")
		year, _ := strconv.Atoi(parts[0])
		monthInt, _ := strconv.Atoi(parts[1])
		month := time.Month(monthInt)
		newMarkup := BuildCalendar(year, month)
		edit := tgbot.NewEditMessageReplyMarkup(chatID, msgID, newMarkup)
		r.bot.Send(edit)
		_, _ = r.bot.Request(tgbot.NewCallback(cb.ID, ""))

	case strings.HasPrefix(data, "calendar:"):
		dateStr := strings.TrimPrefix(data, "calendar:")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			r.Reply(chatID, "❌ Ошибка парсинга даты")
			return
		}
		state := userStates[cb.From.ID]
		if state != nil && state.Step == 1 {
			state.TempTask.DueAt = &date
			state.Step = 2
			r.Reply(chatID, fmt.Sprintf("Вы выбрали дату: %s", date.Format("02-01-2006")))
			r.Reply(chatID, "Введите заголовок задачи:")
		} else if state != nil && state.Step == 11 {
			state.TempTask.DueAt = &date
			r.finishUpdate(ctx, chatID, cb.From.ID, state.TempTask)
		} else {
			r.Reply(chatID, "Вы выбрали дату: "+dateStr)
		}

	case data == "menu:list":
		userID := strconv.Itoa(int(cb.From.ID))
		tasks, err := r.todo.ListTasks(ctx, userID)
		if err != nil {
			_, _ = r.bot.Request(tgbot.NewCallback(cb.ID, "Ошибка получения списка"))
			return
		}
		if len(tasks) == 0 {
			_, _ = r.bot.Request(tgbot.NewCallback(cb.ID, "Список пуст"))
			r.Reply(chatID, "Пусто. Используй ➕ Добавить для новой задачи")
			return
		}
		for _, t := range tasks {
			r.SendTask(chatID, t)
		}
		r.SendMainMenu(chatID)

	case data == "menu:listStory":
		userID := strconv.Itoa(int(cb.From.ID))
		tasks, err := r.todo.ListDoneTasks(ctx, userID)
		if err != nil {
			_, _ = r.bot.Request(tgbot.NewCallback(cb.ID, "Ошибка получения списка"))
			return
		}
		if len(tasks) == 0 {
			_, _ = r.bot.Request(tgbot.NewCallback(cb.ID, "Список пуст"))
			r.Reply(chatID, "Пусто. Используй ➕ Добавить для новой задачи")
			return
		}
		for _, t := range tasks {
			r.SendTask(chatID, t)
		}
		r.SendMainMenu(chatID)
	case strings.HasPrefix(data, "complete:"):
		taskID := strings.TrimPrefix(data, "complete:")
		if err := r.todo.CompleteTask(ctx, taskID); err == nil {
			_, _ = r.bot.Request(tgbot.NewCallback(cb.ID, "Задача завершена ✅"))
			newText := cb.Message.Text + "\n✅ Завершена"
			edit := tgbot.NewEditMessageText(chatID, msgID, newText)
			r.bot.Send(edit)
		} else {
			_, _ = r.bot.Request(tgbot.NewCallback(cb.ID, "Ошибка завершения"))
		}

	case strings.HasPrefix(data, "delete:"):
		taskID := strings.TrimPrefix(data, "delete:")
		if err := r.todo.DeleteTask(ctx, taskID); err == nil {
			_, _ = r.bot.Request(tgbot.NewCallback(cb.ID, "Задача удалена 🗑"))
			del := tgbot.NewDeleteMessage(chatID, msgID)
			r.bot.Send(del)
		} else {
			_, _ = r.bot.Request(tgbot.NewCallback(cb.ID, "Ошибка удаления"))
		}

	case strings.HasPrefix(data, "update:"):
		taskID := strings.TrimPrefix(data, "update:")
		userID := strconv.Itoa(int(cb.From.ID))
		userStates[cb.From.ID] = &UserState{
			TempTask: models.Task{
				ID:     taskID,
				UserID: userID,
			},
		}
		_, _ = r.bot.Request(tgbot.NewCallback(cb.ID, "Редактирование задачи ✏️"))

		// Inline‑кнопки для выбора поля
		dateBtn := tgbot.NewInlineKeyboardButtonData("Дата", "update_field:date")
		titleBtn := tgbot.NewInlineKeyboardButtonData("Заголовок", "update_field:title")
		notesBtn := tgbot.NewInlineKeyboardButtonData("Описание", "update_field:notes")

		keyboard := tgbot.NewInlineKeyboardMarkup(
			tgbot.NewInlineKeyboardRow(dateBtn, titleBtn, notesBtn),
		)

		msg := tgbot.NewMessage(chatID, "Что хотите изменить?")
		msg.ReplyMarkup = keyboard
		r.bot.Send(msg)

	case strings.HasPrefix(data, "update_field:"):
		field := strings.TrimPrefix(data, "update_field:")
		state := userStates[cb.From.ID]
		switch field {
		case "date":
			state.Step = 11
			r.Reply(chatID, "Введите новую дату (DD-MM-YYYY):")
		case "title":
			state.Step = 12
			r.Reply(chatID, "Введите новый заголовок:")
		case "notes":
			state.Step = 13
			r.Reply(chatID, "Введите новое описание:")
		}
	case strings.HasPrefix(cb.Data, "calendar:"):
		dateStr := strings.TrimPrefix(cb.Data, "calendar:")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			r.Reply(chatID, "❌ Ошибка парсинга даты")
			return
		}
		state := userStates[cb.From.ID]
		if state != nil && state.Step == 1 {
			state.TempTask.DueAt = &date
			state.Step = 2
			r.Reply(chatID, "Введите заголовок задачи:")
		} else {
			r.Reply(chatID, "Вы выбрали дату: "+dateStr)
		}
	case data == "menu:back":
		r.SendMainMenu(chatID)
	}
}

func (r *Router) finishUpdate(ctx context.Context, chatID int64, telegramID int64, task models.Task) {
	updatedTask, err := r.todo.UpdateTask(ctx, task)
	if err != nil {
		r.Reply(chatID, "❌ Ошибка при обновлении задачи1")
	} else {
		r.Reply(chatID, fmt.Sprintf("✅ Задача обновлена: %s", updatedTask.Title))
	}
	delete(userStates, telegramID)
	r.SendMainMenu(chatID)
}

func (r *Router) Reply(chatID int64, text string) {
	if r == nil || r.bot == nil {
		fmt.Printf("Нет подключения к ТГ: %s\n", text)
		return
	}
	msg := tgbot.NewMessage(chatID, text)
	if _, err := r.bot.Send(msg); err != nil {
		fmt.Printf("telegram send error: %v\n", err)
	}
}
