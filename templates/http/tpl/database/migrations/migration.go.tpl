package migrations

import (
	"time"

	"github.com/rogeecn/atom/contracts"
	"gorm.io/gorm"
)

var Migrations = []contracts.MigrationProvider{}


type Model struct {
	ID        uint           `gorm:"primarykey;comment:ID"`
	CreatedAt time.Time      `gorm:"comment:创建时间"`
	UpdatedAt time.Time      `gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

type ModelOnlyID struct {
	ID        uint      `gorm:"primarykey;comment:ID"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
}

type ModelNoSoftDelete struct {
	ID        uint      `gorm:"primarykey;comment:ID"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
}

type ModelWithUser struct {
	TenantID uint `gorm:"comment:租户ID"`
	UserID   uint `gorm:"comment:用户ID"`
}
