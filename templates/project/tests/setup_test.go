package tests

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestMain 测试入口点
func TestMain(m *testing.M) {
	// 运行测试
	m.Run()
}

// TestSetup 测试基础设置
func TestSetup(t *testing.T) {
	Convey("测试基础设置", t, func() {
		Convey("当初始化测试环境时", func() {
			// 初始化测试环境
			testEnv := &TestEnvironment{
				Name:    "test-env",
				Version: "1.0.0",
			}

			Convey("那么测试环境应该被正确创建", func() {
				So(testEnv.Name, ShouldEqual, "test-env")
				So(testEnv.Version, ShouldEqual, "1.0.0")
			})
		})
	})
}

// TestEnvironment 测试环境结构
type TestEnvironment struct {
	Name    string
	Version string
	Config  map[string]interface{}
}

// NewTestEnvironment 创建新的测试环境
func NewTestEnvironment(name string) *TestEnvironment {
	return &TestEnvironment{
		Name:   name,
		Config: make(map[string]interface{}),
	}
}

// WithConfig 设置配置
func (e *TestEnvironment) WithConfig(key string, value interface{}) *TestEnvironment {
	e.Config[key] = value
	return e
}

// GetConfig 获取配置
func (e *TestEnvironment) GetConfig(key string) interface{} {
	return e.Config[key]
}

// Setup 设置测试环境
func (e *TestEnvironment) Setup() *TestEnvironment {
	// 初始化测试环境
	e.Config["initialized"] = true
	return e
}

// Cleanup 清理测试环境
func (e *TestEnvironment) Cleanup() {
	// 清理测试环境
	e.Config = make(map[string]interface{})
}

// TestEnvironmentManagement 测试环境管理
func TestEnvironmentManagement(t *testing.T) {
	Convey("测试环境管理", t, func() {
		var env *TestEnvironment

		Convey("当创建新测试环境时", func() {
			env = NewTestEnvironment("test-app")

			Convey("那么环境应该有正确的名称", func() {
				So(env.Name, ShouldEqual, "test-app")
				So(env.Version, ShouldBeEmpty)
				So(env.Config, ShouldNotBeNil)
			})
		})

		Convey("当设置配置时", func() {
			env.WithConfig("debug", true)
			env.WithConfig("port", 8080)

			Convey("那么配置应该被正确设置", func() {
				So(env.GetConfig("debug"), ShouldEqual, true)
				So(env.GetConfig("port"), ShouldEqual, 8080)
			})
		})

		Convey("当初始化环境时", func() {
			env.Setup()

			Convey("那么环境应该被标记为已初始化", func() {
				So(env.GetConfig("initialized"), ShouldEqual, true)
			})
		})

		Reset(func() {
			if env != nil {
				env.Cleanup()
			}
		})
	})
}

// TestConveyBasicUsage 测试 Convey 基础用法
func TestConveyBasicUsage(t *testing.T) {
	Convey("Convey 基础用法测试", t, func() {
		Convey("数字操作", func() {
			num := 42

			Convey("应该能够进行基本比较", func() {
				So(num, ShouldEqual, 42)
				So(num, ShouldBeGreaterThan, 0)
				So(num, ShouldBeLessThan, 100)
			})
		})

		Convey("字符串操作", func() {
			str := "hello world"

			Convey("应该能够进行字符串比较", func() {
				So(str, ShouldEqual, "hello world")
				So(str, ShouldContainSubstring, "hello")
				So(str, ShouldStartWith, "hello")
				So(str, ShouldEndWith, "world")
			})
		})

		Convey("切片操作", func() {
			slice := []int{1, 2, 3, 4, 5}

			Convey("应该能够进行切片操作", func() {
				So(slice, ShouldHaveLength, 5)
				So(slice, ShouldContain, 3)
				So(slice, ShouldNotContain, 6)
			})
		})

		Convey("Map 操作", func() {
			m := map[string]interface{}{
				"name":  "test",
				"value": 123,
			}

			Convey("应该能够进行 Map 操作", func() {
				So(m, ShouldContainKey, "name")
				So(m, ShouldContainKey, "value")
				So(m["name"], ShouldEqual, "test")
				So(m["value"], ShouldEqual, 123)
			})
		})
	})
}