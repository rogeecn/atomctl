package {{.PkgName}}

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

 	// 这里的依赖需要被导入，否则会报错
	_ "atom/providers"


	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/dig"
)

type {{.PascalName}}InjectParams struct {
	dig.In
}

type {{.PascalName}}Suite struct {
	suite.Suite
	{{.PascalName}}InjectParams
}

func Test_{{.PascalName}}Suite(t *testing.T) {
	err := container.Container.Invoke(func(p {{.PascalName}}InjectParams) {
		s := &{{.PascalName}}Suite{}
		s.{{.PascalName}}InjectParams = p

		suite.Run(t, s)
	})
	assert.NoError(t, err)
}

func (s *{{.PascalName}}Suite) SetupSuite() {
	fmt.Println("SetupSuite")
}

func (s *{{.PascalName}}Suite) SetupTest() {
	fmt.Println("SetupTest")
}
func (s *{{.PascalName}}Suite) BeforeTest(suiteName, testName string) {
	fmt.Println("BeforeTest:", suiteName, testName)
}
func (s *{{.PascalName}}Suite) AfterTest(suiteName, testName string) {
	fmt.Println("AfterTest:", suiteName, testName)
}
func (s *{{.PascalName}}Suite) HandleStats(suiteName string, stats *suite.SuiteInformation) {
	fmt.Println("HandleStats:", suiteName, stats)
}
func (s *{{.PascalName}}Suite) TearDownTest() {
	fmt.Println("TearDownTest")
}
func (s *{{.PascalName}}Suite) TearDownSuite() {
	fmt.Println("TearDownSuite")
}
///////////////////
// start testing cases
//////////////////

func (s *{{.PascalName}}Suite) Test_GetName() {

}