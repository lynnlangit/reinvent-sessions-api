package crons

import (
	"github.com/supinf/reinvent-sessions-api/app/logs"
	"github.com/supinf/reinvent-sessions-api/app/models"
)

func init() {
	c := NewCron()

	// SyncReInventSessions
	c.AddFunc("@every 60m", func() {
		models.SyncReInventSessions(true)
	})
	logs.Info("[cron] SyncReInventSessions runs every 60 minutes.")

	c.Start()
}
