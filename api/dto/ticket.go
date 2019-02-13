package dto

type Ticket struct {
	Ticket string `json:"ticket"`
}

func NewTicket(ticket string) Ticket {
	return Ticket{
		Ticket: ticket,
	}
}
