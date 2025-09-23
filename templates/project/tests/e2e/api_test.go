package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"{{.ModuleName}}/app"
	"{{.ModuleName}}/app/config"
	. "github.com/smartystreets/goconvey/convey"
)

// TestAPIHealth 测试 API 健康检查
func TestAPIHealth(t *testing.T) {
	Convey("API 健康检查测试", t, func() {
		var server *httptest.Server
		var testConfig *config.Config

		Convey("当启动测试服务器时", func() {
			testConfig = &config.Config{
				App: config.AppConfig{
					Mode:    "test",
					BaseURI: "http://localhost:8080",
				},
				Http: config.HttpConfig{
					Port: 8080,
				},
				Log: config.LogConfig{
					Level:        "debug",
					Format:       "text",
					EnableCaller: true,
				},
			}

			app := app.New(testConfig)
			server = httptest.NewServer(app)

			Convey("服务器应该成功启动", func() {
				So(server, ShouldNotBeNil)
				So(server.URL, ShouldNotBeEmpty)
			})
		})

		Convey("当访问健康检查端点时", func() {
			resp, err := http.Get(server.URL + "/health")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)

			defer resp.Body.Close()

			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			So(err, ShouldBeNil)

			Convey("响应应该包含正确的状态", func() {
				So(result["status"], ShouldEqual, "ok")
			})

			Convey("响应应该包含时间戳", func() {
				So(result, ShouldContainKey, "timestamp")
			})

			Convey("响应应该是 JSON 格式", func() {
				So(resp.Header.Get("Content-Type"), ShouldEqual, "application/json; charset=utf-8")
			})
		})

		Convey("当访问不存在的端点时", func() {
			resp, err := http.Get(server.URL + "/api/nonexistent")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusNotFound)

			defer resp.Body.Close()

			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			So(err, ShouldBeNil)

			Convey("响应应该包含错误信息", func() {
				So(result, ShouldContainKey, "error")
			})
		})

		Convey("当测试 CORS 支持", func() {
			req, err := http.NewRequest("OPTIONS", server.URL+"/api/test", nil)
			So(err, ShouldBeNil)

			req.Header.Set("Origin", "http://localhost:3000")
			req.Header.Set("Access-Control-Request-Method", "POST")
			req.Header.Set("Access-Control-Request-Headers", "Content-Type,Authorization")

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			defer resp.Body.Close()

			Convey("应该返回正确的 CORS 头", func() {
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(resp.Header.Get("Access-Control-Allow-Origin"), ShouldContainSubstring, "localhost")
				So(resp.Header.Get("Access-Control-Allow-Methods"), ShouldContainSubstring, "POST")
			})
		})

		Reset(func() {
			if server != nil {
				server.Close()
			}
		})
	})
}

// TestAPIPerformance 测试 API 性能
func TestAPIPerformance(t *testing.T) {
	Convey("API 性能测试", t, func() {
		var server *httptest.Server
		var testConfig *config.Config

		Convey("当准备性能测试时", func() {
			testConfig = &config.Config{
				App: config.AppConfig{
					Mode:    "test",
					BaseURI: "http://localhost:8080",
				},
				Http: config.HttpConfig{
					Port: 8080,
				},
				Log: config.LogConfig{
					Level:  "error", // 减少日志输出以提升性能
					Format: "text",
				},
			}

			app := app.New(testConfig)
			server = httptest.NewServer(app)
		})

		Convey("当测试响应时间时", func() {
			start := time.Now()
			resp, err := http.Get(server.URL + "/health")
			So(err, ShouldBeNil)
			defer resp.Body.Close()

			duration := time.Since(start)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)

			Convey("响应时间应该在合理范围内", func() {
				So(duration, ShouldBeLessThan, 100*time.Millisecond)
			})
		})

		Convey("当测试并发请求时", func() {
			const numRequests = 50
			const maxConcurrency = 10
			const timeout = 5 * time.Second

			var wg sync.WaitGroup
			successCount := 0
			errorCount := 0
			var mu sync.Mutex

			// 使用信号量控制并发数
			sem := make(chan struct{}, maxConcurrency)

			start := time.Now()

			for i := 0; i < numRequests; i++ {
				wg.Add(1)
				go func(requestID int) {
					defer wg.Done()

					// 获取信号量
					sem <- struct{}{}
					defer func() { <-sem }()

					ctx, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()

					req, err := http.NewRequestWithContext(ctx, "GET", server.URL+"/health", nil)
					if err != nil {
						mu.Lock()
						errorCount++
						mu.Unlock()
						return
					}

					client := &http.Client{
						Timeout: timeout,
					}

					resp, err := client.Do(req)
					if err != nil {
						mu.Lock()
						errorCount++
						mu.Unlock()
						return
					}
					defer resp.Body.Close()

					mu.Lock()
					if resp.StatusCode == http.StatusOK {
						successCount++
					} else {
						errorCount++
					}
					mu.Unlock()
				}(i)
			}

			wg.Wait()
			duration := time.Since(start)

			Convey("所有请求都应该完成", func() {
				So(successCount+errorCount, ShouldEqual, numRequests)
			})

			Convey("所有请求都应该成功", func() {
				So(errorCount, ShouldEqual, 0)
			})

			Convey("总耗时应该在合理范围内", func() {
				So(duration, ShouldBeLessThan, 10*time.Second)
			})

			Convey("并发性能应该良好", func() {
				avgTime := duration / numRequests
				So(avgTime, ShouldBeLessThan, 200*time.Millisecond)
			})
		})

		Reset(func() {
			if server != nil {
				server.Close()
			}
		})
	})
}

