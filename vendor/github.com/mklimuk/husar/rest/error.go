package rest

import (
	"net/http"
	"runtime/debug"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

//ErrorHandler writes unexpected errors to the HTTP response. To be used with defer.
func ErrorHandler(ctx *gin.Context) {
	if r := recover(); r != nil {
		clog := log.WithFields(log.Fields{"logger": "rest.error", "method": "ErrorHandler"})
		clog.WithField("error", r).Error("Recovered from unexpected error.")
		debug.PrintStack()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected internal error occurred.", "details": r})
	}
}
