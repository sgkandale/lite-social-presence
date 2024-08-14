package server

import (
	"encoding/json"
	"log"
	"net/http"
	"socialite/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type MessageType string

const (
	MessageType_Ping                 = "ping"
	MessageType_Pong                 = "pong"
	MessageType_FriendsOnline        = "friends_online"
	MessageType_FriendsOnlineInParty = "friends_online_in_party"
)

type WebsocketStatusIncomingMessage struct {
	MsgType MessageType `json:"msg_type"`
}

type WebsocketStatusOutgoingMessage struct {
	MsgType  MessageType `json:"msg_type"`
	UserName string      `json:"user_name"`
}

func (s *Server) WebsocketStatus(ginCtx *gin.Context) {
	// get user from context
	user, exists := ginCtx.Get(Header_AuthUserKey)
	if !exists || user == nil {
		ginCtx.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
		return
	}
	userInstance := user.(*database.User)
	s.userOnlineStatus <- userInstance.Name

	// upgrade to websocket
	conn, err := s.upgrader.Upgrade(ginCtx.Writer, ginCtx.Request, nil)
	if err != nil {
		ginCtx.JSON(
			http.StatusInternalServerError,
			GeneralResponse{Message: err.Error()},
		)
		return
	}
	defer conn.Close()
	closeChan := make(chan bool)
	writeChan := make(chan []byte)
	s.rwmutex.Lock()
	s.userWebsocketChannels[userInstance.Name] = writeChan
	s.rwmutex.Unlock()

	go func() {
		defer func() { closeChan <- true }()
		for {
			select {
			case <-ginCtx.Done():
				log.Printf("[ERROR] request context expired in status websocket")
				return
			default:
				msgType, msg, err := conn.ReadMessage()
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					log.Printf("[ERROR] websocket conn is closed : %s", err.Error())
					return
				}
				if err != nil {
					if err == websocket.ErrCloseSent {
						return
					}
					log.Printf("[ERROR] reading message from websocket : %s", err.Error())
				}
				if msgType != websocket.TextMessage {
					log.Printf("[ERROR] websocket message type is not text : %d", msgType)
					continue
				}
				incomingMsg := WebsocketStatusIncomingMessage{}
				err = json.Unmarshal(msg, &incomingMsg)
				if err != nil {
					log.Printf("[ERROR] unmarshaling websocket message : %s", err.Error())
					continue
				}
				// handle incoming message
				switch incomingMsg.MsgType {
				case MessageType_Ping:
					resp := WebsocketStatusOutgoingMessage{MsgType: MessageType_Pong}
					respBytes, err := json.Marshal(resp)
					if err != nil {
						log.Printf("[ERROR] marshaling websocket message : %s", err.Error())
						continue
					}
					writeChan <- respBytes
					s.userOnlineStatus <- userInstance.Name
				default:
					log.Printf("[ERROR] unknown message type : %s", incomingMsg.MsgType)
					continue
				}
			}
		}
	}()

	// Goroutine for writing messages
	go func() {
		for {
			select {
			case <-ginCtx.Done():
				return
			case resp := <-writeChan:
				err := conn.WriteMessage(websocket.TextMessage, resp)
				if err != nil {
					log.Printf("[ERROR] writing message to websocket : %s", err.Error())
					return
				}
			}
		}
	}()

	<-closeChan
}

func (s *Server) WebsocketParty(ginCtx *gin.Context) {
	// get user from context
	user, exists := ginCtx.Get(Header_AuthUserKey)
	if !exists || user == nil {
		ginCtx.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
		return
	}
	userInstance := user.(*database.User)
	s.userOnlineStatus <- userInstance.Name

	partyName := ginCtx.Param("party_id")
	if partyName == "" {
		ginCtx.JSON(
			http.StatusBadRequest,
			GeneralResponse{Message: "party_id is required"},
		)
		return
	}

	// check party membership
	_, err := s.db.GetPartyMembership(ginCtx, partyName, userInstance.Name)
	if err != nil {
		if err == database.Err_NotFound {
			ginCtx.JSON(http.StatusNotFound, Err_PartyMembershipNotFound)
			return
		}
		log.Fatalf("[ERROR] getting party membership for partyname %s and user %s", partyName, userInstance.Name)
		ginCtx.JSON(
			http.StatusInternalServerError,
			Err_SomethingWrong,
		)
		return
	}

	// upgrade to websocket
	conn, err := s.upgrader.Upgrade(ginCtx.Writer, ginCtx.Request, nil)
	if err != nil {
		ginCtx.JSON(
			http.StatusInternalServerError,
			GeneralResponse{Message: err.Error()},
		)
		return
	}
	defer conn.Close()

	for {
		time.Sleep(time.Second * 5)
		// get users friends from the cache
		userfriendsList, err := s.cache.GetUserFriendsList(ginCtx, userInstance.Name)
		if err != nil {
			log.Printf("[ERROR] getting user friends from cache : %s", err.Error())
			continue
		}
		// create a map of user's friends
		usersOnlineMap := make(map[string]bool, len(userfriendsList))
		for _, eachFriend := range userfriendsList {
			usersOnlineMap[eachFriend] = false
		}
		// fetch party members from cache
		partyMembers, err := s.cache.GetPartyMembersList(ginCtx, partyName)
		if err != nil {
			log.Printf("[ERROR] getting party members from cache : %s", err.Error())
			continue
		}
		// get a list of users who are online
		for _, eachMember := range partyMembers {
			isUserOnline, err := s.cache.IsUserOnline(ginCtx, eachMember)
			if err != nil {
				log.Printf("[ERROR] checking if user %s is online : %s", eachMember, err.Error())
				continue
			}
			if isUserOnline {
				usersOnlineMap[eachMember] = true
			}
		}
		// delete users who are offline
		for eachUser := range usersOnlineMap {
			if !usersOnlineMap[eachUser] {
				delete(usersOnlineMap, eachUser)
			}
		}
		// send the list to the user
		resp, err := json.Marshal(usersOnlineMap)
		if err != nil {
			log.Printf("[ERROR] marshaling users online map : %s", err.Error())
			continue
		}
		err = conn.WriteMessage(websocket.TextMessage, resp)
		if err != nil {
			log.Printf("[ERROR] writing message to websocket : %s", err.Error())
			return
		}
	}
}
