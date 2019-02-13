package ws

import (
	"encoding/base64"
	"time"

	"crypto/rand"

	"github.com/hednowley/sound/config"
)

type Ticket struct {
	user    *config.User
	expires time.Time
}

func (t *Ticket) hasExpired() bool {
	return t.expires.Before(time.Now())
}

type Ticketer struct {
	duration time.Duration
	tickets  map[string]Ticket
}

func NewTicketer(config *config.Config) *Ticketer {
	return &Ticketer{
		duration: time.Second * time.Duration(config.WebsocketTicketExpiry),
		tickets:  make(map[string]Ticket),
	}
}

func (t *Ticketer) MakeTicket(user *config.User) string {
	t.cleanTickets()

	b := make([]byte, 20)
	rand.Read(b)
	s := base64.URLEncoding.EncodeToString(b)

	t.tickets[s] = Ticket{
		expires: time.Now().Add(t.duration),
		user:    user,
	}

	return s
}

func (t *Ticketer) SubmitTicket(value string) *config.User {
	if ticket, ok := t.tickets[value]; ok {
		// Delete so ticket can't be used twice
		delete(t.tickets, value)
		return ticket.user
	}

	return nil
}

func (t *Ticketer) cleanTickets() {
	for value, ticket := range t.tickets {
		if ticket.hasExpired() {
			delete(t.tickets, value)
		}
	}
}
