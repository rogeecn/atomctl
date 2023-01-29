package seeders

import (
	"atom/container"
	"atom/contracts"
	"log"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

func init() {
	if err := container.Container.Provide(New{{.PascalSeederName}}Seeder, dig.Group("seeders")); err != nil {
		log.Fatal(err)
	}
}

type {{.PascalSeederName}}Seeder struct {
}

func New{{.PascalSeederName}}Seeder() contracts.Seeder {
	return &{{.PascalSeederName}}Seeder{}
}

func (s *{{.PascalSeederName}}Seeder) Run(db *gorm.DB) {
	times := 10
	for i := 0; i < times; i++ {
		data := s.Generate(i)
		if i == 0 {
			stmt := &gorm.Statement{DB: db}
			_ = stmt.Parse(&data)
			log.Printf("seeding %s for %d times", stmt.Schema.Table, times)
		}
		db.Create(&data)
	}
}

func (s *{{.PascalSeederName}}Seeder) Generate(idx int) models.{{.PascalSeederName}} {
	return models.{{.PascalSeederName}}{
        // fill model fields
	}
}
