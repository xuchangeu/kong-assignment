package configuration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/xuchangeu/kong-assignment/src/constant"
)

type MockupOrganization struct {
	lock          sync.Mutex
	Organizations map[string]*Organization `json:"organization"`
}

type Organization struct {
	lock           sync.Mutex
	LintingRule    []string           `json:"linting-rule"`
	Projects       map[string]Project `json:"projects"`
	OrganizationId string             `json:"-"`
}

type Project struct {
	Active string `json:"active"`
	ProjId string `json:"-"`
}

var mockupOrganization *MockupOrganization

// GetMockUpOrganization Load all mockup organization into object
func GetMockUpOrganization(ctx context.Context) (*MockupOrganization, error) {
	ctxId := ctx.Value(constant.KeyContextId)
	if mockupOrganization == nil {
		data, err := retrieveMockupOrganization()
		if err != nil {
			log.WithFields(log.Fields{
				constant.LogKeyCode:      constant.RetrieveMockupOrgErrorCode,
				constant.LogKeyError:     err.Error(),
				constant.LogKeyContextId: fmt.Sprintf("%v", ctxId),
			}).Fatal()
			return nil, err
		}
		mockupOrganization = data
	}
	return mockupOrganization, nil
}

// GetOrganizationById Get Organization By Organization Id
func (m *MockupOrganization) GetOrganizationById(ctx context.Context, orgId string) (*Organization, error) {
	ctxId := ctx.Value(constant.KeyContextId)
	m.lock.Lock()
	defer func() {
		m.lock.Unlock()
	}()
	org, exist := m.Organizations[orgId]
	if !exist {
		log.WithFields(log.Fields{
			constant.LogKeyCode:      constant.GetOrganizationByIdNotExistErrorCode,
			constant.LogKeyMessage:   constant.GetOrganizationByIdNotExistErrorMsg,
			constant.LogKeyContextId: fmt.Sprintf("%s", ctxId),
		}).Error()
		return nil, errors.New(constant.GetOrganizationByIdNotExistErrorMsg)
	}
	return org, nil
}

func (m *MockupOrganization) writeMockupOrganization(ctx context.Context) error {
	ctxId := ctx.Value(constant.KeyContextId)
	m.lock.Lock()
	log.WithFields(log.Fields{
		constant.LogKeyMessage:   "write mockup organization start, do lock",
		constant.LogKeyContextId: fmt.Sprintf("%s", ctxId),
	}).Info()
	defer func() {
		m.lock.Unlock()
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   "write mockup organization end, do unlock",
			constant.LogKeyContextId: fmt.Sprintf("%s", ctxId),
		}).Info()
	}()
	var path, err = getMockupFilePath()
	if err != nil {
		return err
	}
	replaceBytes, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, replaceBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (org *Organization) WriteRule(ctx context.Context, orgId string, lintingId string, rule io.Reader) error {
	ctxId := ctx.Value(constant.KeyContextId)
	org.lock.Lock()
	defer org.lock.Unlock()

	filePath := path.Join(fmt.Sprintf("%s.yaml", lintingId))
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	if err != nil {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   "open linting file error",
			constant.LogKeyContextId: ctxId,
			"file-path":              filePath,
			constant.LogKeyError:     err.Error(),
		}).Error()
		return err
	}

	bytes, err := io.ReadAll(rule)
	if err != nil {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   "write linting file error",
			constant.LogKeyContextId: ctxId,
			"file-path":              filePath,
			constant.LogKeyError:     err.Error(),
		}).Error()
		return err
	}

	n, err := file.Write(bytes)
	if err != nil {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   "open linting file error",
			constant.LogKeyContextId: ctxId,
			"file-path":              filePath,
			constant.LogKeyError:     err.Error(),
		}).Error()
		return err
	}

	log.WithFields(log.Fields{
		constant.LogKeyMessage:   "write linting file success",
		constant.LogKeyContextId: ctxId,
		"file-path":              filePath,
		"write-bytes":            n,
	}).Info()

	//only write when linting id is empty
	if lintingId == "" {
		config, _ := GetMockUpOrganization(ctx)
		linting := config.Organizations[orgId].LintingRule
		linting = append(linting, lintingId)
		config.Organizations[orgId].LintingRule = linting
		err = config.writeMockupOrganization(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				constant.LogKeyMessage:   "write configuration error",
				constant.LogKeyContextId: ctxId,
				"file-path":              filePath,
				constant.LogKeyError:     err.Error(),
			}).Error()
		}
	}
	return nil
}

