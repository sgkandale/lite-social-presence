package server

import (
	"log"
	"net/http"

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
