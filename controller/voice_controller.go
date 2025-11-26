package controller

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	MaxRoomCapacity = 2
	DefaultRoomName = "Voice Call"
	RoomTypeGroup   = "group"
	RoomTypePrivate = "private"
)

const (
	MsgTypeError      = "error"
	MsgTypeRoomInfo   = "room-info"
	MsgTypeUserLeft   = "user-left"
	ErrorRoomFull     = "Room is full (max 2 users)"
	ErrorRoomExists   = "Voice room already exists for this team"
	ErrorRoomNotFound = "Voice room not found"
	ErrorUnauthorized = "You are not invited to this call"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type RoomResponse struct {
	entity.VoiceRoom
	UserCount int `json:"userCount" example:"2"`
}

type VoiceController struct {
	mu    sync.RWMutex
	rooms map[string]*entity.VoiceRoom
}

func NewVoiceController() *VoiceController {
	return &VoiceController{
		rooms: make(map[string]*entity.VoiceRoom),
	}
}

// StartPrivateCall creates a new private voice call between two users
//
//	@Summary		Start a private voice call
//	@Description	Creates a private voice room for two users with restricted access
//	@Tags			voice
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			callerId	query		string	true	"ID of the user initiating the call"
//	@Param			targetId	query		string	true	"ID of the user being called"
//	@Param			teamId		query		string	false	"Team ID for context"
//	@Success		201			{object}	entity.VoiceRoom
//	@Failure		400			{object}	map[string]string
//	@Router			/voice/private/call [post]
func (vc *VoiceController) StartPrivateCall(c *gin.Context) {
	callerId := c.Query("callerId")
	targetId := c.Query("targetId")
	teamId := c.Query("teamId")

	if callerId == "" || targetId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both callerId and targetId are required"})
		return
	}

	roomId := generateUniqueId()

	whitelist := make(map[string]bool)
	whitelist[callerId] = true
	whitelist[targetId] = true

	newRoom := &entity.VoiceRoom{
		Id:           roomId,
		TeamId:       teamId,
		Name:         "Private Call",
		Type:         RoomTypePrivate,
		CreatedBy:    callerId,
		CreatedAt:    time.Now().UnixMilli(),
		AllowedUsers: whitelist,
		Clients:      make(map[*websocket.Conn]string),
	}

	vc.mu.Lock()
	vc.rooms[roomId] = newRoom
	vc.mu.Unlock()

	c.JSON(http.StatusCreated, newRoom)
}

// GetJoinableRooms returns all active voice rooms that a user can join
//
//	@Summary		Get joinable voice rooms
//	@Description	Returns all group and private rooms that the user is authorized to join and are not full
//	@Tags			voice
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			userId	query		string	true	"User ID of the client"
//	@Success		200		{array}		controller.RoomResponse
//	@Failure		400		{object}	map[string]string
//	@Router			/voice/joinable [get]
func (vc *VoiceController) GetJoinableRooms(c *gin.Context) {
	userId := c.Query("userId")

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId query parameter is required"})
		return
	}

	var responseList []RoomResponse

	vc.mu.RLock()
	defer vc.mu.RUnlock()

	for _, room := range vc.rooms {
		room.Mutex.RLock()
		userCount := len(room.Clients)
		room.Mutex.RUnlock()

		if userCount >= MaxRoomCapacity {
			continue
		}

		isJoinable := false

		if room.Type == RoomTypeGroup {
			isJoinable = true
		}

		if room.Type == RoomTypePrivate {
			if room.AllowedUsers != nil && room.AllowedUsers[userId] {
				isJoinable = true
			}
		}

		if isJoinable {
			responseList = append(responseList, RoomResponse{
				VoiceRoom: *room,
				UserCount: userCount,
			})
		}
	}

	if responseList == nil {
		responseList = []RoomResponse{}
	}

	c.JSON(http.StatusOK, responseList)
}

// CreateVoiceRoom creates a new group voice room for a team
//
//	@Summary		Create a group voice room
//	@Description	Creates a new voice room for team members
//	@Tags			voice
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			teamId	path		string	true	"Team ID"
//	@Param			userId	query		string	true	"User ID of the creator"
//	@Param			name	query		string	false	"Room name (optional)"
//	@Success		201		{object}	entity.VoiceRoom
//	@Failure		409		{object}	map[string]string
//	@Router			/voice/rooms/{teamId} [post]
func (vc *VoiceController) CreateVoiceRoom(c *gin.Context) {
	teamId := c.Param("teamId")
	userId := c.Query("userId")
	roomName := c.Query("name")

	if roomName == "" {
		roomName = DefaultRoomName
	}

	vc.mu.RLock()
	if _, exists := vc.rooms[teamId]; exists {
		vc.mu.RUnlock()
		c.JSON(http.StatusConflict, gin.H{"error": ErrorRoomExists})
		return
	}
	vc.mu.RUnlock()

	newRoom := &entity.VoiceRoom{
		Id:        teamId,
		TeamId:    teamId,
		Name:      roomName,
		Type:      RoomTypeGroup,
		CreatedBy: userId,
		CreatedAt: time.Now().UnixMilli(),
		Clients:   make(map[*websocket.Conn]string),
	}

	vc.mu.Lock()
	vc.rooms[teamId] = newRoom
	vc.mu.Unlock()

	c.JSON(http.StatusCreated, newRoom)
}

