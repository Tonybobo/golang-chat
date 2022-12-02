package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/service"
)

func GetGroup(ctx *gin.Context) {

}

func SaveGroup(ctx *gin.Context) {
	uid := ctx.Param("uid")
	var group entity.GroupChat

	if err := ctx.ShouldBindJSON(&group); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}

	if err := service.GroupSerivce.SaveGroup(uid, &group); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func JoinGroup(ctx *gin.Context) {

}

func GetGroupUser(ctx *gin.Context) {

}
