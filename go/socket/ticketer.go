package socket

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/hednowley/sound/config"
	"time"
)

// Ticket allows a user to negotiate a websocket session.
type Ticket struct {
	user    *config.User
	expires time.Time
}

func (t *Ticket) hasExpired() bool {
	return t.expires.Before(time.Now())
}

// Ticketer creates and monitors tickets.
type Ticketer struct {
	// How long after its creation a ticket expires
	duration time.Duration
	tickets  map[string]Ticket
}

// NewTicketer creates a new ticketer.
func NewTicketer(config *config.Config) *Ticketer {
	return &Ticketer{
		duration: time.Second * time.Duration(config.WebsocketTicketExpiry),
		tickets:  make(map[string]Ticket),
	}
}

// MakeTicket creates a new ticket.
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

// SubmitTicket sees if there is a ticket with the given key.
// If there is, then the user which created the ticket is returned.
// Otherwise returns nil.
func (t *Ticketer) SubmitTicket(key string) *config.User {
	if ticket, ok := t.tickets[key]; ok {
		// Delete so ticket can't be used twice
		delete(t.tickets, key)
		return ticket.user
	}

	return nil
}

// Removes all expired tickets.
func (t *Ticketer) cleanTickets() {
	for value, ticket := range t.tickets {
		if ticket.hasExpired() {
			delete(t.tickets, value)
		}
	}
}
