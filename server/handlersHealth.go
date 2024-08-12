package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Health(ginCtx *gin.Context) {
	ginCtx.JSON(
		http.StatusOK,
		HealthResponse{
			Status: "OK",
		},
	)
}
