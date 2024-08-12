package server

import (
	"log"
	"net/http"
	"socialite/database"

	"github.com/gin-gonic/gin"
)

func (s *Server) Login(ginCtx *gin.Context) {
	// read request body
	var reqBody LoginRequest
	err := ginCtx.BindJSON(&reqBody)
	if err != nil {
		log.Printf("[ERROR] server.Login: reading request body: %s", err.Error())
		ginCtx.JSON(http.StatusBadRequest, Err_ReadingRequest)
		return
	}

	// get user from database
	userInstance, err := s.db.GetUser(ginCtx, reqBody.Name)
	if err != nil {
		log.Printf("[ERROR] server.Login: getting user from database: %s", err.Error())
		if err == database.Err_NotFound {
			ginCtx.JSON(http.StatusNotFound, Err_UserNotFound)
			return
		}
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	// return token
	ginCtx.JSON(http.StatusOK, LoginResponse{Token: userInstance.Name})
}

func (s *Server) Register(ginCtx *gin.Context) {
	// read request body
	var reqBody RegisterRequest
	err := ginCtx.BindJSON(&reqBody)
	if err != nil {
		log.Printf("[ERROR] server.Register: reading request body: %s", err.Error())
		ginCtx.JSON(http.StatusBadRequest, Err_ReadingRequest)
		return
	}

	// create new user
	newUserInstance, err := database.NewUser(reqBody.Name)
	if err != nil {
		log.Printf("[ERROR] server.Register: creating new user: %s", err.Error())
		ginCtx.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
		return
	}

	// put user in database
	err = s.db.PutUser(ginCtx, newUserInstance)
	if err != nil {
		log.Printf("[ERROR] server.Register: putting user in database: %s", err.Error())
		if err == database.Err_DuplicatePrimaryKey {
			ginCtx.JSON(http.StatusConflict, Err_UserAlreadyRegistered)
			return
		}
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	// return success
	ginCtx.JSON(http.StatusOK, Resp_Success)
}
