package socket

import (
	"net/http"

	"github.com/hednowley/sound/socket/dto"
)

type MockHub struct {
}

func NewMockHub() IHub {
	return &MockHub{}
}

func (h *MockHub) SetHandler(method string, handler Handler) {
}

// Run starts the hub.
func (h *MockHub) Run() {
}

// Notify sends a notification to all clients.
func (h *MockHub) Notify(notification *dto.Notification) {
}

// Notify sends a notification to all clients.
func (h *MockHub) AddClient(ticketer *Ticketer, w http.ResponseWriter, r *http.Request) {
}
