package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/xuchangeu/invalid"
	"github.com/xuchangeu/kong-assignment/src/configuration"
	"github.com/xuchangeu/kong-assignment/src/constant"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type LintingService struct {
}

type LintingCreateResp struct {
	Type   string
	ErrMsg string
	Range  *invalid.Range
}

// LintingCreate crate new linting rule if it's not exists, otherwise override the old one.
func (s LintingService) LintingCreate(ctx *gin.Context) {
	userInfo, exists := s.checkUserInContextExistsAndReturn(ctx)
	if !exists {
		ctx.Abort()
		return
	}

	lintId := ctx.Param(constant.HTTPReqKeyLintInPath)
	var org *configuration.Organization
	var exist bool
	if lintId == "new" {
		lintId = ""
		org, exist = s.organizationExist(ctx, userInfo.UserId, userInfo.UserOrganization)
		if !exist {
			ctx.Abort()
			return
		}
	} else {
		org, exist = s.checkOrganizationAndLintExistsAndReturn(ctx, userInfo.UserId, userInfo.UserOrganization, lintId)
		if !exist {
			ctx.Abort()
			return
		}
	}

	formFile, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			constant.HTTPRespKeyCode:    constant.HTTPRespCodeParamMissing,
			constant.HTTPRespKeyMessage: err.Error(),
		})

		ctx.Abort()
		return
	}

	file, err := formFile.Open()
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	by, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			constant.HTTPRespKeyCode:    constant.ServiceInternalErrorCode,
			constant.HTTPRespKeyMessage: err.Error(),
		})
		ctx.Abort()
		return
	}

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			constant.HTTPRespKeyCode:    constant.ServiceInternalErrorCode,
			constant.HTTPRespKeyMessage: err.Error(),
		})
		ctx.Abort()
		return
	}

	buffer := bytes.NewReader(by)
	field, err := invalid.NewYAML(buffer)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			constant.HTTPRespKeyCode:    constant.ServiceInternalErrorCode,
			constant.HTTPRespKeyMessage: fmt.Sprintf("parse yaml content error : %s", err.Error()),
		})
		ctx.Abort()
		return
	}

	ruleFile, err := os.OpenFile(filepath.Join([]string{"yaml", "openapi_rule.yaml"}...), os.O_RDONLY, os.ModeSticky)
	defer func() {
		if ruleFile != nil {
			ruleFile.Close()
		}
	}()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			constant.HTTPRespKeyCode:    constant.ServiceInternalErrorCode,
			constant.HTTPRespKeyMessage: fmt.Sprintf("parse yaml content error : %s", err.Error()),
		})
		ctx.Abort()
		return
	}
	rule, err := invalid.NewRule(ruleFile)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			constant.HTTPRespKeyCode:    constant.ServiceInternalErrorCode,
			constant.HTTPRespKeyMessage: fmt.Sprintf("parse rule error : %s", err.Error()),
		})
		ctx.Abort()
		return
	}

	validResult := rule.Validate(field)
	resp := make([]LintingCreateResp, len(validResult))
	for i := range validResult {
		resp = append(resp, LintingCreateResp{
			Type:   string(validResult[i].Type),
			ErrMsg: validResult[i].Error.Error(),
			Range:  validResult[i].Range,
		})
	}

	buffer = bytes.NewReader(by)
	if lintId == "" {
		lintId = uuid.New().String()
	}
	err = org.WriteRule(ctx, userInfo.UserOrganization, lintId, buffer)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			constant.HTTPRespKeyCode:    constant.ServiceInternalErrorCode,
			constant.HTTPRespKeyMessage: fmt.Sprintf("write rule error : %s", err.Error()),
		})
		ctx.Abort()
		return
	}
	ctx.JSON(200, gin.H{
		constant.HTTPRespKeyCode:     "linting create ok",
		constant.HTTPRespKeyUserId:   userInfo.UserId,
		constant.HTTPRespKeyUserRole: userInfo.UserRole,
		constant.HTTPRespKeyUserOrg:  userInfo.UserOrganization,
		constant.HTTPRespKeyMessage:  resp,
	})
}

