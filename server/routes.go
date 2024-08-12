package server

import (
	"net/http"

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
	friendRequestsGroup.POST("/user/:user_id", s.SendFriendRequest) // send friend request
	friendRequestsGroup.POST("/:request_id/accept", nil)            // accept friend request
	friendRequestsGroup.POST("/:request_id/reject", nil)            // reject friend request

	// party routes
	partyGroup := securedRoutes.Group("/party")

	// each party group
	eachPartyGroup := partyGroup.Group("/:party_id")
	eachPartyGroup.POST("/invite", nil)          // invite party
	eachPartyGroup.POST("/join", nil)            // join party
	eachPartyGroup.POST("/leave", nil)           // leave party
	eachPartyGroup.DELETE("/user/:user_id", nil) // remove user from party

	// party invitation group
	partyInvitationsGroup := partyGroup.Group("/invitations")
	partyInvitationsGroup.POST("/:invitation_id", nil) // act on party invitation
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
