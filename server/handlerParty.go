package server

import (
	"log"
	"net/http"

	"socialite/database"

	"github.com/gin-gonic/gin"
)

func (s *Server) CreateParty(ginCtx *gin.Context) {
	// get user from context
	user, exists := ginCtx.Get(Header_AuthUserKey)
	if !exists || user == nil {
		ginCtx.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
		return
	}
	userInstance := user.(*database.User)

	// read request body
	var reqBody CreatePartyRequest
	err := ginCtx.BindJSON(&reqBody)
	// handle error
	if err != nil {
		log.Printf("[ERROR] server.CreateParty: reading request body: %s", err.Error())
		ginCtx.JSON(http.StatusBadRequest, Err_ReadingRequest)
		return
	}

	partyInstance, err := database.NewParty(reqBody.Name, userInstance.Name)
	if err != nil {
		log.Printf("[ERROR] server.CreateParty: creating party instance: %s", err.Error())
		ginCtx.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
		return
	}

	// create party
	err = s.db.PutParty(ginCtx, partyInstance)
	if err != nil {
		log.Printf("[ERROR] server.CreateParty: putting party in db: %s", err.Error())
		ginCtx.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
		return
	}

	ginCtx.JSON(http.StatusOK, GeneralResponse{Message: partyInstance.Name})
}

func (s *Server) GetCreatedParties(ginCtx *gin.Context) {
	// get user from context
	user, exists := ginCtx.Get(Header_AuthUserKey)
	if !exists || user == nil {
		ginCtx.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
		return
	}
	userInstance := user.(*database.User)

	parties, err := s.db.GetCreatedParties(ginCtx, userInstance.Name)
	if err != nil {
		log.Printf("[ERROR] server.GetCreatedParties: getting parties from db: %s", err.Error())
		ginCtx.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
		return
	}

	ginCtx.JSON(http.StatusOK, parties)
}
