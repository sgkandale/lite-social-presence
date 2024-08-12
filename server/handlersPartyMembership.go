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

	// read request body
	var reqBody InviteUserToPartyRequest
	err := ginCtx.BindJSON(&reqBody)
	if err != nil {
		log.Printf("[ERROR] server.InviteUserToParty: reading request body: %s", err.Error())
		ginCtx.JSON(http.StatusBadRequest, Err_ReadingRequest)
		return
	}

	party, err := s.db.GetParty(ginCtx, reqBody.PartyName)
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
