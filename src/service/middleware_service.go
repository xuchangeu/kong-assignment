package service

import (
	"github.com/gin-gonic/gin"
	"github.com/xuchangeu/kong-assignment/src/configuration"
	"github.com/xuchangeu/kong-assignment/src/constant"
	"github.com/xuchangeu/kong-assignment/src/external"
	"github.com/xuchangeu/kong-assignment/src/utility"

	log "github.com/sirupsen/logrus"

	"fmt"
	"net/http"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c = utility.GenerateUUIDForContext(c)
		authWeatherAdminRequire(true, c)
	}
}

func AuthMiddleWare2() gin.HandlerFunc {
	return func(c *gin.Context) {
		c = utility.GenerateUUIDForContext(c)
		authWeatherAdminRequire(false, c)
	}
}

func authWeatherAdminRequire(require bool, c *gin.Context) {
	//suppose to use user-id as token which retrieve from cookie and doing verify
	ctxId, _ := c.Get(constant.KeyContextId)
	token, err := c.Cookie("user-id")
	if err != nil {
		log.WithFields(log.Fields{
			constant.KeyMessage:   "retrieve token error",
			constant.KeyError:     err.Error(),
			constant.KeyContextId: ctxId,
		}).Info()
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			constant.KeyMessage: constant.TokenEmptyErrorMsg,
			constant.KeyCode:    constant.TokenEmptyErrorCode,
		})
		return
	}
	log.WithFields(log.Fields{
		constant.KeyMessage:   "before making authentication",
		constant.KeyToken:     token,
		constant.KeyContextId: ctxId,
	}).Info()

	resp, err := external.UserAuth(c, token)
	if err != nil {
		log.WithFields(log.Fields{
			constant.KeyMessage:   fmt.Sprintf("making user auth failure : %s", err.Error()),
			constant.KeyCode:      resp.Code,
			constant.KeyToken:     token,
			constant.KeyContextId: ctxId,
		}).Error()
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"msg":  constant.UserAuthErrorMsg,
			"code": constant.UserAuthErrorCode,
		})
		return
	}

	if resp != nil && resp.Code != constant.UserAuthStatusOK {
		log.WithFields(log.Fields{
			constant.KeyMessage:   fmt.Sprintf("making user auth failure : %s", resp.Msg),
			constant.KeyCode:      resp.Code,
			constant.KeyToken:     token,
			constant.KeyContextId: ctxId,
		}).Error()
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			constant.KeyMessage: constant.UserAuthErrorMsg,
			constant.KeyCode:    constant.UserAuthErrorCode,
		})
		return
	}

	if require {
		if resp != nil && resp.Code == constant.UserAuthStatusOK && resp.UserRole != constant.UserRoleAdministrator {
			log.WithFields(log.Fields{
				constant.KeyMessage:   "authentication failure, role of user's not administrator",
				constant.KeyCode:      resp.Code,
				constant.KeyToken:     token,
				constant.KeyContextId: ctxId,
			}).Warn()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg":  constant.UserAuthErrorNonAdministratorMsg,
				"code": constant.UserAuthErrorNonAdministratorCode,
			})
			return
		}

	}
	c.Set(constant.ContextKeyUserInfo, configuration.UserSession{
		UserRole:         resp.UserRole,
		UserOrganization: resp.UserOrganization,
		UserId:           resp.UserId,
	})
	c.Next()
}
