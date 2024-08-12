package server

import (
	"log"
	"net/http"

	"socialite/database"

	"github.com/gin-gonic/gin"
)

func (s *Server) InviteUserToParty(ginCtx *gin.Context) {
	// get user from context
	user, exists := ginCtx.Get(Header_AuthUserKey)
	if !exists || user == nil {
		ginCtx.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
		return
	}
	userInstance := user.(*database.User)

	// get party name from path
	partyName := ginCtx.Param("party_id")
	if partyName == "" {
		ginCtx.JSON(http.StatusBadRequest, Err_ReadingRequest)
		return
	}

	// read request body
	var reqBody InviteUserToPartyRequest
	err := ginCtx.BindJSON(&reqBody)
	if err != nil {
		log.Printf("[ERROR] server.InviteUserToParty: reading request body: %s", err.Error())
		ginCtx.JSON(http.StatusBadRequest, Err_ReadingRequest)
		return
	}

	party, err := s.db.GetParty(ginCtx, partyName)
	if err != nil {
		if err == database.Err_NotFound {
			ginCtx.JSON(http.StatusNotFound, Err_PartyNotFound)
			return
		}
		log.Printf("[ERROR] server.InviteUserToParty: getting party from db: %s", err.Error())
		ginCtx.JSON(http.StatusBadRequest, Err_ReadingRequest)
		return
	}

	// only creator should be able to invite
	if party.Creator != userInstance.Name {
		ginCtx.JSON(http.StatusUnauthorized, Err_NotPartyCreator)
		return
	}

	partyMembership, err := database.NewPartyMembership(party.Name, reqBody.UserName)
	if err != nil {
		log.Printf("[ERROR] server.InviteUserToParty: creating party membership instance: %s", err.Error())
		ginCtx.JSON(http.StatusBadRequest, Err_ReadingRequest)
		return
	}

	err = s.db.PutPartyMembership(ginCtx, partyMembership)
	if err != nil {
		if err == database.Err_DuplicatePrimaryKey {
			ginCtx.JSON(http.StatusConflict, Err_UserAlreadyInParty)
			return
		}
		log.Printf("[ERROR] server.InviteUserToParty: putting party membership to db: %s", err.Error())
		ginCtx.JSON(http.StatusBadRequest, Err_ReadingRequest)
		return
	}

	ginCtx.JSON(http.StatusOK, Resp_Success)
}

func (s *Server) JoinParty(ginCtx *gin.Context) {
	// get user from context
	user, exists := ginCtx.Get(Header_AuthUserKey)
	if !exists || user == nil {
		ginCtx.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
		return
	}
	userInstance := user.(*database.User)

	// get party name from path
	partyName := ginCtx.Param("party_id")
	if partyName == "" {
		ginCtx.JSON(http.StatusBadRequest, Err_ReadingRequest)
		return
	}

	// get party membership
	partyMembership, err := s.db.GetPartyMembership(ginCtx, partyName, userInstance.Name)
	if err != nil {
		if err == database.Err_NotFound {
			ginCtx.JSON(http.StatusNotFound, Err_PartyInvitationNotFound)
			return
		}
		log.Printf("[ERROR] server.JoinParty: getting party membership from db: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	if partyMembership.Status != database.PartyMembership_Status_Invited {
		ginCtx.JSON(http.StatusBadRequest, Err_PartyInvitationNotFound)
		return
	}

	partyMembership.Status = database.PartyMembership_Status_Active

	err = s.db.UpdatePartyMembership(ginCtx, partyMembership)
	if err != nil {
		log.Printf("[ERROR] server.JoinParty: updating party membership to db: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	ginCtx.JSON(http.StatusOK, Resp_Success)
}

func (s *Server) LeaveParty(ginCtx *gin.Context) {
	// get user from context
	user, exists := ginCtx.Get(Header_AuthUserKey)
	if !exists || user == nil {
		ginCtx.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
		return
	}
	userInstance := user.(*database.User)

	// get party name from path
	partyName := ginCtx.Param("party_id")
	if partyName == "" {
		ginCtx.JSON(http.StatusBadRequest, Err_ReadingRequest)
		return
	}

	// get party membership
	partyMembership, err := s.db.GetPartyMembership(ginCtx, partyName, userInstance.Name)
	if err != nil {
		if err == database.Err_NotFound {
			ginCtx.JSON(http.StatusNotFound, Err_PartyMembershipNotFound)
			return
		}
		log.Printf("[ERROR] server.LeaveParty: getting party membership from db: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	err = s.db.DeletePartyMembership(ginCtx, partyMembership)
	if err != nil {
		log.Printf("[ERROR] server.LeaveParty: deleting party membership from db: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	ginCtx.JSON(http.StatusOK, Resp_Success)
}

func (s *Server) RemoveUserFromParty(ginCtx *gin.Context) {
	// get user from context
	user, exists := ginCtx.Get(Header_AuthUserKey)
	if !exists || user == nil {
		ginCtx.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
		return
	}
	userInstance := user.(*database.User)

	// get party name from path
	partyName := ginCtx.Param("party_id")
	if partyName == "" {
		ginCtx.JSON(http.StatusBadRequest, Err_ReadingRequest)
		return
	}

	// get user name from path
	userName := ginCtx.Param("user_id")
	if userName == "" {
		ginCtx.JSON(http.StatusBadRequest, Err_ReadingRequest)
		return
	}

	// get party
	party, err := s.db.GetParty(ginCtx, partyName)
	// check if party exists
	if err != nil {
		if err == database.Err_NotFound {
			ginCtx.JSON(http.StatusNotFound, Err_PartyNotFound)
			return
		}
		log.Printf("[ERROR] server.RemoveUserFromParty: getting party from db: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	// only creator should be able to remove
	if party.Creator != userInstance.Name {
		ginCtx.JSON(http.StatusUnauthorized, Err_NotPartyCreator)
		return
	}

	// get party membership
	partyMembership, err := s.db.GetPartyMembership(ginCtx, partyName, userName)
	if err != nil {
		if err == database.Err_NotFound {
			ginCtx.JSON(http.StatusNotFound, Err_PartyMembershipNotFound)
			return
		}
		log.Printf("[ERROR] server.RemoveUserFromParty: getting party membership from db: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	err = s.db.DeletePartyMembership(ginCtx, partyMembership)
	if err != nil {
		log.Printf("[ERROR] server.RemoveUserFromParty: deleting party membership from db: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	ginCtx.JSON(http.StatusOK, Resp_Success)
}
