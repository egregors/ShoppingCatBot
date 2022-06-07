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

// ShoppingList represents a list of items
type ShoppingList struct {
	items []string
	curr  int

	mu *sync.Mutex
}

// Next returns item from ShoppingList if it's possible
func (l *ShoppingList) Next() (string, bool) {
	if l.curr > len(l.items)-1 {
		l.curr = 0
		return "", false
	}
	l.curr++
	return l.items[l.curr-1], true
}

// Add parses and adds items in ShoppingList
func (l *ShoppingList) Add(text string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	var count int
	rows := strings.Split(text, "\n")
	for _, r := range rows {
		var xs []string
		for _, x := range strings.Split(r, ",") {
			x = strings.TrimSpace(x)
			if x != "" {
				xs = append(xs, x)
				count++
			}
		}
		l.items = append(l.items, xs...)
	}
	return count
}

// Wipe removes all items from ShoppingList
func (l *ShoppingList) Wipe() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.curr = 0
	l.items = []string{}
}

// Len return amount of items in ShoppingList
func (l *ShoppingList) Len() int {
	return len(l.items)
}

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("SCBOT_TG_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Init
	list := ShoppingList{mu: &sync.Mutex{}}

	// Commands
	b.Handle("/list", func(c tele.Context) error {
		if list.Len() == 0 {
			return c.Send("Список пуст")
		}

		if list.Len() == 1 {
			return c.Send("Один пункт ты и так запомнишь")
		}

		pages := len(list.items)/10 + 1
		page := 1
		var ps []*tele.Poll
		p := &tele.Poll{
			Type:            tele.PollRegular,
			MultipleAnswers: true,
			Question:        fmt.Sprintf("Страница %d/%d", 1, pages),
		}

		for item, ok := list.Next(); ok; item, ok = list.Next() {
			if len(p.Options) == 10 {
				page++
				ps = append(ps, p)
				p = &tele.Poll{
					Type:            tele.PollRegular,
					MultipleAnswers: true,
					Question:        fmt.Sprintf("Страница %d/%d", page, pages),
				}
			}
			p.AddOptions(item)
		}
		ps = append(ps, p)

		for _, p := range ps {
			err = c.Send(p)
			if err != nil {
				return err
			}
		}

		return nil
	})

	b.Handle("/add", func(c tele.Context) error {
		t := strings.TrimPrefix(c.Message().Text, "/add")
		return c.Send(fmt.Sprintf("Добавлено %d позиций", list.Add(t)))
	})

	b.Handle("/clean", func(c tele.Context) error {
		list.Wipe()
		return c.Send("Список теперь пуст")
	})

	b.Start()
}