// GetActiveRooms returns all active voice rooms for a specific team
//
//	@Summary		Get active voice rooms for a team
//	@Description	Returns all group voice rooms belonging to a specific team
//	@Tags			voice
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Param			teamId	path	string	true	"Team ID"
//	@Success		200		{array}	controller.RoomResponse
//	@Router			/voice/rooms/{teamId} [get]
func (vc *VoiceController) GetActiveRooms(c *gin.Context) {
	teamId := c.Param("teamId")

	var responseList []RoomResponse

	vc.mu.RLock()
	defer vc.mu.RUnlock()

	for _, room := range vc.rooms {
		if room.TeamId == teamId && room.Type == RoomTypeGroup {
			room.Mutex.RLock()
			count := len(room.Clients)
			room.Mutex.RUnlock()

			responseList = append(responseList, RoomResponse{
				VoiceRoom: *room,
				UserCount: count,
			})
		}
	}

	if responseList == nil {
		responseList = []RoomResponse{}
	}

	c.JSON(http.StatusOK, responseList)
}

// JoinVoiceRoom allows a user to join a voice room via WebSocket connection
//
//	@Summary		Join a voice room via WebSocket
//	@Description	Establishes a WebSocket connection for voice communication in a room
//	@Tags			voice
//	@Security		Bearer
//	@Param			roomId	path		string	true	"Room ID to join"
//	@Param			userId	query		string	true	"User ID joining the room"
//	@Success		101		{string}	string	"Switching Protocols"
//	@Failure		400		{object}	map[string]string
//	@Failure		403		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Router			/voice/join/{roomId} [get]
func (vc *VoiceController) JoinVoiceRoom(c *gin.Context) {
	roomId := c.Param("roomId")
	userId := c.Query("userId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	vc.mu.RLock()
	room, exists := vc.rooms[roomId]
	vc.mu.RUnlock()

	if !exists {
		vc.sendErrorAndClose(conn, ErrorRoomNotFound)
		return
	}

	if room.Type == RoomTypePrivate {
		if room.AllowedUsers == nil || !room.AllowedUsers[userId] {
			vc.sendErrorAndClose(conn, ErrorUnauthorized)
			return
		}
	}

	if !vc.canJoinRoom(room) {
		vc.sendErrorAndClose(conn, ErrorRoomFull)
		return
	}
	room.Mutex.Lock()
	room.Clients[conn] = userId
	room.Mutex.Unlock()

	defer vc.handleUserDisconnect(room, conn, userId, roomId)

	vc.sendRoomInfo(conn, room, true)

	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		if vc.handleRoomInfoRequest(msg, conn, room) {
			continue
		}
		vc.routeToOtherUser(room, conn, msg)
	}
}

func (vc *VoiceController) canJoinRoom(room *entity.VoiceRoom) bool {
	room.Mutex.RLock()
	defer room.Mutex.RUnlock()
	return len(room.Clients) < MaxRoomCapacity
}

func (vc *VoiceController) sendErrorAndClose(conn *websocket.Conn, errorMsg string) {
	vc.sendError(conn, errorMsg)
	conn.Close()
}

func (vc *VoiceController) sendError(conn *websocket.Conn, errorMsg string) {
	conn.WriteJSON(map[string]interface{}{
		"type":  MsgTypeError,
		"error": errorMsg,
	})
}

func (vc *VoiceController) sendRoomInfo(conn *websocket.Conn, room *entity.VoiceRoom, canJoin bool) {
	room.Mutex.RLock()
	userCount := len(room.Clients)
	room.Mutex.RUnlock()

	response := map[string]interface{}{
		"type":      MsgTypeRoomInfo,
		"userCount": userCount,
		"canJoin":   canJoin,
	}
	conn.WriteJSON(response)
}

func (vc *VoiceController) handleRoomInfoRequest(msg map[string]interface{}, conn *websocket.Conn, room *entity.VoiceRoom) bool {
	if msg["type"] == MsgTypeRoomInfo {
		vc.sendRoomInfo(conn, room, vc.canJoinRoom(room))
		return true
	}
	return false
}

func (vc *VoiceController) routeToOtherUser(room *entity.VoiceRoom, sender *websocket.Conn, msg map[string]interface{}) {
	room.Mutex.RLock()
	defer room.Mutex.RUnlock()
	senderId, _ := room.Clients[sender]
	forwardMsg := make(map[string]interface{})
	for k, v := range msg {
		forwardMsg[k] = v
	}
	forwardMsg["from"] = senderId

	for conn := range room.Clients {
		if conn != sender {
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			conn.WriteJSON(forwardMsg)
			conn.SetWriteDeadline(time.Time{})
			break
		}
	}
}

func (vc *VoiceController) handleUserDisconnect(room *entity.VoiceRoom, conn *websocket.Conn, userId, roomId string) {
	room.Mutex.Lock()
	delete(room.Clients, conn)
	room.Mutex.Unlock()

	vc.notifyUserLeft(room, userId)
	if len(room.Clients) == 0 {
		vc.mu.Lock()
		delete(vc.rooms, roomId)
		vc.mu.Unlock()
	}
}

func (vc *VoiceController) notifyUserLeft(room *entity.VoiceRoom, userId string) {
	msg := map[string]interface{}{
		"type":   MsgTypeUserLeft,
		"userId": userId,
	}
	vc.broadcastToAll(room, msg)
}

func (vc *VoiceController) broadcastToAll(room *entity.VoiceRoom, msg map[string]interface{}) {
	room.Mutex.RLock()
	defer room.Mutex.RUnlock()

	for client := range room.Clients {
		client.SetWriteDeadline(time.Now().Add(10 * time.Second))
		client.WriteJSON(msg)
		client.SetWriteDeadline(time.Time{})
	}
}

func generateUniqueId() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "id-" + time.Now().String()
	}
	return hex.EncodeToString(bytes)
}
