package server

import "github.com/gin-gonic/gin"

func (s *Server) AddRoutes() {
	// health routes
	s.engine.Any("/health", s.Health)
	s.engine.Any("/liveness", s.Health)

	// friends routes
	friendsGroup := s.engine.Group("/friends")
	friendsGroup.GET("/")            // get all friends
	friendsGroup.DELETE("/:user_id") // remove friend

	// friend requests group
	friendRequestsGroup := friendsGroup.Group("/requests")
	friendRequestsGroup.POST("/user/:user_id", nil) // send friend request
	friendRequestsGroup.POST("/:request_id", nil)   // act on friend request

	// party routes
	partyGroup := s.engine.Group("/party")

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
		} else {
			c.Next()
		}
	}
}
