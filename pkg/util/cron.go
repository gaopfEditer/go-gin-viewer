package util

import (
	"github.com/robfig/cron/v3"
)

var context *cron.Cron

func init() {
	// Seconds field, optional
	context = cron.New(cron.WithSeconds())

	context.Start()
}
