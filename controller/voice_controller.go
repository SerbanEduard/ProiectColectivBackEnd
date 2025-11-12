package controller

import (
	"net/http"
	"sync"
	"time"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	MaxRoomCapacity = 2
	DefaultRoomName = "Voice Chat"
)

const (
	MsgTypeError      = "error"
	MsgTypeRoomInfo   = "room-info"
	MsgTypeUserLeft   = "user-left"
	ErrorRoomFull     = "Room is full (max 2 users)"
	ErrorRoomExists   = "Voice room already exists for this team"
	ErrorRoomNotFound = "Voice room not found"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type VoiceController struct {
	rooms map[string]*VoiceRoom
}

type VoiceRoom struct {
	Id      string
	Clients map[*websocket.Conn]string
	Mutex   sync.RWMutex
}

func NewVoiceController() *VoiceController {
	return &VoiceController{
		rooms: make(map[string]*VoiceRoom),
	}
}

// JoinVoiceRoom
// @Summary Join voice chat room
// @Param teamId path string true "Team ID"
// @Param userId query string true "User ID"
// @Router /voice/{teamId} [get]
func (vc *VoiceController) JoinVoiceRoom(c *gin.Context) {
	teamId := c.Param("teamId")
	userId := c.Query("userId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	if vc.rooms[teamId] == nil {
		vc.rooms[teamId] = vc.createNewRoom(teamId)
	}

	room := vc.rooms[teamId]

	vc.cleanupDeadConnections(room)

	if !vc.canJoinRoom(room) {
		vc.sendError(conn, ErrorRoomFull)
		return
	}

	vc.addClientToRoom(room, conn, userId)
	vc.sendRoomInfo(conn, room, true)

	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			delete(room.Clients, conn)
			if len(room.Clients) == 0 {
				delete(vc.rooms, teamId)
			}
			break
		}

		if vc.handleRoomInfoRequest(msg, conn, room) {
			continue
		}

		vc.broadcastMessage(room, conn, msg)
	}
}

// CreateVoiceRoom
// @Summary Create voice room
// @Accept json
// @Produce json
// @Param teamId path string true "Team ID"
// @Success 201 {object} entity.VoiceRoom
// @Router /voice/rooms/{teamId} [post]
// Helper methods

func (vc *VoiceController) canJoinRoom(room *VoiceRoom) bool {
	return len(room.Clients) < MaxRoomCapacity
}

func (vc *VoiceController) sendError(conn *websocket.Conn, errorMsg string) {
	conn.WriteJSON(map[string]interface{}{
		"type":  MsgTypeError,
		"error": errorMsg,
	})
}

func (vc *VoiceController) addClientToRoom(room *VoiceRoom, conn *websocket.Conn, userId string) {
	room.Clients[conn] = userId
}

func (vc *VoiceController) sendRoomInfo(conn *websocket.Conn, room *VoiceRoom, canJoin bool) {
	response := map[string]interface{}{
		"type":      MsgTypeRoomInfo,
		"userCount": len(room.Clients),
		"canJoin":   canJoin,
	}
	conn.WriteJSON(response)
}

func (vc *VoiceController) handleRoomInfoRequest(msg map[string]interface{}, conn *websocket.Conn, room *VoiceRoom) bool {
	if msg["type"] == MsgTypeRoomInfo {
		vc.sendRoomInfo(conn, room, vc.canJoinRoom(room))
		return true
	}
	return false
}

func (vc *VoiceController) broadcastMessage(room *VoiceRoom, sender *websocket.Conn, msg map[string]interface{}) {
	room.Mutex.RLock()
	var clients []*websocket.Conn
	for client := range room.Clients {
		if client != sender {
			clients = append(clients, client)
		}
	}
	room.Mutex.RUnlock()

	for _, client := range clients {
		client.WriteJSON(msg)
	}
}

func (vc *VoiceController) cleanupDeadConnections(room *VoiceRoom) {
	for conn := range room.Clients {
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			delete(room.Clients, conn)
		}
	}
}

func (vc *VoiceController) createNewRoom(teamId string) *VoiceRoom {
	return &VoiceRoom{
		Id:      teamId,
		Clients: make(map[*websocket.Conn]string),
		Mutex:   sync.RWMutex{},
	}
}

func (vc *VoiceController) buildRoomResponse(teamId, userId string) *entity.VoiceRoom {
	return &entity.VoiceRoom{
		Id:           teamId,
		TeamId:       teamId,
		Name:         DefaultRoomName,
		IsActive:     true,
		Participants: []string{},
		CreatedBy:    userId,
		CreatedAt:    time.Now().UnixMilli(),
	}
}

func (vc *VoiceController) removeUserFromRoom(room *VoiceRoom, userId string) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	for conn, id := range room.Clients {
		if id == userId {
			delete(room.Clients, conn)
			conn.Close()
			break
		}
	}
}

func (vc *VoiceController) notifyUserLeft(room *VoiceRoom, userId string) {
	msg := map[string]interface{}{
		"type":   MsgTypeUserLeft,
		"userId": userId,
	}
	vc.broadcastToAll(room, msg)
}

func (vc *VoiceController) broadcastToAll(room *VoiceRoom, msg map[string]interface{}) {
	room.Mutex.RLock()
	var clients []*websocket.Conn
	for client := range room.Clients {
		clients = append(clients, client)
	}
	room.Mutex.RUnlock()

	for _, client := range clients {
		client.WriteJSON(msg)
	}
}

func (vc *VoiceController) CreateVoiceRoom(c *gin.Context) {
	teamId := c.Param("teamId")
	userId := c.Query("userId")

	if vc.rooms[teamId] != nil {
		c.JSON(http.StatusConflict, gin.H{"error": ErrorRoomExists})
		return
	}

	vc.rooms[teamId] = vc.createNewRoom(teamId)
	room := vc.buildRoomResponse(teamId, userId)
	c.JSON(http.StatusCreated, room)
}

// LeaveVoiceRoom
// @Summary Leave voice chat room
// @Param teamId path string true "Team ID"
// @Param userId query string true "User ID"
// @Router /voice/{teamId}/leave [delete]
func (vc *VoiceController) LeaveVoiceRoom(c *gin.Context) {
	teamId := c.Param("teamId")
	userId := c.Query("userId")

	room := vc.rooms[teamId]
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrorRoomNotFound})
		return
	}

	vc.removeUserFromRoom(room, userId)
	vc.notifyUserLeft(room, userId)

	if len(room.Clients) == 0 {
		delete(vc.rooms, teamId)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left room successfully"})
}
