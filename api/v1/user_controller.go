package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/service"
	"github.com/tonybobo/go-chat/pkg/global/log"
)

func Login(ctx *gin.Context) {
	var login entity.Login
	if err := ctx.BindJSON(&login); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}

	if user, ok := service.UserService.Login(&login); ok {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": entity.FilteredResponse(user)})
		return
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "Error": "Incorrect Username/Password"})
		return
	}
}

func Register(ctx *gin.Context) {
	var register entity.Register
	if err := ctx.ShouldBindJSON(&register); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}

	user, err := service.UserService.Register(&register)
	if err != nil {
		log.Logger.Error("Fail To Register User", log.String("Error: ", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}

	log.Logger.Info("Successfully Register User", log.Any("User : ", user.Username))
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": entity.FilteredResponse(user)})
}

func EditUserDetail(ctx *gin.Context) {
	var user entity.EditUser

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}
	if err := service.UserService.EditUserDetail(&user); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "Message": "User Detail has been updated"})

}

func GetUserDetail(ctx *gin.Context) {
	uid := ctx.Param("uid")
	log.Logger.Info("uid", log.String("uid", uid))
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": service.UserService.GetUserDetails(uid)})
}

func GetUserOrGroupByName(ctx *gin.Context) {
	name := ctx.Param("name")
	result := service.UserService.GetUsersOrGroupBy(name)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}

func AddFriend(ctx *gin.Context) {
	var request entity.FriendRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}

	if err := service.UserService.AddFriend(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetFriends(ctx *gin.Context) {
	uid := ctx.Query("uid")

	friends, err := service.UserService.GetFriends(uid)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": friends})
}
