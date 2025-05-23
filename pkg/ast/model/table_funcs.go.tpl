package model

import (
	"context"
	"time"

	"{{ .PkgName }}/database/table"

	. "github.com/go-jet/jet/v2/postgres"
	log "github.com/sirupsen/logrus"
)


func (m *{{.PascalTable}}) log() *log.Entry {
	return log.WithField("model", "{{.PascalTable}}")
}

func (m *{{.PascalTable}}) Create(ctx context.Context) error {
	m.CreatedAt = time.Now()
	stmt := table.Medias.INSERT(table.{{.PascalTable}}.MutableColumns).MODEL(m).RETURNING(table.Medias.AllColumns)
	m.log().WithField("func","Create").Info( stmt.DebugSql())

	if err := stmt.QueryContext(ctx, db, m); err != nil {
		m.log().WithField("func","Create").Errorf("error creating {{.PascalTable}} item: %v", err)
		return err
	}

	m.log().WithField("func","Create").Infof("{{.PascalTable}} item created successfully")
	return nil
}


func (m *{{.PascalTable}}) BatchCreate(ctx context.Context, models []*{{.PascalTable}}) error {
	stmt := table.{{.PascalTable}}.INSERT(table.{{.PascalTable}}.MutableColumns).MODELS(models)
	m.log().WithField("func", "BatchCreate").Info(stmt.DebugSql())

	if _, err := stmt.ExecContext(ctx, db); err != nil {
		m.log().WithField("func","Create").Errorf("error creating {{.PascalTable}} item: %v", err)
		return err
	}

	m.log().WithField("func", "BatchCreate").Infof("{{.PascalTable}} items created successfully")
	return nil
}

// Delete
func (m *{{.PascalTable}}) Delete(ctx context.Context) error {
	stmt := table.{{.PascalTable}}.DELETE().WHERE(table.{{.PascalTable}}.ID.EQ(m.ID))
	m.log().WithField("func", "Delete").Info(stmt.DebugSql())

	if _, err := stmt.ExecContext(ctx, db); err != nil {
		m.log().WithField("func","Delete").Errorf("error deleting {{.PascalTable}} item: %v", err)
		return err
	}

	m.log().WithField("func", "Delete").Infof("{{.PascalTable}} item deleted successfully")
	return nil
}

func (m *{{.PascalTable}}) BatchDelete(ctx context.Context, ids []int64) error {
	stmt := table.{{.PascalTable}}.DELETE().WHERE(table.{{.PascalTable}}.ID.IN(ids))
	m.log().WithField("func", "BatchDelete").Info(stmt.DebugSql())

	if _, err := stmt.ExecContext(ctx, db); err != nil {
		m.log().WithField("func","BatchDelete").Errorf("error deleting {{.PascalTable}} items: %v", err)
		return err
	}

	m.log().WithField("func", "BatchDelete").Infof("{{.PascalTable}} items deleted successfully")
	return nil
}

func (m *{{.PascalTable}}) Update(ctx context.Context) error {
	stmt := table.{{.PascalTable}}.UPDATE(table.{{.PascalTable}}.MutableColumns).SET(m).WHERE(table.{{.PascalTable}}.ID.EQ(m.ID)).RETURNING(table.{{.PascalTable}}.AllColumns)
	m.log().WithField("func", "Update").Info(stmt.DebugSql())

	if err := stmt.QueryContext(ctx, db, m); err != nil {
		m.log().WithField("func","Update").Errorf("error updating {{.PascalTable}} item: %v", err)
		return err
	}

	m.log().WithField("func", "Update").Infof("{{.PascalTable}} item updated successfully")
	return nil
}
// GetByID
func (m *{{.PascalTable}}) GetByID(ctx context.Context, id int64) (*{{.PascalTable}}, error) {
	stmt := table.{{.PascalTable}}.SELECT(table.{{.PascalTable}}.AllColumns).WHERE(table.{{.PascalTable}}.ID.EQ(id))
	m.log().WithField("func", "GetByID").Info(stmt.DebugSql())

	if err := stmt.QueryContext(ctx, db, m); err != nil {
		m.log().WithField("func","GetByID").Errorf("error getting {{.PascalTable}} item by ID: %v", err)
		return nil, err
	}

	m.log().WithField("func", "GetByID").Infof("{{.PascalTable}} item retrieved successfully")
	return m, nil
}


// Count
func (m *{{.PascalTable}}) Count(ctx context.Context, conds ...BoolExpression) (int64, error) {
	cond := Bool(true)
	if len(conds) > 0 {
		for _, c := range conds {
			cond = cond.AND(c)
		}
	}

	tbl := table.{{.PascalTable}}
	stmt := tbl.SELECT(COUNT(tbl.ID).AS("count")).WHERE(cond)
	m.log().Infof("sql: %s", stmt.DebugSql())

	var count struct {
		Count int64
	}

	if err := stmt.QueryContext(ctx, db, &count); err != nil {
		m.log().Errorf("error counting {{.PascalTable}} items: %v", err)
		return 0, err
	}

	return count.Count, nil
}

