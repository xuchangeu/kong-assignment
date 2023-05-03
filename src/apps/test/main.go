package main

import (
	"fmt"
	"sync"
	"time"
)

type MockupOrganization struct {
	Organizations map[string]Organization `json:"organization"`
	lock          sync.Mutex
}

type Organization struct {
	LintingRule []string           `json:"linting-rule"`
	Projects    map[string]Project `json:"projects"`
}

type Project struct {
	Active string `json:"active"`
}

var orgData *MockupOrganization

var orgId = "org-111"
var projId = "proj-222"

func main() {
	time.Sleep(2 * time.Second)
}

func init() {
	orgData = &MockupOrganization{
		Organizations: map[string]Organization{
			orgId: {
				LintingRule: []string{"r1", "r2", "r3"},
				Projects: map[string]Project{
					projId: {
						Active: "r1",
					},
				},
			},
		},
	}

	var immutOrg = map[string]Organization{}
	for k, v := range orgData.Organizations {
		immutOrg[k] = v
	}
	var immut = MockupOrganization{
		Organizations: immutOrg,
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		orgData.Organizations[orgId].Projects[projId] = Project{
			Active: "r100",
		}
		proj := orgData.Organizations[orgId].Projects[projId]
		projMap := orgData.Organizations[orgId].Projects

		fmt.Printf("%T,%p\n", proj, &proj)
		fmt.Printf("%T,%p\n", projMap, &projMap)
	}()

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println(immut.Organizations[orgId].Projects[projId].Active)
		proj := immut.Organizations[orgId].Projects[projId]
		projMap := immut.Organizations[orgId].Projects
		fmt.Printf("%p\n", &proj)
		fmt.Printf("%p\n", &projMap)
	}()

}
