package cron

import (
	"license/internal/jetbrains/code/task"
	"license/internal/logger"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	c := cron.New(cron.WithSeconds())

	jetbrainsTask := task.NewTask()

	// Add scheduled task
	_, err := c.AddFunc("0 0 1 * * ?", func() {
		jetbrainsTask.FetchProductLatest()
	})
	if err != nil {
		logger.Error("Failed to add cron job:", err)
	}
	c.Start()
	logger.Sys("Cron job started")
}
