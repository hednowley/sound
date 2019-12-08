package ws

import (
	"net/http"

	"github.com/hednowley/sound/ws/dto"
)

type MockHub struct {
}

func NewMockHub() *MockHub {
	return &MockHub{}
}

func (h *MockHub) SetHandler(method string, handler WsHandler) {
}

// Run starts the hub.
func (h *MockHub) Run() {
}

// Notify sends a notification to all clients.
func (h *MockHub) Notify(notification *dto.Notification) {
}

// Notify sends a notification to all clients.
func (h *MockHub) AddClient(ticketer Ticketer, w http.ResponseWriter, r *http.Request) {
}
