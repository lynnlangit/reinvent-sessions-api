// Package crons manages application's cron goroutines
package crons

import (
	"github.com/robfig/cron"
)

var crons []*cron.Cron

// Crons lists all crons
func Crons() []*cron.Cron {
	return crons
}

// NewCron creates new cron
func NewCron() *cron.Cron {
	cron := cron.New()
	crons = append(crons, cron)
	return cron
}
