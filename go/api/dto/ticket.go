package dto

// Ticket is a ticket which allows a new websocket session to be negotiated.
type Ticket struct {
	Ticket string `json:"ticket"`
}

// NewTicket makes a new ticket.
func NewTicket(ticket string) Ticket {
	return Ticket{
		Ticket: ticket,
	}
}
