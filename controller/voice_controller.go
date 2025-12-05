package controller

import (
	"crypto/rand"
	"encoding/hex"

	"log"
	"net/http"
	"sync"
	"time"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	MaxRoomCapacity = 10
	DefaultRoomName = "Voice Call"
	RoomTypeGroup   = "group"
	RoomTypePrivate = "private"
)

const (
	MsgTypeError         = "error"
	MsgTypeRoomInfo      = "room-info"
	MsgTypeUserLeft      = "user-left"
	MsgTypeScreenState   = "screen-state"
	MsgTypeScreenStart   = "screenshare-start"
	MsgTypeScreenStop    = "screenshare-stop"
	ErrorRoomFull        = "Room is full (max users)"
	ErrorRoomExists      = "Voice room already exists for this team"
	ErrorRoomNotFound    = "Voice room not found"
	ErrorUnauthorized    = "You are not invited to this call"
	ErrorPresenterActive = "A presenter is already active"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type RoomResponse struct {
	entity.VoiceRoom
	UserCount int `json:"userCount" example:"2"`
}

type VoiceController struct {
	userService  UserServiceInterface
	mu           sync.RWMutex
	rooms        map[string]*entity.VoiceRoom
	pendingDel   map[string]bool // tracks rooms scheduled for deletion
	cleanupDelay time.Duration   // deletion grace period
}

// NewVoiceController constructs the controller
func NewVoiceController() *VoiceController {
	return &VoiceController{
		userService:  service.NewUserService(),
		rooms:        make(map[string]*entity.VoiceRoom),
		pendingDel:   make(map[string]bool),
		cleanupDelay: 5 * time.Second,
	}
}

// StartPrivateCall creates a new private voice call between two users
//
//	@Summary		Start a private voice call
//	@Description	Creates a private voice room for two users with restricted access
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

	if vc.pendingDel[roomId] {
		delete(vc.pendingDel, roomId)
	}
	vc.mu.Unlock()

	c.JSON(http.StatusCreated, newRoom)
}

// GetJoinableRooms returns all active voice rooms that a user can join
//
//	@Summary		Get joinable voice rooms
//	@Description	Returns all group and private rooms that the user is authorized to join and are not full
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
		log.Printf("[voice] CreateVoiceRoom: room already exists for teamId=%s", teamId)
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

	if vc.pendingDel[teamId] {
		delete(vc.pendingDel, teamId)
	}
	vc.mu.Unlock()

	log.Printf("[voice] CreateVoiceRoom: created roomId=%s name=%q by userId=%s", teamId, roomName, userId)
	c.JSON(http.StatusCreated, newRoom)
}

