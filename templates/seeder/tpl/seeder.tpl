package seeders

import (
	"log"

	"{{.Package}}/database/models"

	"gorm.io/gorm"
	"github.com/rogeecn/atom/contracts"
	"github.com/brianvoe/gofakeit/v6"
)

type {{.PascalSeederName}}Seeder struct {
}

func New{{.PascalSeederName}}Seeder() contracts.Seeder {
	return &{{.PascalSeederName}}Seeder{}
}

func (s *{{.PascalSeederName}}Seeder) Run(faker *gofakeit.Faker, db *gorm.DB) {
	times := 10
	for i := 0; i < times; i++ {
		data := s.Generate(faker, i)
		if i == 0 {
			stmt := &gorm.Statement{DB: db}
			_ = stmt.Parse(&data)
			log.Printf("seeding %s for %d times", stmt.Schema.Table, times)
		}
		db.Create(&data)
	}
}

func (s *{{.PascalSeederName}}Seeder) Generate(faker *gofakeit.Faker,idx int) models.{{.PascalSeederName}} {
	return models.{{.PascalSeederName}}{
        // fill model fields
	}
}
