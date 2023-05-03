package configuration

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/xuchangeu/kong-assignment/src/constant"

	log "github.com/sirupsen/logrus"
)

type AppConfiguration struct {
}

/**mockup data structure**/
type MockupSession struct {
	Session map[string]UserSession `json:"session"`
}

type UserSession struct {
	UserId           string `json:"user"`
	UserRole         string `json:"role"`
	UserOrganization string `json:"organization"`
}

/**mockup data structure end**/

var conf *AppConfiguration
var mockup *MockupSession

func GetConfiguration() *AppConfiguration {
	if conf == nil {
		c, err := retrieveConfiguration()
		if err != nil {

			os.Exit(constant.RetrieveAppConfigError)
		}
		conf = c
	}
	return conf
}

func retrieveConfiguration() (conf *AppConfiguration, err error) {
	return &AppConfiguration{}, nil
}

func GetMockUpSession() *MockupSession {
	if mockup == nil {
		data, err := retrieveMockupSession()
		if err != nil {
			log.WithFields(log.Fields{
				constant.KeyCode:  constant.RetrieveMockupSessionErrorCode,
				constant.KeyError: err.Error(),
			}).Fatal()
			os.Exit(constant.RetrieveMockupSessionErrorCode)
		}
		mockup = data
	}
	return mockup
}

func retrieveMockupSession() (*MockupSession, error) {
	var path string
	var data MockupSession
	if os.Getenv("MOCK_SESSION") != "" {
		path = os.Getenv("MOCK_SESSION")
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, errors.New(constant.RetrieveMockupSessionErrorMsg)
		}
		//hardcode relation path for dev & test, while program should read system env on prod stage.
		path = filepath.Join(pwd, "src", "conf", "session.json")
	}
	log.WithFields(log.Fields{
		"message": "retrieving mockup session",
		"path":    path,
	}).Info()

	bytesData, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New(constant.RetrieveMockupSessionErrorMsg)
	}

	err = json.Unmarshal(bytesData, &data)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New(constant.RetrieveMockupDataUnmarshalErrorMsg)
	}

	return &data, nil
}