// GetActiveRooms returns all active voice rooms for a specific team
//
//	@Summary		Get active voice rooms for a team
//	@Description	Returns all group voice rooms belonging to a specific team
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

	log.Printf("[voice] JoinVoiceRoom: request roomId=%s userId=%s", roomId, userId)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	vc.mu.RLock()
	room, exists := vc.rooms[roomId]
	// Fallback: if not found by key, search by TeamId match (group rooms)
	if !exists {
		for _, r := range vc.rooms {
			if r != nil && r.TeamId == roomId && r.Type == RoomTypeGroup {
				room = r
				exists = true
				break
			}
		}
	}
	vc.mu.RUnlock()

	if !exists {
		// Log existing room IDs to help diagnose mismatched ids
		vc.mu.RLock()
		keys := make([]string, 0, len(vc.rooms))
		for k := range vc.rooms {
			keys = append(keys, k)
		}
		vc.mu.RUnlock()
		log.Printf("[voice] JoinVoiceRoom: room not found for id=%s; existing=%v", roomId, keys)
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

	vc.mu.Lock()
	if vc.pendingDel[roomId] {
		delete(vc.pendingDel, roomId)
	}
	vc.mu.Unlock()

	room.Mutex.Lock()
	room.Clients[conn] = userId
	room.Mutex.Unlock()

	log.Printf("[voice] JoinVoiceRoom: joined roomId=%s userId=%s currentUsers=%d", room.Id, userId, func() int { room.Mutex.RLock(); defer room.Mutex.RUnlock(); return len(room.Clients) }())

	defer vc.handleUserDisconnect(room, conn, userId, room.Id)

	// Emit a hello message to confirm delivery path
	_ = vc.safeWriteToConn(conn, userId, map[string]interface{}{
		"type":    "hello",
		"message": "welcome",
	})
	// Also broadcast hello to others
	vc.broadcastToAll(room, map[string]interface{}{
		"type":    "hello",
		"from":    userId,
		"message": "user joined",
	})

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

func (vc *VoiceController) sendError(conn *websocket.Conn, errorMsg string) {
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_ = conn.WriteJSON(map[string]interface{}{
		"type":  MsgTypeError,
		"error": errorMsg,
	})
	conn.SetWriteDeadline(time.Time{})
}

// sendErrorAndClose sends an error message and closes the websocket connection.
func (vc *VoiceController) sendErrorAndClose(conn *websocket.Conn, errorMsg string) {
	vc.sendError(conn, errorMsg)
	_ = conn.Close()
}

// safeWriteToConn sends JSON on conn, sets deadline and logs the error.
func (vc *VoiceController) safeWriteToConn(conn *websocket.Conn, uid string, msg interface{}) error {
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err := conn.WriteJSON(msg)
	conn.SetWriteDeadline(time.Time{})
	if err != nil {
		log.Printf("websocket write to user %s failed: %v", uid, err)
	}
	return err
}

// getUsername tries to resolve a username for a given userId.
func (vc *VoiceController) getUsername(userId string) string {
	user, err := vc.userService.GetUserByID(userId)
	if err != nil {
		return ""
	}
	return user.Username
}

// sendRoomInfo writes the room-info message to a connection.
func (vc *VoiceController) sendRoomInfo(conn *websocket.Conn, room *entity.VoiceRoom, canJoin bool) {
	log.Printf("[voice] sendRoomInfo: canJoin=%v presenterId=%s", canJoin, room.ScreenPresenter)
	room.Mutex.RLock()
	userCount := len(room.Clients)
	var usersList []map[string]string
	for _, uid := range room.Clients {
		uname := vc.getUsername(uid)
		usersList = append(usersList, map[string]string{
			"userId":   uid,
			"username": uname,
		})
	}
	presenter := room.ScreenPresenter
	room.Mutex.RUnlock()

	response := map[string]interface{}{
		"type":        MsgTypeRoomInfo,
		"userCount":   userCount,
		"users":       usersList,
		"canJoin":     canJoin,
		"presenterId": presenter,
	}

	room.Mutex.RLock()
	uid, ok := room.Clients[conn]
	room.Mutex.RUnlock()
	if !ok {
		return
	}

	if err := vc.safeWriteToConn(conn, uid, response); err != nil {
		log.Printf("[voice] sendRoomInfo: write failed to uid=%s err=%v", uid, err)
		vc.handleUserDisconnect(room, conn, uid, room.Id)
		_ = conn.Close()
	}
}

func (vc *VoiceController) handleRoomInfoRequest(msg map[string]interface{}, conn *websocket.Conn, room *entity.VoiceRoom) bool {
	if msg["type"] == MsgTypeRoomInfo {
		vc.sendRoomInfo(conn, room, vc.canJoinRoom(room))
		return true
	}
	return false
}

func (vc *VoiceController) routeToOtherUser(room *entity.VoiceRoom, sender *websocket.Conn, msg map[string]interface{}) {
	// Intercept screenshare control messages
	if t, ok := msg["type"].(string); ok {
		switch t {
		case MsgTypeScreenStart:
			log.Printf("[voice] routeToOtherUser: screenshare-start from sender")
			vc.handleScreenStart(room, sender)
			return
		case MsgTypeScreenStop:
			log.Printf("[voice] routeToOtherUser: screenshare-stop from sender")
			vc.handleScreenStop(room, sender)
			return
		}
	}

	room.Mutex.RLock()
	senderId, _ := room.Clients[sender]
	forwardMsg := make(map[string]interface{})
	for k, v := range msg {
		forwardMsg[k] = v
	}
	forwardMsg["from"] = senderId

	var failed []struct {
		conn *websocket.Conn
		uid  string
	}

	if toRaw, ok := msg["to"]; ok {
		if toId, ok2 := toRaw.(string); ok2 {
			for conn, uid := range room.Clients {
				if uid == toId {
					log.Printf("[voice] routeToOtherUser: targeted send to uid=%s type=%v", uid, msg["type"])
					if err := vc.safeWriteToConn(conn, uid, forwardMsg); err != nil {
						failed = append(failed, struct {
							conn *websocket.Conn
							uid  string
						}{conn, uid})
					}
					break
				}
			}
			room.Mutex.RUnlock()
			for _, f := range failed {
				vc.handleUserDisconnect(room, f.conn, f.uid, room.Id)
				_ = f.conn.Close()
			}
			return
		}
	}

	for conn, uid := range room.Clients {
		if conn != sender {
			log.Printf("[voice] routeToOtherUser: broadcast send to uid=%s type=%v", uid, msg["type"])
			if err := vc.safeWriteToConn(conn, uid, forwardMsg); err != nil {
				failed = append(failed, struct {
					conn *websocket.Conn
					uid  string
				}{conn, uid})
			}
		}
	}
	room.Mutex.RUnlock()

	for _, f := range failed {
		vc.handleUserDisconnect(room, f.conn, f.uid, room.Id)
		_ = f.conn.Close()
	}
}

func (vc *VoiceController) handleUserDisconnect(room *entity.VoiceRoom, conn *websocket.Conn, userId, roomId string) {
	room.Mutex.Lock()
	delete(room.Clients, conn)
	// Clear presenter if the user was sharing the screen
	if room.ScreenPresenter == userId {
		room.ScreenPresenter = ""
	}
	room.Mutex.Unlock()

	vc.notifyUserLeft(room, userId)

	// If presenter cleared, notify peers of screen-state change
	if userId != "" {
		vc.broadcastToAll(room, map[string]interface{}{
			"type":        MsgTypeScreenState,
			"presenterId": room.ScreenPresenter,
			"active":      room.ScreenPresenter != "",
		})
	}

	room.Mutex.RLock()
	remaining := len(room.Clients)
	room.Mutex.RUnlock()

	if remaining == 0 {
		vc.mu.Lock()
		if !vc.pendingDel[roomId] {
			vc.pendingDel[roomId] = true
			delay := vc.cleanupDelay
			go func(rid string, d time.Duration) {
				time.Sleep(d)
				vc.mu.Lock()
				defer vc.mu.Unlock()
				if r, ok := vc.rooms[rid]; ok {
					r.Mutex.RLock()
					cnt := len(r.Clients)
					r.Mutex.RUnlock()
					if cnt == 0 {
						delete(vc.rooms, rid)
					}
				}
				delete(vc.pendingDel, rid)
			}(roomId, delay)
		}
		vc.mu.Unlock()
	}
}

func (vc *VoiceController) handleScreenStart(room *entity.VoiceRoom, sender *websocket.Conn) {
	room.Mutex.Lock()
	uid := room.Clients[sender]
	if room.ScreenPresenter != "" && room.ScreenPresenter != uid {
		_ = vc.safeWriteToConn(sender, uid, map[string]interface{}{
			"type":  MsgTypeError,
			"error": ErrorPresenterActive,
		})
		room.Mutex.Unlock()
		return
	}
	room.ScreenPresenter = uid
	users := len(room.Clients)
	room.Mutex.Unlock()
	log.Printf("[voice] handleScreenStart: presenter=%s users=%d", uid, users)
	vc.broadcastToAll(room, map[string]interface{}{
		"type":        MsgTypeScreenState,
		"presenterId": uid,
		"active":      true,
	})
}

func (vc *VoiceController) handleScreenStop(room *entity.VoiceRoom, sender *websocket.Conn) {
	room.Mutex.Lock()
	uid := room.Clients[sender]
	shouldBroadcast := false
	if room.ScreenPresenter == uid {
		room.ScreenPresenter = ""
		shouldBroadcast = true
	}
	room.Mutex.Unlock()
	if shouldBroadcast {
		log.Printf("[voice] handleScreenStop: presenter cleared by %s", uid)
		vc.broadcastToAll(room, map[string]interface{}{
			"type":        MsgTypeScreenState,
			"presenterId": "",
			"active":      false,
		})
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
	log.Printf("[voice] broadcastToAll: type=%v", msg["type"])
	room.Mutex.RLock()
	var failed []struct {
		conn *websocket.Conn
		uid  string
	}
	for client, uid := range room.Clients {
		log.Printf("[voice] broadcastToAll: sending to uid=%s", uid)
		if err := vc.safeWriteToConn(client, uid, msg); err != nil {
			log.Printf("[voice] broadcastToAll: write failed to uid=%s err=%v", uid, err)
			failed = append(failed, struct {
				conn *websocket.Conn
				uid  string
			}{client, uid})
		}
	}
	room.Mutex.RUnlock()

	for _, f := range failed {
		log.Printf("[voice] broadcastToAll: disconnecting failed uid=%s", f.uid)
		vc.handleUserDisconnect(room, f.conn, f.uid, room.Id)
		_ = f.conn.Close()
	}
}

func generateUniqueId() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "id-" + time.Now().String()
	}
	return hex.EncodeToString(bytes)
}
