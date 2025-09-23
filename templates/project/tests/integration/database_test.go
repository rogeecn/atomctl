package integration

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"{{.ModuleName}}/app/config"
	"{{.ModuleName}}/app/database"
	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestUser 测试用户模型
type TestUser struct {
	ID        int       `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;unique;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TestDatabaseConnection 测试数据库连接
func TestDatabaseConnection(t *testing.T) {
	Convey("数据库连接测试", t, func() {
		var db *gorm.DB
		var sqlDB *sql.DB
		var testConfig *config.Config
		var testDBName string

		Convey("当准备测试数据库时", func() {
			testDBName = "{{.ProjectName}}_test_integration"
			testConfig = &config.Config{
				Database: config.DatabaseConfig{
					Host:            "localhost",
					Port:            5432,
					Database:        testDBName,
					Username:        "postgres",
					Password:        "password",
					SslMode:         "disable",
					MaxIdleConns:    5,
					MaxOpenConns:    20,
					ConnMaxLifetime: 30 * time.Minute,
				},
			}

			Convey("应该能够连接到数据库", func() {
				dsn := testConfig.Database.GetDSN()
				var err error
				db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
				So(err, ShouldBeNil)
				So(db, ShouldNotBeNil)

				sqlDB, err = db.DB()
				So(err, ShouldBeNil)
				So(sqlDB, ShouldNotBeNil)

				// 设置连接池
				sqlDB.SetMaxIdleConns(testConfig.Database.MaxIdleConns)
				sqlDB.SetMaxOpenConns(testConfig.Database.MaxOpenConns)
				sqlDB.SetConnMaxLifetime(testConfig.Database.ConnMaxLifetime)

				// 测试连接
				err = sqlDB.Ping()
				So(err, ShouldBeNil)
			})

			Convey("应该能够创建测试表", func() {
				err := db.Exec(`
					CREATE TABLE IF NOT EXISTS integration_test_users (
						id SERIAL PRIMARY KEY,
						name VARCHAR(100) NOT NULL,
						email VARCHAR(100) UNIQUE NOT NULL,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
					)
				`).Error
				So(err, ShouldBeNil)
			})
		})

		Convey("当测试数据库操作时", func() {
			Convey("应该能够创建记录", func() {
				user := TestUser{
					Name:  "Integration Test User",
					Email: "integration@example.com",
				}

				result := db.Create(&user)
				So(result.Error, ShouldBeNil)
				So(result.RowsAffected, ShouldEqual, 1)
				So(user.ID, ShouldBeGreaterThan, 0)
			})

			Convey("应该能够查询记录", func() {
				// 先插入测试数据
				user := TestUser{
					Name:  "Query Test User",
					Email: "query@example.com",
				}
				db.Create(&user)

				// 查询记录
				var result TestUser
				err := db.First(&result, "email = ?", "query@example.com").Error
				So(err, ShouldBeNil)
				So(result.Name, ShouldEqual, "Query Test User")
				So(result.Email, ShouldEqual, "query@example.com")
			})

			Convey("应该能够更新记录", func() {
				// 先插入测试数据
				user := TestUser{
					Name:  "Update Test User",
					Email: "update@example.com",
				}
				db.Create(&user)

				// 更新记录
				result := db.Model(&user).Update("name", "Updated Integration User")
				So(result.Error, ShouldBeNil)
				So(result.RowsAffected, ShouldEqual, 1)

				// 验证更新
				var updatedUser TestUser
				err := db.First(&updatedUser, user.ID).Error
				So(err, ShouldBeNil)
				So(updatedUser.Name, ShouldEqual, "Updated Integration User")
			})

			Convey("应该能够删除记录", func() {
				// 先插入测试数据
				user := TestUser{
					Name:  "Delete Test User",
					Email: "delete@example.com",
				}
				db.Create(&user)

				// 删除记录
				result := db.Delete(&user)
				So(result.Error, ShouldBeNil)
				So(result.RowsAffected, ShouldEqual, 1)

				// 验证删除
				var deletedUser TestUser
				err := db.First(&deletedUser, user.ID).Error
				So(err, ShouldEqual, gorm.ErrRecordNotFound)
			})
		})

		Convey("当测试事务时", func() {
			Convey("应该能够执行事务操作", func() {
				// 开始事务
				tx := db.Begin()
				So(tx, ShouldNotBeNil)

				// 在事务中插入数据
				user := TestUser{
					Name:  "Transaction Test User",
					Email: "transaction@example.com",
				}
				result := tx.Create(&user)
				So(result.Error, ShouldBeNil)
				So(result.RowsAffected, ShouldEqual, 1)

				// 查询事务中的数据
				var count int64
				tx.Model(&TestUser{}).Count(&count)
				So(count, ShouldEqual, 1)

				// 提交事务
				err := tx.Commit().Error
				So(err, ShouldBeNil)

				// 验证数据已提交
				db.Model(&TestUser{}).Count(&count)
				So(count, ShouldBeGreaterThan, 0)
			})

			Convey("应该能够回滚事务", func() {
				// 开始事务
				tx := db.Begin()

				// 在事务中插入数据
				user := TestUser{
					Name:  "Rollback Test User",
					Email: "rollback@example.com",
				}
				tx.Create(&user)

				// 回滚事务
				err := tx.Rollback().Error
				So(err, ShouldBeNil)

				// 验证数据已回滚
				var count int64
				db.Model(&TestUser{}).Where("email = ?", "rollback@example.com").Count(&count)
				So(count, ShouldEqual, 0)
			})
		})

		Convey("当测试批量操作时", func() {
			Convey("应该能够批量插入记录", func() {
				users := []TestUser{
					{Name: "Batch User 1", Email: "batch1@example.com"},
					{Name: "Batch User 2", Email: "batch2@example.com"},
					{Name: "Batch User 3", Email: "batch3@example.com"},
				}

				result := db.Create(&users)
				So(result.Error, ShouldBeNil)
				So(result.RowsAffected, ShouldEqual, 3)

				// 验证批量插入
				var count int64
				db.Model(&TestUser{}).Where("email LIKE ?", "batch%@example.com").Count(&count)
				So(count, ShouldEqual, 3)
			})

			Convey("应该能够批量更新记录", func() {
				// 先插入测试数据
				users := []TestUser{
					{Name: "Batch Update 1", Email: "batchupdate1@example.com"},
					{Name: "Batch Update 2", Email: "batchupdate2@example.com"},
				}
				db.Create(&users)

				// 批量更新
				result := db.Model(&TestUser{}).
					Where("email LIKE ?", "batchupdate%@example.com").
					Update("name", "Batch Updated User")
				So(result.Error, ShouldBeNil)
				So(result.RowsAffected, ShouldEqual, 2)

				// 验证更新
				var updatedCount int64
				db.Model(&TestUser{}).
					Where("name = ?", "Batch Updated User").
					Count(&updatedCount)
				So(updatedCount, ShouldEqual, 2)
			})
		})

		Convey("当测试查询条件时", func() {
			Convey("应该能够使用各种查询条件", func() {
				// 插入测试数据
				testUsers := []TestUser{
					{Name: "Alice", Email: "alice@example.com"},
					{Name: "Bob", Email: "bob@example.com"},
					{Name: "Charlie", Email: "charlie@example.com"},
					{Name: "Alice Smith", Email: "alice.smith@example.com"},
				}
				db.Create(&testUsers)

				Convey("应该能够使用 LIKE 查询", func() {
					var users []TestUser
					err := db.Where("name LIKE ?", "Alice%").Find(&users).Error
					So(err, ShouldBeNil)
					So(len(users), ShouldEqual, 2)
				})

				Convey("应该能够使用 IN 查询", func() {
					var users []TestUser
					err := db.Where("name IN ?", []string{"Alice", "Bob"}).Find(&users).Error
					So(err, ShouldBeNil)
					So(len(users), ShouldEqual, 2)
				})

				Convey("应该能够使用 BETWEEN 查询", func() {
					var users []TestUser
					err := db.Where("id BETWEEN ? AND ?", 1, 3).Find(&users).Error
					So(err, ShouldBeNil)
					So(len(users), ShouldBeGreaterThan, 0)
				})

				Convey("应该能够使用多条件查询", func() {
					var users []TestUser
					err := db.Where("name LIKE ? AND email LIKE ?", "%Alice%", "%example.com").Find(&users).Error
					So(err, ShouldBeNil)
					So(len(users), ShouldEqual, 2)
				})
			})
		})

		Reset(func() {
			// 清理测试表
			if db != nil {
				db.Exec("DROP TABLE IF EXISTS integration_test_users")
			}
			// 关闭数据库连接
			if sqlDB != nil {
				sqlDB.Close()
			}
		})
	})
}

// TestDatabaseConnectionPool 测试数据库连接池
func TestDatabaseConnectionPool(t *testing.T) {
	Convey("数据库连接池测试", t, func() {
		var db *gorm.DB
		var sqlDB *sql.DB

		Convey("当配置连接池时", func() {
			testConfig := &config.Config{
				Database: config.DatabaseConfig{
					Host:            "localhost",
					Port:            5432,
					Database:        "{{.ProjectName}}_test_pool",
					Username:        "postgres",
					Password:        "password",
					SslMode:         "disable",
					MaxIdleConns:    5,
					MaxOpenConns:    10,
					ConnMaxLifetime: 5 * time.Minute,
				},
			}

			dsn := testConfig.Database.GetDSN()
			var err error
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			So(err, ShouldBeNil)

			sqlDB, err = db.DB()
			So(err, ShouldBeNil)

			Convey("应该能够设置连接池参数", func() {
				sqlDB.SetMaxIdleConns(testConfig.Database.MaxIdleConns)
				sqlDB.SetMaxOpenConns(testConfig.Database.MaxOpenConns)
				sqlDB.SetConnMaxLifetime(testConfig.Database.ConnMaxLifetime)

				// 验证设置
				stats := sqlDB.Stats()
				So(stats.MaxOpenConns, ShouldEqual, testConfig.Database.MaxOpenConns)
				So(stats.MaxIdleConns, ShouldEqual, testConfig.Database.MaxIdleConns)
			})

			Convey("应该能够监控连接池状态", func() {
				// 获取初始状态
				initialStats := sqlDB.Stats()
				So(initialStats.OpenConnections, ShouldEqual, 0)

				// 执行一些查询来创建连接
				for i := 0; i < 3; i++ {
					sqlDB.Ping()
				}

				// 获取使用后的状态
				afterStats := sqlDB.Stats()
				So(afterStats.OpenConnections, ShouldBeGreaterThan, 0)
				So(afterStats.InUse, ShouldBeGreaterThan, 0)
			})
		})

		Reset(func() {
			// 关闭数据库连接
			if sqlDB != nil {
				sqlDB.Close()
			}
		})
	})
}