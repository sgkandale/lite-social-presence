package server

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"socialite/database"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetFriends(ginCtx *gin.Context) {
	// get user from context
	user, exists := ginCtx.Get(Header_AuthUserKey)
	if !exists || user == nil {
		ginCtx.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
		return
	}
	userInstance := user.(*database.User)

	// get friends from database
	friends, err := s.db.GetUserFriends(ginCtx, userInstance.Name)
	if err != nil {
		log.Printf("[ERROR] server.GetFriends: getting friends: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	ginCtx.JSON(http.StatusOK, friends)
}

func (s *Server) RemoveFriend(ginCtx *gin.Context) {
	// get user from context
	user, exists := ginCtx.Get(Header_AuthUserKey)
	if !exists || user == nil {
		ginCtx.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
		return
	}
	userInstance := user.(*database.User)

	// get friend name from query
	friendName := ginCtx.Query("user_id")
	if friendName == "" {
		ginCtx.JSON(http.StatusBadRequest, Err_UserIdMissing)
		return
	}

	// get friendship details from db
	friendship, err := s.db.GetFriendship(ginCtx, userInstance.Name, friendName)
	if err != nil {
		if err == database.Err_NotFound {
			ginCtx.JSON(http.StatusNotFound, Err_FriendshipNotFound)
			return
		}
		log.Printf("[ERROR] server.RemoveFriend: getting friendship: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	// remove friendship
	err = s.db.DeleteFriendship(ginCtx, friendship.Id)
	if err != nil {
		if err == database.Err_NotFound {
			ginCtx.JSON(http.StatusNotFound, Err_FriendshipNotFound)
			return
		}
		log.Printf("[ERROR] server.RemoveFriend: deleting friendship: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	ginCtx.JSON(http.StatusOK, Resp_Success)
}

func (s *Server) SendFriendRequest(ginCtx *gin.Context) {
	// get user from context
	user, exists := ginCtx.Get(Header_AuthUserKey)
	if !exists || user == nil {
		ginCtx.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
		return
	}
	userInstance := user.(*database.User)

	// get friend name from query
	friendName := ginCtx.Query("user_id")
	if friendName == "" {
		ginCtx.JSON(http.StatusBadRequest, Err_UserIdMissing)
		return
	}

	// check if friendship already exists
	_, err := s.db.GetFriendship(ginCtx, userInstance.Name, friendName)
	if err == nil {
		ginCtx.JSON(http.StatusConflict, Err_FriendshipAlreadyExists)
		return
	} else if err != database.Err_NotFound {
		log.Printf("[ERROR] server.SendFriendRequest: getting friendship: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	friendshipInstance, err := database.NewFriendship(userInstance.Name, friendName)
	if err != nil {
		log.Printf("[ERROR] server.SendFriendRequest: creating friendship instance: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	// create friendship
	err = s.db.PutFriendship(ginCtx, friendshipInstance)
	if err != nil {
		log.Printf("[ERROR] server.SendFriendRequest: creating friendship: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	ginCtx.JSON(http.StatusOK, Resp_Success)
}

func (s *Server) AcceptFriendRequest(ginCtx *gin.Context) {
	// get user from context
	user, exists := ginCtx.Get(Header_AuthUserKey)
	if !exists || user == nil {
		ginCtx.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
		return
	}
	userInstance := user.(*database.User)

	// get request id from path
	requestId := ginCtx.Param("request_id")
	if requestId == "" {
		ginCtx.JSON(http.StatusBadRequest, Err_RequestIdMissing)
		return
	}

	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, Err_RequestIdInvalid)
		return
	}

	friendshipInstance, err := s.db.GetFriendshipById(ginCtx, int32(requestIdInt))
	if err != nil {
		if err == database.Err_NotFound {
			ginCtx.JSON(http.StatusNotFound, Err_FriendshipRequestNotFound)
			return
		}
		log.Printf("[ERROR] server.AcceptFriendRequest: getting friendship request: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	if friendshipInstance.User1 != userInstance.Name && friendshipInstance.User2 != userInstance.Name {
		ginCtx.JSON(http.StatusNotFound, Err_FriendshipRequestNotFound)
		return
	}

	friendshipInstance.Status = database.Friendship_Status_Confirmed
	friendshipInstance.UpdatedAt = time.Now()

	// get friendship details from db
	err = s.db.UpdateFriendship(ginCtx, friendshipInstance)
	if err != nil {
		if err == database.Err_NotFound {
			ginCtx.JSON(http.StatusNotFound, Err_FriendshipRequestNotFound)
			return
		}
		log.Printf("[ERROR] server.AcceptFriendRequest: getting friendship request: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	ginCtx.JSON(http.StatusOK, Resp_Success)
}

func (s *Server) RejectFriendRequest(ginCtx *gin.Context) {
	// get user from context
	user, exists := ginCtx.Get(Header_AuthUserKey)
	if !exists || user == nil {
		ginCtx.JSON(http.StatusUnauthorized, Err_AuthHeaderMissing)
		return
	}
	userInstance := user.(*database.User)

	// get request id from path
	requestId := ginCtx.Param("request_id")
	if requestId == "" {
		ginCtx.JSON(http.StatusBadRequest, Err_RequestIdMissing)
		return
	}

	requestIdInt, err := strconv.Atoi(requestId)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, Err_RequestIdInvalid)
		return
	}

	friendshipInstance, err := s.db.GetFriendshipById(ginCtx, int32(requestIdInt))
	if err != nil {
		if err == database.Err_NotFound {
			ginCtx.JSON(http.StatusNotFound, Err_FriendshipRequestNotFound)
			return
		}
		log.Printf("[ERROR] server.RejectFriendRequest: getting friendship request: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	if friendshipInstance.User1 != userInstance.Name && friendshipInstance.User2 != userInstance.Name {
		ginCtx.JSON(http.StatusNotFound, Err_FriendshipRequestNotFound)
		return
	}

	err = s.db.DeleteFriendship(ginCtx, friendshipInstance.Id)
	if err != nil {
		if err == database.Err_NotFound {
			ginCtx.JSON(http.StatusNotFound, Err_FriendshipRequestNotFound)
			return
		}
		log.Printf("[ERROR] server.RejectFriendRequest: deleting friendship request: %s", err.Error())
		ginCtx.JSON(http.StatusInternalServerError, Err_SomethingWrong)
		return
	}

	ginCtx.JSON(http.StatusOK, Resp_Success)
}