// LintingApply This function apply an existing linting rule to a project inside an organization.
// The relationship between them is 1 to n which means a project must should have only 1 applying rule,
// while a rule could be applied to any project inside the organization.
func (s LintingService) LintingApply(ctx *gin.Context) {
	userInfo, exists := s.checkUserInContextExistsAndReturn(ctx)
	if !exists {
		ctx.Abort()
		return
	}

	lintId, exists := s.paramCheckAndReturn(ctx)
	if !exists {
		ctx.Abort()
		return
	}

	org, exists := s.checkOrganizationAndLintExistsAndReturn(ctx, userInfo.UserId, userInfo.UserOrganization, lintId)
	if !exists {
		ctx.Abort()
		return
	}

	proj, exists := s.checkProjectIdAndReturn(ctx, org)
	if !exists {
		ctx.Abort()
		return
	}

	err := org.ApplyLintToProj(ctx, userInfo.UserOrganization, proj.ProjId, lintId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			constant.HTTPRespKeyCode:    constant.HTTPRespCodeInternalError,
			constant.HTTPRespKeyMessage: constant.HTTPRespMessageInternalError,
		})
		return
	}

	ctx.JSON(200, gin.H{
		constant.HTTPRespKeyCode:     "linting apply ok",
		constant.HTTPRespKeyUserId:   userInfo.UserId,
		constant.HTTPRespKeyUserRole: userInfo.UserRole,
		constant.HTTPRespKeyUserOrg:  userInfo.UserOrganization,
	})
}

// LintingView This function return linting rule according to lint id
func (s LintingService) LintingView(ctx *gin.Context) {
	userInfo, exist := s.checkUserInContextExistsAndReturn(ctx)
	if !exist {
		ctx.Abort()
		return
	}

	lintId, exist := s.paramCheckAndReturn(ctx)
	if !exist {
		ctx.Abort()
		return
	}

	org, exist := s.checkOrganizationAndLintExistsAndReturn(ctx, userInfo.UserId, userInfo.UserOrganization, lintId)
	if !exist {
		ctx.Abort()
		return
	}

	bytes, err := org.ViewFile(ctx, lintId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			constant.HTTPRespKeyCode:    constant.HTTPRespCodeInternalError,
			constant.HTTPRespKeyMessage: err.Error(),
		})
		ctx.Abort()
		return
	}

	encodeStr := base64.StdEncoding.EncodeToString(bytes)

	ctx.JSON(200, gin.H{
		constant.HTTPRespKeyCode:     "linting view ok",
		constant.HTTPRespKeyUserId:   userInfo.UserId,
		constant.HTTPRespKeyUserRole: userInfo.UserRole,
		constant.HTTPRespKeyUserOrg:  userInfo.UserOrganization,
		constant.HTTPRespKeyUserFile: encodeStr,
	})
}

func (s LintingService) checkProjectIdAndReturn(ctx *gin.Context, org *configuration.Organization) (*configuration.Project, bool) {
	ctxId, _ := ctx.Get(constant.KeyContextId)
	projId := ctx.PostForm(constant.HTTPReqKeyProjIdInForm)

	//check param projId exists in form data
	if projId == "" {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   fmt.Sprintf("request param missing : [%s]", constant.HTTPReqKeyProjIdInForm),
			constant.LogKeyProjId:    projId,
			constant.LogKeyContextId: ctxId,
		}).Warn()
		ctx.JSON(http.StatusOK, gin.H{
			constant.HTTPRespKeyCode:    constant.HTTPRespCodeParamMissing,
			constant.HTTPRespKeyMessage: fmt.Sprintf(constant.HTTPRespMessageParamMissing, constant.HTTPReqKeyProjIdInForm),
		})
		return nil, false
	}

	//check projId exists in organization
	projMap := org.Projects
	for k, v := range projMap {
		if k == projId {
			return &configuration.Project{
				Active: v.Active,
				ProjId: k,
			}, true
		}
	}

	log.WithFields(log.Fields{
		constant.LogKeyMessage:   "projId not found in organization",
		constant.LogKeyProjId:    projId,
		constant.LogKeyContextId: ctxId,
	}).Warn()
	ctx.JSON(http.StatusOK, gin.H{
		constant.HTTPRespKeyCode:    constant.HTTPRespCodeProjNotFound,
		constant.HTTPRespKeyMessage: constant.HTTPRespMessageProjNotFound,
	})

	return nil, false
}

