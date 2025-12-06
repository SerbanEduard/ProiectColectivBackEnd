package entity

import (
	"sync"

	"github.com/gorilla/websocket"
)

type VoiceRoom struct {
	Id        string `json:"id"`
	TeamId    string `json:"teamId"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedBy string `json:"createdBy"`
	CreatedAt int64  `json:"createdAt"`

	AllowedUsers map[string]bool            `json:"-"`
	Clients      map[*websocket.Conn]string `json:"-"`
	Mutex        sync.RWMutex               `json:"-"`

	// Current screenshare presenter userId (empty if none)
	ScreenPresenter string `json:"-"`
}
