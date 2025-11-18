package telegram

import (
	"TODO/adapter/driven/models"
	"TODO/services/contracts"
	"context"
	"fmt"
	"log"
	"strings"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	bot  *tgbot.BotAPI
	todo contracts.TodoRepository
}

func NewRouter(bot *tgbot.BotAPI, todo contracts.TodoRepository) *Router {
	return &Router{
		bot:  bot,
		todo: todo}
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
			if upd.Message == nil {
				continue
			}
			r.HandleMessage(ctx, upd.Message)
		}
	}
}

func (r *Router) HandleMessage(ctx context.Context, msg *tgbot.Message) {
	if msg == nil || msg.From == nil {
		return
	}
	telegramID := msg.From.ID
	log.Printf("User %s (tg_%d) написал: %q", msg.From.UserName, msg.From.ID, msg.Text)
	userID := fmt.Sprintf("tg_%d", telegramID)
	// Очищаю текст от пробелов
	text := strings.TrimSpace(msg.Text)
	switch {
	case strings.HasPrefix(text, "/add"):
		title := parseAddCommand(text)
		if title == "" {
			r.Reply(msg.Chat.ID, "Ошибка: заголовок пустой. Используй: /add Купить хлеб")
			return
		}
		task, _ := r.todo.CreateTask(ctx, models.Task{
			UserID: userID,
			Title:  title,
		})
		r.Reply(msg.Chat.ID, "Добавлено:"+task.Title+task.ID)
	case strings.HasPrefix(text, "/list"):
		tasks, err := r.todo.ListTasks(ctx, userID)
		if err != nil {
			log.Printf("ListTasks error: %v", err)
			r.Reply(msg.Chat.ID, "Да , тут ошибка списка")
			return
		}
		if len(tasks) == 0 {
			r.Reply(msg.Chat.ID, "Пусто. Используй /add <задача>")
			return
		}
		var b strings.Builder
		for _, t := range tasks {
			status := "⏳"
			if t.Done {
				status = "✅"
			}
			fmt.Fprintf(&b, "%s %s — %s\n", status, t.ID, t.Title)
		}
		r.Reply(msg.Chat.ID, b.String())
	case strings.HasPrefix(text, "/done"):
		id := strings.TrimSpace(strings.TrimPrefix(text, "/done"))
		if id == "" {
			r.Reply(msg.Chat.ID, "Ошибка: укажи ID задачи. Пример: /done tsk_20251116213200")
			return
		}
		err := r.todo.CompleteTask(ctx, id)
		r.Reply(msg.Chat.ID, "Готово: %s"+id)
		if err != nil {
			r.Reply(msg.Chat.ID, "Ошибка: укажи ID задачи. Пример: /done tsk_20251116213200")
			return
		}
	default:
		r.Reply(msg.Chat.ID, "Команды: /add, /list, /done")
	}
}

func (r *Router) Reply(chatID int64, text string) {
	if r == nil || r.bot == nil {
		fmt.Printf("telegram bot is nil, reply skipped: %s\n", text)
		return
	}
	msg := tgbot.NewMessage(chatID, text)
	if _, err := r.bot.Send(msg); err != nil {
		fmt.Printf("telegram send error: %v\n", err)
	}
}

func parseAddCommand(text string) string {
	const cmd = "/add"
	if !strings.HasPrefix(text, cmd) {
		return ""
	}
	return text[len(cmd):]
}
