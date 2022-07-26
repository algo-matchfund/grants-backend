package watchdog

import (
	"context"
	"errors"
	"log"

	"github.com/algo-matchfund/grants-backend/internal/config"
	"github.com/algo-matchfund/grants-backend/internal/database"
)

type Watchdog interface {
	// Watchdog initialization, should be done once
	Init() bool
	// Starts tracking goroutine
	StartWatch(context.Context)
	// Track state of something
	Watch(interface{}) bool
}

type WatchdogFactory struct {
	config *config.Config
	db     *database.GrantsDatabase

	instances map[string]Watchdog
}

func NewWatchdogFactory(config *config.Config, db *database.GrantsDatabase) *WatchdogFactory {
	return &WatchdogFactory{
		config:    config,
		db:        db,
		instances: make(map[string]Watchdog),
	}
}

func (wf *WatchdogFactory) GetWatchdog(watchdogType string) (Watchdog, error) {
	// return existing instances
	if watchdogInstance, ok := wf.instances[watchdogType]; ok {
		return watchdogInstance, nil
	}

	// create new watchdog instance if it doesn't exist yet
	switch watchdogType {
	case "algorand":
		wd, err := NewAlgorandWatchdog(wf.config, wf.db)
		if err != nil {
			log.Printf("WatchdogFactory.GetWatchdog: failed to create watchdog %s, reason - %s\n", watchdogType, err)
			return nil, err
		}

		initSuccess := wd.Init()
		if !initSuccess {
			return nil, errors.New("WatchdogFactory.GetWatchdog: failed to initialize watchdog of type " + watchdogType)
		}
		wf.instances[watchdogType] = wd
		return wf.instances[watchdogType], nil
	default:
		return nil, errors.New("WatchdogFactory.GetWatchdog: Unknown watchdog type " + watchdogType)
	}
}
