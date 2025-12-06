package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
)

func (r *Router) CheckDueTasks(ctx context.Context) {
	now := time.Now().UTC()
	tasks, err := r.todo.ListDueTasks(ctx, now)
	if err != nil {
		log.Printf("Ошибка при получении задач по дате: %v", err)
		return
	}

	for _, task := range tasks {
		chatID, err := strconv.ParseInt(task.UserID, 10, 64)
		if err != nil {
			log.Printf("Невалидный UserID: %v", task.UserID)
			continue
		}

		text := fmt.Sprintf("⏰ Напоминание: задача \"%s\" на сегодня (%s)", task.Title, task.DueAt.Format("02-01-2006"))
		r.Reply(chatID, text)
	}
}
func (r *Router) StartDueTaskNotifier(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Оповещающий воркер остановлен")
				return
			case <-ticker.C:
				r.CheckDueTasks(ctx)
			}
		}
	}()
}