func (s LintingService) paramCheckAndReturn(ctx *gin.Context) (string, bool) {
	ctxId, _ := ctx.Get(constant.KeyContextId)
	lintId := ctx.Param(constant.HTTPReqKeyLintInPath)

	if lintId == "" {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   fmt.Sprintf("request param missing : [%s]", constant.HTTPReqKeyLintInPath),
			constant.LogKeyLintId:    lintId,
			constant.LogKeyContextId: ctxId,
		}).Warn()
		ctx.JSON(http.StatusOK, gin.H{
			constant.HTTPRespKeyCode:    constant.HTTPRespCodeParamMissing,
			constant.HTTPRespKeyMessage: fmt.Sprintf(constant.HTTPRespMessageParamMissing, constant.HTTPReqKeyLintInPath),
		})
		return "", false
	}

	return lintId, true
}

func (s LintingService) checkUserInContextExistsAndReturn(ctx *gin.Context) (*configuration.UserSession, bool) {
	ctxId, _ := ctx.Get(constant.KeyContextId)
	userInfo, exists := s.getUserInfoFromContext(ctx)
	if !exists {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   "fail to find userinfo in context",
			constant.LogKeyContextId: ctxId,
		}).Error()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			constant.LogKeyCode:    constant.ServiceInternalErrorCode,
			constant.LogKeyMessage: constant.ServiceInternalErrorMsg,
		})
		return nil, false
	}

	return userInfo, true
}

func (s LintingService) organizationExist(ctx *gin.Context, userId, orgId string) (*configuration.Organization, bool) {
	ctxId, _ := ctx.Get(constant.KeyContextId)
	mockupOrg, _ := configuration.GetMockUpOrganization(ctx)
	org, exist := mockupOrg.Organizations[orgId]
	if !exist {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   "fail to find organization",
			constant.LogKeyOrgId:     orgId,
			constant.LogKeyUserId:    userId,
			constant.LogKeyContextId: ctxId,
		}).Error()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			constant.LogKeyCode:    constant.ServiceInternalErrorCode,
			constant.LogKeyMessage: constant.ServiceInternalErrorMsg,
		})
		return nil, false
	}

	return org, true
}

func (s LintingService) checkOrganizationAndLintExistsAndReturn(ctx *gin.Context, userId, orgId, lintId string) (*configuration.Organization, bool) {
	ctxId, _ := ctx.Get(constant.KeyContextId)
	mockupOrg, _ := configuration.GetMockUpOrganization(ctx)
	org, exist := mockupOrg.Organizations[orgId]
	if !exist {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   "fail to find organization",
			constant.LogKeyOrgId:     orgId,
			constant.LogKeyUserId:    userId,
			constant.LogKeyLintId:    lintId,
			constant.LogKeyContextId: ctxId,
		}).Error()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			constant.LogKeyCode:    constant.ServiceInternalErrorCode,
			constant.LogKeyMessage: constant.ServiceInternalErrorMsg,
		})
		return nil, false
	}

	var lintExist = false
	for _, v := range org.LintingRule {
		if v == lintId {
			lintExist = true
			break
		}
	}

	if !lintExist {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   "fail to find lint in org",
			constant.LogKeyOrgId:     orgId,
			constant.LogKeyUserId:    userId,
			constant.LogKeyContextId: ctxId,
			constant.LogKeyLintId:    lintId,
		}).Error()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			constant.LogKeyCode:    constant.ServiceInternalErrorCode,
			constant.LogKeyMessage: constant.ServiceInternalErrorMsg,
		})
		return nil, false
	}

	return org, true
}

func (s LintingService) getUserInfoFromContext(ctx *gin.Context) (*configuration.UserSession, bool) {
	userInfo, exist := ctx.Get(constant.ContextKeyUserInfo)
	if exist {
		info := userInfo.(configuration.UserSession)
		return &info, true
	} else {
		return nil, false
	}
}
