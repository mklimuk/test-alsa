package rest

import (
	"github.com/gin-gonic/gin"
)

//API is an interface for publishing routes
type API interface {
	AddRoutes(router *gin.Engine)
}
