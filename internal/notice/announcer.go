package notice

import (
	"sync"

	"github.com/li4n0/revsuit/internal/record"
	log "unknwon.dev/clog/v2"
)

var (
	lock      sync.RWMutex
	announcer *Announcer
)

// Bot is used to send notice
type Bot interface {
	notice(record.Record) error
	buildPayload(record.Record) string
}

// Announcer is used for storage and schedule bots
type Announcer struct {
	Bots []Bot
}

// New initializes a new Announcer
func New() *Announcer {
	announcer = &Announcer{Bots: make([]Bot, 0)}
	return announcer
}

// AddBot add a new bot to Announcer
func (a *Announcer) AddBot(b Bot) *Announcer {
	lock.RLock()
	a.Bots = append(a.Bots, b)
	lock.RUnlock()
	return a
}

// Notice let bots to send notice
func Notice(r record.Record) {
	for _, bot := range announcer.Bots {
		if err := bot.notice(r); err != nil {
			log.Error(err.Error())
		}
	}
}
