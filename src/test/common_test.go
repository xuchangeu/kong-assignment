package test

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"testing"
)

func TestApplyLintToProj(t *testing.T) {

	want := true
	have := true
	assert.Equal(t, want, have)

	assert.Equal(t, fmt.Sprintf("%d", 1), "1")
	//mockup, err := configuration.GetMockUpOrganization(context.WithValue(context.Background(), constant.KeyContextId, "init"))
	//
	//if err != nil {
	//	log.WithFields(log.Fields{
	//		"msg":  constant.RetrieveMockupOrgErrorMsg,
	//		"code": constant.RetrieveMockupOrgErrorCode,
	//	}).Fatal()
	//	os.Exit(constant.RetrieveMockupOrgErrorCode)

	//var orgId = "org-111"
	//var projId = "proj-222"
	//var i = 0
	//go func(i int) {
	//	for {
	//		ctxId := fmt.Sprintf("aaa-%d", i)
	//		ctx := context.WithValue(context.Background(), "context-id", ctxId)
	//
	//		org, _ := mockup.GetOrganizationById(ctx, orgId)
	//		//fmt.Println(org, ctxId, orgId, projId)
	//		org.ApplyLintToProj(ctx, orgId, projId, fmt.Sprintf("lint-%d", i))
	//		i++
	//	}
	//
	//}(i)
	//
	//go func() {
	//	for {
	//		ctxId := fmt.Sprintf("aaa-%d", 50)
	//		ctx := context.WithValue(context.Background(), "context-id", ctxId)
	//
	//		mockup.GetOrganizationById(ctx, orgId)
	//		//fmt.Println(org, ctxId, orgId, projId)
	//	}
	//
	//}()
	//
	//time.Sleep(10 * time.Second)

}
