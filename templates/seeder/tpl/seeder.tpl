package seeders

import (
	"atom/container"
	"atom/contracts"
	"atom/database/models"
	"log"

	"go.uber.org/dig"
	"gorm.io/gorm"
	"github.com/brianvoe/gofakeit/v6"
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