// TestAPIBehavior 测试 API 行为
func TestAPIBehavior(t *testing.T) {
	Convey("API 行为测试", t, func() {
		var server *httptest.Server
		var testConfig *config.Config

		Convey("当准备行为测试时", func() {
			testConfig = &config.Config{
				App: config.AppConfig{
					Mode:    "test",
					BaseURI: "http://localhost:8080",
				},
				Http: config.HttpConfig{
					Port: 8080,
				},
				Log: config.LogConfig{
					Level:        "debug",
					Format:       "text",
					EnableCaller: true,
				},
			}

			app := app.New(testConfig)
			server = httptest.NewServer(app)
		})

		Convey("当测试不同 HTTP 方法时", func() {
			testURL := server.URL + "/health"

			Convey("GET 请求应该成功", func() {
				resp, err := http.Get(testURL)
				So(err, ShouldBeNil)
				defer resp.Body.Close()
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
			})

			Convey("POST 请求应该被处理", func() {
				resp, err := http.Post(testURL, "application/json", bytes.NewBuffer([]byte{}))
				So(err, ShouldBeNil)
				defer resp.Body.Close()
				// 健康检查端点通常支持所有方法
				So(resp.StatusCode, ShouldBeIn, []int{http.StatusOK, http.StatusMethodNotAllowed})
			})

			Convey("PUT 请求应该被处理", func() {
				req, err := http.NewRequest("PUT", testURL, bytes.NewBuffer([]byte{}))
				So(err, ShouldBeNil)
				resp, err := http.DefaultClient.Do(req)
				So(err, ShouldBeNil)
				defer resp.Body.Close()
				So(resp.StatusCode, ShouldBeIn, []int{http.StatusOK, http.StatusMethodNotAllowed})
			})

			Convey("DELETE 请求应该被处理", func() {
				req, err := http.NewRequest("DELETE", testURL, nil)
				So(err, ShouldBeNil)
				resp, err := http.DefaultClient.Do(req)
				So(err, ShouldBeNil)
				defer resp.Body.Close()
				So(resp.StatusCode, ShouldBeIn, []int{http.StatusOK, http.StatusMethodNotAllowed})
			})
		})

		Convey("当测试自定义请求头时", func() {
			req, err := http.NewRequest("GET", server.URL+"/health", nil)
			So(err, ShouldBeNil)

			// 设置各种请求头
			req.Header.Set("User-Agent", "E2E-Test-Agent/1.0")
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
			req.Header.Set("X-Custom-Header", "test-value")
			req.Header.Set("X-Request-ID", "test-request-123")
			req.Header.Set("Authorization", "Bearer test-token")

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			defer resp.Body.Close()

			Convey("请求应该成功", func() {
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
			})

			Convey("响应应该是 JSON 格式", func() {
				So(resp.Header.Get("Content-Type"), ShouldEqual, "application/json; charset=utf-8")
			})
		})

		Convey("当测试错误处理时", func() {
			Convey("访问不存在的路径应该返回 404", func() {
				resp, err := http.Get(server.URL + "/api/v1/nonexistent")
				So(err, ShouldBeNil)
				defer resp.Body.Close()
				So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
			})

			Convey("访问非法路径应该返回 404", func() {
				resp, err := http.Get(server.URL + "/../etc/passwd")
				So(err, ShouldBeNil)
				defer resp.Body.Close()
				So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
			})
		})

		Reset(func() {
			if server != nil {
				server.Close()
			}
		})
	})
}

// TestAPIDocumentation 测试 API 文档
func TestAPIDocumentation(t *testing.T) {
	Convey("API 文档测试", t, func() {
		var server *httptest.Server
		var testConfig *config.Config

		Convey("当准备文档测试时", func() {
			testConfig = &config.Config{
				App: config.AppConfig{
					Mode:    "test",
					BaseURI: "http://localhost:8080",
				},
				Http: config.HttpConfig{
					Port: 8080,
				},
				Log: config.LogConfig{
					Level:        "debug",
					Format:       "text",
					EnableCaller: true,
				},
			}

			app := app.New(testConfig)
			server = httptest.NewServer(app)
		})

		Convey("当访问 Swagger UI 时", func() {
			resp, err := http.Get(server.URL + "/swagger/index.html")
			So(err, ShouldBeNil)
			defer resp.Body.Close()

			Convey("应该能够访问 Swagger UI", func() {
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
			})

			Convey("响应应该是 HTML 格式", func() {
				contentType := resp.Header.Get("Content-Type")
				So(contentType, ShouldContainSubstring, "text/html")
			})
		})

		Convey("当访问 OpenAPI 规范时", func() {
			resp, err := http.Get(server.URL + "/swagger/doc.json")
			So(err, ShouldBeNil)
			defer resp.Body.Close()

			Convey("应该能够访问 OpenAPI 规范", func() {
				// 如果存在则返回 200，不存在则返回 404
				So(resp.StatusCode, ShouldBeIn, []int{http.StatusOK, http.StatusNotFound})
			})

			Convey("如果存在，响应应该是 JSON 格式", func() {
				if resp.StatusCode == http.StatusOK {
					contentType := resp.Header.Get("Content-Type")
					So(contentType, ShouldContainSubstring, "application/json")
				}
			})
		})

		Reset(func() {
			if server != nil {
				server.Close()
			}
		})
	})
}