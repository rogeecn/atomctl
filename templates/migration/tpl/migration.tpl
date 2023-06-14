package migrations

import (
	"github.com/rogeecn/atom/contracts"
	"gorm.io/gorm"
)

func (m *Migration{{.ID}}{{.PascalMigrationName}}) table() interface{} {
	type TableName struct {
		FieldName string
	}

	return TableName{}
}

func (m *Migration{{.ID}}{{.PascalMigrationName}}) Up(tx *gorm.DB) error {
	return tx.AutoMigrate(m.table())
}

func (m *Migration{{.ID}}{{.PascalMigrationName}}) Down(tx *gorm.DB) error {
    return tx.Migrator().DropTable(m.table())
	// return tx.Migrator().DropColumn(m.table(), "input_column_name")
}

// DO NOT EDIT BLOW CODES!!
// DO NOT EDIT BLOW CODES!!
// DO NOT EDIT BLOW CODES!!
// DO NOT EDIT BLOW CODES!!
// DO NOT EDIT BLOW CODES!!
// DO NOT EDIT BLOW CODES!!
func init() {
	Migrations = append(Migrations, New{{.ID}}{{.PascalMigrationName}}Migration)
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
