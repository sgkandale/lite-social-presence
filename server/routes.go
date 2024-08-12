package server

func (s *Server) AddRoutes() {
	// health routes
	s.engine.Any("/health", nil)
	s.engine.Any("/liveness", nil)

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
