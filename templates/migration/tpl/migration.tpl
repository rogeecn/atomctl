package migrations

import (
	"atom/container"
	"atom/contracts"
	"atom/providers/log"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func init() {
	if err := container.Container.Provide(New{{.ID}}{{.PascalMigrationName}}Migration, dig.Group("migrations")); err != nil {
		log.Fatal(err)
	}
}

type Migration{{.ID}}{{.PascalMigrationName}} struct {
	id string
}

func New{{.ID}}{{.PascalMigrationName}}Migration() contracts.Migration {
	return &Migration{{.ID}}{{.PascalMigrationName}}{id: "{{.ID}}_{{.SnakeMigrationName}}"}
}

func (m *Migration{{.ID}}{{.PascalMigrationName}}) ID() string {
	return m.id
}

func (m *Migration{{.ID}}{{.PascalMigrationName}}) Up(tx *gorm.DB) error {
    table := m.table()
	return tx.AutoMigrate(&table)
}

func (m *Migration{{.ID}}{{.PascalMigrationName}}) Down(tx *gorm.DB) error {
    return tx.Migrator().DropTable(m.table())
}

func (m *Migration{{.ID}}{{.PascalMigrationName}}) table() interface{} {
	type TableName struct {
		FieldName string
	}

	return TableName{}
}