package external

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/xuchangeu/kong-assignment/src/constant"
	"github.com/xuchangeu/kong-assignment/src/utility"
	"net/url"
)

type UserAuthResp struct {
	UserId           string `json:"userId"`
	UserRole         string `json:"userRole"`
	UserOrganization string `json:"organization"`
	Msg              string `json:"msg"`
	Code             int64  `json:"code"`
}

type UserProjectsResp struct {
}

func UserAuth(ctx *gin.Context, token string) (*UserAuthResp, error) {
	ctx = utility.GenerateUUIDForContext(ctx)
	ctxId, _ := ctx.Get(constant.KeyContextId)
	var absPath = "http://127.0.0.1:8080/user/auth"
	reqPath, _ := url.JoinPath(absPath, token)
	client := resty.New()
	log.WithFields(log.Fields{
		"message":             "making external request for auth",
		"path":                reqPath,
		constant.KeyContextId: ctxId,
	}).Info()
	resp, err := client.R().Get(reqPath)
	if err != nil {
		log.WithFields(log.Fields{
			"message": err.Error(),
		}).Info()
		return nil, err
	}
	var authResp UserAuthResp
	err = json.Unmarshal(resp.Body(), &authResp)
	if err != nil {
		log.WithFields(log.Fields{
			"message":             err.Error(),
			"code":                constant.DataUnmarshalErrorCode,
			constant.KeyContextId: ctxId,
		}).Error()
	}
	return &authResp, nil
}

func RetrieveUserProjects() {

}
