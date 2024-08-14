package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"socialite/database"

	"github.com/gin-gonic/gin"
)

func (s *Server) AddRoutes() {
	// health routes
	s.engine.Any("/health", s.Health)
	s.engine.Any("/liveness", s.Health)

	// auth routes
	authGroup := s.engine.Group("/auth")
	authGroup.POST("/login", s.Login)
	authGroup.POST("/register", s.Register)

	// all routes below are secured with a middleware
	securedRoutes := s.engine.Group("/")
	securedRoutes.Use(AuthMiddleware())

	// friends routes
	friendsGroup := securedRoutes.Group("/friends")
	friendsGroup.GET("/", s.GetFriends)              // get all friends
	friendsGroup.DELETE("/:user_id", s.RemoveFriend) // remove friend

	// friend requests group
	friendRequestsGroup := friendsGroup.Group("/requests")
	friendRequestsGroup.GET("/", s.GetFriendRequests)                      // get pending friend request
	friendRequestsGroup.POST("/user/:user_id", s.SendFriendRequest)        // send friend request
	friendRequestsGroup.POST("/:request_id/accept", s.AcceptFriendRequest) // accept friend request
	friendRequestsGroup.POST("/:request_id/reject", s.RejectFriendRequest) // reject friend request

	// party routes
	partyGroup := securedRoutes.Group("/party")
	partyGroup.POST("/", s.CreateParty)             // create party
	partyGroup.GET("/created", s.GetCreatedParties) // get created party

	// each party group
	eachPartyGroup := partyGroup.Group("/:party_id")
	eachPartyGroup.POST("/invite", s.InviteUserToParty)            // invite party
	eachPartyGroup.POST("/join", s.JoinParty)                      // join party
	eachPartyGroup.POST("/leave", s.LeaveParty)                    // leave party
	eachPartyGroup.DELETE("/user/:user_id", s.RemoveUserFromParty) // remove user from party

	// websocket group
	websocketGroup := securedRoutes.Group("/ws")
	websocketGroup.Any("/party/:party_id", s.WebsocketParty)
	websocketGroup.Any("/status", s.WebsocketStatus)
}

func (s *Server) AddMiddlewares() {
	s.engine.Use(CORSMiddleware())
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		} else {
			c.Next()
		}
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check for auth token in header
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			c.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
			c.Abort()
			return
		}

		// create user instance from auth token and add to context
		userInstance, err := database.NewUser(authToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, GeneralResponse{Message: err.Error()})
			c.Abort()
			return
		}
		c.Set(Header_AuthUserKey, userInstance)
		c.Next()
	}
}

func (s *Server) StartCrons(ctx context.Context) {
	go s.MonitorOnlineUsers(ctx)
	go s.UpdateUserFriendsListCron(ctx)
	go s.UpdatePartyMembersCron(ctx)
}

func (s *Server) UpdateUserFriendsListCron(ctx context.Context) {
	log.Printf("[INFO] starting cron for updating user's friends list in cache")
	// update user friends list in cache
	err := s.UpdateUserFriendsList(ctx)
	if err != nil {
		log.Printf("[ERROR] updating user friends list in cache : %s", err.Error())
	}

	for {
		time.Sleep(time.Minute)
		if ctx.Err() != nil {
			return
		}
		err = s.UpdateUserFriendsList(ctx)
		if err != nil {
			log.Printf("[ERROR] updating user friends list in cache : %s", err.Error())
		}
	}
}

func (s *Server) UpdateUserFriendsList(ctx context.Context) error {
	log.Print("[INFO] fetching users friends list from db")
	friendsMap, err := s.db.GetUserFriendsList(ctx)
	if err != nil {
		return err
	}
	if len(friendsMap) > 0 {
		log.Print("[INFO] putting user's friends list in cache for total users : ", len(friendsMap))
		for userName, friendsList := range friendsMap {
			err = s.cache.PutUserFriendsList(ctx, userName, friendsList)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

func (s *Server) MonitorOnlineUsers(ctx context.Context) {
	log.Printf("[INFO] starting cron for monitoring user's online status")
	for userName := range s.userOnlineStatus {
		// put user as online in cache
		go s.cache.PutUserOnline(ctx, userName)
		// handle user's online status
		go s.HandleUserOnlineStatus(ctx, userName)
	}
}

func (s *Server) HandleUserOnlineStatus(ctx context.Context, userName string) {
	friendsList, err := s.cache.GetUserFriendsList(ctx, userName)
	if err != nil {
		log.Printf("[ERROR] getting user's friends list from cache : %s", err.Error())
		return
	}

	resp := WebsocketStatusOutgoingMessage{
		MsgType:  MessageType_FriendsOnline,
		UserName: userName,
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[ERROR] marshalling websocket status message : %s", err.Error())
		return
	}

	s.rwmutex.RLock()
	for _, friendName := range friendsList {
		userChan, exist := s.userWebsocketChannels[friendName]
		if !exist || userChan == nil {
			continue
		}
		userChan <- respJson
	}
	s.rwmutex.RUnlock()
}

func (s *Server) UpdatePartyMembersCron(ctx context.Context) {
	err := s.UpdatePartyMembers(ctx)
	if err != nil {
		log.Printf("[ERROR] updating party members list in cache : %s", err.Error())
	}
	for {
		time.Sleep(time.Minute)
		if ctx.Err() != nil {
			return
		}
		err = s.UpdatePartyMembers(ctx)
		if err != nil {
			log.Printf("[ERROR] updating party members list in cache : %s", err.Error())
		}
	}
}

func (s *Server) UpdatePartyMembers(ctx context.Context) error {
	partyMembersMap, err := s.db.GetAllPartyMembers(ctx)
	if err != nil {
		return err
	}
	for partyName, membersList := range partyMembersMap {
		err = s.cache.PutPartyMembersList(ctx, partyName, membersList)
		if err != nil {
			return err
		}
	}
	return nil
}
