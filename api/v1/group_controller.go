package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/service"
)

func GetGroups(ctx *gin.Context) {
	uid := ctx.Param("uid")

	groups, err := service.GroupSerivce.GetGroups(uid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": groups})
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
	userUid := ctx.Param("userUid")
	groupUid := ctx.Param("groupUid")

	group , err := service.GroupSerivce.JoinGroup(userUid, groupUid);

	if  err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success" , "data" : group})
}

func GetGroupUser(ctx *gin.Context) {
	groupUid := ctx.Param("uid")

	users, err := service.GroupSerivce.GetGroupUsers(groupUid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": users})
}
