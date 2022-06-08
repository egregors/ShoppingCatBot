package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	bunchSize = 10
)

// ItemStorager represents storage for items
type ItemStorager interface {
	// Add adds item into particular chatID bucket
	Add(chatID int64, item string)
	// Wipe removes all items for particular Chat
	Wipe(chatID int64)
	// GetAll return collection of bunches (size <= 10, because tg Poll could contain only <= 10 option)
	// with items. As long, as I use tg Polls to show lists it should be so. ¯\_(ツ)_/¯
	GetAll(chatID int64) [][]string
}

// Srv is runnable instance if shopping list
type Srv struct {
	db  ItemStorager
	bot *tele.Bot
}

// NewServer takes ItemStorager and tele.Bot and initializes all handlers
func NewServer(db ItemStorager, b *tele.Bot) *Srv {
	const (
		addCmd  = "/add"
		listCmd = "/list"
		wipeCmd = "/wipe"
	)

	b.Handle(listCmd, func(c tele.Context) error {
		items := db.GetAll(c.Message().Chat.ID)

		if len(items) == 0 {
			return c.Send("Список пуст")
		}

		// we should do it because Tg Polls can't have less than two options
		if len(items) == 1 && len(items[0]) == 1 {
			return c.Send(fmt.Sprintf("В списке пока только один пункт: %s", items[0][0]))
		}

		var polls []*tele.Poll

		for page, group := range items {
			p := &tele.Poll{
				Type:            tele.PollRegular,
				MultipleAnswers: true,
				Question:        fmt.Sprintf("Список покупок, страница: %d/%d", page+1, len(items)),
			}
			for _, item := range group {
				p.AddOptions(item)
			}
			polls = append(polls, p)
		}

		// sand all polls one by one back to tg
		for _, p := range polls {
			err := c.Send(p)
			if err != nil {
				return fmt.Errorf("can't send poll: %w", err)
			}
		}

		return nil
	})

	b.Handle(addCmd, func(c tele.Context) error {
		text := strings.TrimPrefix(c.Message().Text, addCmd)
		rows := strings.Split(text, "\n")
		itemCount := 0
		for _, r := range rows {
			ws := strings.Split(r, ",")
			for _, w := range ws {
				item := strings.TrimSpace(w)
				if item != "" {
					db.Add(c.Message().Chat.ID, item)
					itemCount++
				}
			}
		}
		return c.Send(fmt.Sprintf("Добавлено %d пунктов", itemCount))
	})

	b.Handle(wipeCmd, func(c tele.Context) error {
		db.Wipe(c.Message().Chat.ID)
		return c.Send("Список теперь пуст")
	})

	return &Srv{
		db:  db,
		bot: b,
	}
}

// Run stars a Srv
func (s *Srv) Run() error {
	s.bot.Start()
	return nil
}

// Inmem is in-memory implementation of ItemStorager
type Inmem struct {
	items map[int64][]string

	mu *sync.Mutex
}

// todo: do something with Add, Wipe, GetAll comments
var _ ItemStorager = (*Inmem)(nil)

//nolint:revive // see interface docs
func (db *Inmem) Add(chatID int64, item string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.items[chatID] = append(db.items[chatID], item)
}

//nolint:revive // see interface docs
func (db *Inmem) Wipe(chatID int64) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.items, chatID)
}

//nolint:revive // see interface docs
func (db *Inmem) GetAll(chatID int64) [][]string {
	if items, ok := db.items[chatID]; !ok || len(items) == 0 {
		return nil
	}

	var (
		items [][]string
		loc   []string
	)

	for i := 0; i < len(db.items[chatID]); i++ {
		if len(loc) == bunchSize {
			items = append(items, loc)
			loc = []string{}
		}
		loc = append(loc, db.items[chatID][i])
	}
	items = append(items, loc)

	return items
}

func makeBot(timeOut time.Duration) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  os.Getenv("SCBOT_TG_TOKEN"),
		Poller: &tele.LongPoller{Timeout: timeOut},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("can't create a bot: %w", err)
	}
	return b, nil
}

func makeInmemStore() *Inmem {
	return &Inmem{
		items: make(map[int64][]string),
		mu:    &sync.Mutex{},
	}
}

func main() {
	b, err := makeBot(10 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
	db := makeInmemStore()
	srv := NewServer(db, b)

	// todo: catch SIGTERM to commit data to disk before shutdown
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
