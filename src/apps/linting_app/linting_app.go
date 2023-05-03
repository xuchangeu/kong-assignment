package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuchangeu/kong-assignment/src/constant"
	"github.com/xuchangeu/kong-assignment/src/service"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {

	router := gin.Default()

	// //these are mockup service which skip auth middleware
	router.GET("/user/auth/:userid", service.UserAuth)
	router.GET("/user/project/:userid", service.UserProject)

	// //function to apply linting, include middleware to auth user
	router.POST(fmt.Sprintf("/linting/apply/:%s", constant.HTTPReqKeyLintInPath), service.AuthMiddleWare(), service.LintingService{}.LintingApply)
	router.POST(fmt.Sprintf("/linting/create/:%s", constant.HTTPReqKeyLintInPath), service.AuthMiddleWare(), service.LintingService{}.LintingCreate)
	router.GET(fmt.Sprintf("/linting/view/:%s", constant.HTTPReqKeyLintInPath), service.AuthMiddleWare2(), service.LintingService{}.LintingView)
	err := router.Run()
	if err != nil {
		log.Fatalf("run gin service error : %s, exit", err.Error())
		os.Exit(1)
	}

}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
}
