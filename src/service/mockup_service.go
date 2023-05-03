package service

import (
	log "github.com/sirupsen/logrus"
	"github.com/xuchangeu/kong-assignment/src/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuchangeu/kong-assignment/src/configuration"
	"github.com/xuchangeu/kong-assignment/src/constant"
)

func UserAuth(ctx *gin.Context) {
	ctx = utility.GenerateUUIDForContext(ctx)
	token := ctx.Param(constant.HTTPReqKeyUserInPath)
	ctxId, _ := ctx.Get(constant.KeyContextId)
	mockData := configuration.GetMockUpSession().Session

	for k, v := range mockData {
		if k == token {
			log.WithFields(log.Fields{
				constant.KeyMessage:          "user found in session",
				constant.KeyToken:            token,
				constant.KeyUserRole:         v.UserRole,
				constant.KeyUserOrganization: v.UserOrganization,
				constant.KeyContextId:        ctxId,
			}).Info()
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				constant.KeyUserId:           k,
				constant.KeyUserRole:         v.UserRole,
				constant.KeyUserOrganization: v.UserOrganization,
				constant.KeyCode:             constant.UserAuthStatusOK,
				constant.KeyMessage:          "",
			})
			return
		}

	}

	ctx.JSON(http.StatusOK, gin.H{
		constant.KeyUserId:           "",
		constant.KeyUserRole:         "",
		constant.KeyUserOrganization: "",
		constant.KeyCode:             constant.UserAuthTokenNotExistCode,
		constant.KeyMessage:          constant.UserAuthTokenNotExistMsg,
	})
}

func UserProject(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"projects": []int{111, 222, 333},
	})
}
