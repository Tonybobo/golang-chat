package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tonybobo/go-chat/config"
	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/service"
	"github.com/tonybobo/go-chat/pkg/common/utils"
	"github.com/tonybobo/go-chat/pkg/global/log"
)

func Login(ctx *gin.Context) {

}

func Register(ctx *gin.Context) {
	var user entity.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}

	if err := service.UserService.Register(&user); err != nil {
		log.Logger.Error("Fail To Register User", log.String("Error: ", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "Error": err.Error()})
		return
	}
	log.Logger.Info("Successfully Register User", log.Any("User : ", user.Username))
	access_token, _ := utils.CreateToken(config.GetConfig().Token.AccessTokenExpiresIn, user.Id, config.GetConfig().Token.AccessTokenPrivateKey)
	ctx.SetCookie("access_token", access_token, config.GetConfig().Token.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": user})
}

func EditUserDetail(c *gin.Context) {

}

func GetUserDetail(c *gin.Context) {

}
