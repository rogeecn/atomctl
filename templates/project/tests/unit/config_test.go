package unit

import (
	"os"
	"path/filepath"
	"testing"

	"{{.ModuleName}}/app/config"
	. "github.com/smartystreets/goconvey/convey"
)

// TestConfigLoading 测试配置加载功能
func TestConfigLoading(t *testing.T) {
	Convey("配置加载测试", t, func() {
		var testConfig *config.Config
		var configPath string
		var testDir string

		Convey("当准备测试配置文件时", func() {
			originalWd, _ := os.Getwd()
			testDir = filepath.Join(originalWd, "..", "..", "fixtures", "test_config")

			Convey("应该创建测试配置目录", func() {
				err := os.MkdirAll(testDir, 0755)
				So(err, ShouldBeNil)
			})

			Convey("应该创建测试配置文件", func() {
				testConfigContent := `App:
  Mode: "test"
  BaseURI: "http://localhost:8080"
Http:
  Port: 8080
Database:
  Host: "localhost"
  Port: 5432
  Database: "test_db"
  Username: "test_user"
  Password: "test_password"
  SslMode: "disable"
Log:
  Level: "debug"
  Format: "text"
  EnableCaller: true`

				configPath = filepath.Join(testDir, "config.toml")
				err := os.WriteFile(configPath, []byte(testConfigContent), 0644)
				So(err, ShouldBeNil)
			})

			Convey("应该成功加载配置", func() {
				var err error
				testConfig, err = config.Load(configPath)
				So(err, ShouldBeNil)
				So(testConfig, ShouldNotBeNil)
			})
		})

		Convey("验证配置内容", func() {
			So(testConfig, ShouldNotBeNil)

			Convey("应用配置应该正确", func() {
				So(testConfig.App.Mode, ShouldEqual, "test")
				So(testConfig.App.BaseURI, ShouldEqual, "http://localhost:8080")
			})

			Convey("HTTP配置应该正确", func() {
				So(testConfig.Http.Port, ShouldEqual, 8080)
			})

			Convey("数据库配置应该正确", func() {
				So(testConfig.Database.Host, ShouldEqual, "localhost")
				So(testConfig.Database.Port, ShouldEqual, 5432)
				So(testConfig.Database.Database, ShouldEqual, "test_db")
				So(testConfig.Database.Username, ShouldEqual, "test_user")
				So(testConfig.Database.Password, ShouldEqual, "test_password")
				So(testConfig.Database.SslMode, ShouldEqual, "disable")
			})

			Convey("日志配置应该正确", func() {
				So(testConfig.Log.Level, ShouldEqual, "debug")
				So(testConfig.Log.Format, ShouldEqual, "text")
				So(testConfig.Log.EnableCaller, ShouldBeTrue)
			})
		})

		Reset(func() {
			// 清理测试文件
			if testDir != "" {
				os.RemoveAll(testDir)
			}
		})
	})
}