func (org *Organization) ViewFile(ctx context.Context, lintingId string) ([]byte, error) {
	ctxId := ctx.Value(constant.KeyContextId)
	filePath := path.Join(fmt.Sprintf("%s.yaml", lintingId))
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	if err != nil {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   "open linting file error",
			constant.LogKeyContextId: ctxId,
			"file-path":              filePath,
			constant.LogKeyError:     err.Error(),
		}).Error()
		return nil, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   "read linting file error",
			constant.LogKeyContextId: ctxId,
			"file-path":              filePath,
			constant.LogKeyError:     err.Error(),
		}).Error()
		return nil, err
	}

	return bytes, nil
}

// ApplyLintToProj Apply Lint rule to project
func (org *Organization) ApplyLintToProj(ctx context.Context, orgId string, projId string, ruleId string) error {
	ctxId := ctx.Value(constant.KeyContextId)
	org.lock.Lock()
	log.WithFields(log.Fields{
		constant.LogKeyMessage:   "apply lint to project start, do lock",
		constant.LogKeyOrgId:     orgId,
		constant.LogKeyProjId:    projId,
		constant.LogKeyContextId: fmt.Sprintf("%s", ctxId),
	}).Info()
	defer func() {
		org.lock.Unlock()
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   "apply lint to project end, do unlock",
			constant.LogKeyContextId: fmt.Sprintf("%s", ctxId),
		}).Info()
	}()
	org, exist := mockupOrganization.Organizations[orgId]
	if !exist {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   "apply rule fail, organization not exists",
			constant.LogKeyOrgId:     orgId,
			constant.LogKeyProjId:    projId,
			constant.LogKeyContextId: fmt.Sprintf("%s", ctxId),
		}).Error()
		return errors.New(constant.ApplyRuleOrganizationNotExistErrorMsg)
	}
	lastProject, exists := org.Projects[projId]
	if !exists {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   "apply rule fail, project not exists",
			constant.LogKeyOrgId:     orgId,
			constant.LogKeyProjId:    projId,
			constant.LogKeyContextId: fmt.Sprintf("%s", ctxId),
		}).Error()
		return errors.New(constant.ApplyRuleProjectNotExistErrorMsg)
	}
	org.Projects[projId] = Project{
		Active: ruleId,
	}

	mockup, _ := GetMockUpOrganization(ctx)
	err := mockup.writeMockupOrganization(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			constant.LogKeyMessage:   constant.ApplyRuleWriteMockupDataErrorMsg,
			constant.LogKeyOrgId:     orgId,
			constant.LogKeyProjId:    projId,
			constant.LogKeyContextId: fmt.Sprintf("%s", ctxId),
			constant.LogKeyCode:      constant.ApplyRuleWriteMockupDataErrorCode,
		}).Error()
		org.Projects[projId] = lastProject
		return errors.New(constant.ApplyRuleWriteMockupDataErrorMsg)
	}

	return nil
}

func getMockupFilePath() (string, error) {
	var path string
	if os.Getenv("MOCK_ORG") != "" {
		path = os.Getenv("MOCK_ORG")
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			return path, errors.New(constant.RetrieveMockupSessionErrorMsg)
		}
		//hardcode relation path for dev & test, while program should read system env on prod stage.
		path = filepath.Join(pwd, "src", "conf", "organization.json")
	}
	return path, nil
}

func retrieveMockupOrganization() (*MockupOrganization, error) {

	var data MockupOrganization
	var path, err = getMockupFilePath()
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		constant.LogKeyMessage: "retrieving mockup organization",
		constant.LogKeyPath:    path,
	}).Info()

	bytesData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytesData, &data)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &data, nil
}
