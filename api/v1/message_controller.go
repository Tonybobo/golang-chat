package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/service"
)

func GetMessage(ctx *gin.Context) {
	var messageRequest *entity.MessageRequest

	if err := ctx.BindJSON(&messageRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}
	limit := ctx.Query("limit")
	pages := ctx.Query("pages")

	limitInt, err := strconv.Atoi(limit)

	if err != nil {
		limitInt = 20
	}

	pagesInt, err := strconv.Atoi(pages)

	if err != nil {
		pagesInt = 1
	}

	message, err := service.MessageService.GetMessages(limitInt, pagesInt, messageRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": message})

}