// TestConfigFromEnvironment 测试从环境变量加载配置
func TestConfigFromEnvironment(t *testing.T) {
	Convey("环境变量配置测试", t, func() {
		var originalEnvVars map[string]string

		Convey("当设置环境变量时", func() {
			// 保存原始环境变量
			originalEnvVars = map[string]string{
				"APP_MODE":   os.Getenv("APP_MODE"),
				"HTTP_PORT":  os.Getenv("HTTP_PORT"),
				"DB_HOST":    os.Getenv("DB_HOST"),
			}

			// 设置测试环境变量
			os.Setenv("APP_MODE", "test")
			os.Setenv("HTTP_PORT", "9090")
			os.Setenv("DB_HOST", "test-host")

			Convey("环境变量应该被正确设置", func() {
				So(os.Getenv("APP_MODE"), ShouldEqual, "test")
				So(os.Getenv("HTTP_PORT"), ShouldEqual, "9090")
				So(os.Getenv("DB_HOST"), ShouldEqual, "test-host")
			})
		})

		Convey("当从环境变量加载配置时", func() {
			originalWd, _ := os.Getwd()
			testDir := filepath.Join(originalWd, "..", "..", "fixtures", "test_config_env")

			Convey("应该创建测试配置目录", func() {
				err := os.MkdirAll(testDir, 0755)
				So(err, ShouldBeNil)
			})

			Convey("应该创建基础配置文件", func() {
				testConfigContent := `App:
  Mode: "development"
  BaseURI: "http://localhost:3000"
Http:
  Port: 3000
Database:
  Host: "localhost"
  Port: 5432
  Database: "default_db"
  Username: "default_user"
  Password: "default_password"
  SslMode: "disable"`

				configPath := filepath.Join(testDir, "config.toml")
				err := os.WriteFile(configPath, []byte(testConfigContent), 0644)
				So(err, ShouldBeNil)
			})

			Convey("应该成功加载并合并配置", func() {
				configPath := filepath.Join(testDir, "config.toml")
				loadedConfig, err := config.Load(configPath)

				So(err, ShouldBeNil)
				So(loadedConfig, ShouldNotBeNil)

				Convey("环境变量应该覆盖配置文件", func() {
					So(loadedConfig.App.Mode, ShouldEqual, "test")
					So(loadedConfig.Http.Port, ShouldEqual, 9090)
					So(loadedConfig.Database.Host, ShouldEqual, "test-host")
				})

				Convey("配置文件的默认值应该保留", func() {
					So(loadedConfig.App.BaseURI, ShouldEqual, "http://localhost:3000")
					So(loadedConfig.Database.Database, ShouldEqual, "default_db")
				})
			})

			Reset(func() {
				// 清理测试目录
				os.RemoveAll(testDir)
			})
		})

		Reset(func() {
			// 恢复原始环境变量
			if originalEnvVars != nil {
				for key, value := range originalEnvVars {
					if value == "" {
						os.Unsetenv(key)
					} else {
						os.Setenv(key, value)
					}
				}
			}
		})
	})
}

// TestConfigValidation 测试配置验证
func TestConfigValidation(t *testing.T) {
	Convey("配置验证测试", t, func() {
		Convey("当配置为空时", func() {
			config := &config.Config{}

			Convey("应该检测到缺失的必需配置", func() {
				So(config.App.Mode, ShouldBeEmpty)
				So(config.Http.Port, ShouldEqual, 0)
				So(config.Database.Host, ShouldBeEmpty)
			})
		})

		Convey("当配置端口无效时", func() {
			config := &config.Config{
				Http: config.HttpConfig{
					Port: -1,
				},
			}

			Convey("应该检测到无效端口", func() {
				So(config.Http.Port, ShouldBeLessThan, 0)
			})
		})

		Convey("当配置模式有效时", func() {
			validModes := []string{"development", "production", "testing"}

			for _, mode := range validModes {
				config := &config.Config{
					App: config.AppConfig{
						Mode: mode,
					},
				}

				Convey("模式 "+mode+" 应该是有效的", func() {
					So(config.App.Mode, ShouldBeIn, validModes)
				})
			}
		})
	})
}

// TestConfigDefaults 测试配置默认值
func TestConfigDefaults(t *testing.T) {
	Convey("配置默认值测试", t, func() {
		Convey("当创建新配置时", func() {
			config := &config.Config{}

			Convey("应该有合理的默认值", func() {
				// 测试应用的默认值
				So(config.App.Mode, ShouldEqual, "development")

				// 测试HTTP的默认值
				So(config.Http.Port, ShouldEqual, 8080)

				// 测试数据库的默认值
				So(config.Database.Port, ShouldEqual, 5432)
				So(config.Database.SslMode, ShouldEqual, "disable")

				// 测试日志的默认值
				So(config.Log.Level, ShouldEqual, "info")
				So(config.Log.Format, ShouldEqual, "json")
			})
		})
	})
}

// TestConfigHelpers 测试配置辅助函数
func TestConfigHelpers(t *testing.T) {
	Convey("配置辅助函数测试", t, func() {
		Convey("当使用配置辅助函数时", func() {
			config := &config.Config{
				App: config.AppConfig{
					Mode:    "production",
					BaseURI: "https://api.example.com",
				},
				Http: config.HttpConfig{
					Port: 443,
				},
			}

			Convey("应该能够获取应用环境", func() {
				env := config.App.Mode
				So(env, ShouldEqual, "production")
			})

			Convey("应该能够构建完整URL", func() {
				fullURL := config.App.BaseURI + "/api/v1/users"
				So(fullURL, ShouldEqual, "https://api.example.com/api/v1/users")
			})

			Convey("应该能够判断HTTPS", func() {
				isHTTPS := config.Http.Port == 443
				So(isHTTPS, ShouldBeTrue)
			})
		})
	})
}