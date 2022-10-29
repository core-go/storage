package uploads

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	q "github.com/core-go/sql"
	"reflect"
)

type StorageRepository interface {
	Load(ctx context.Context, id string) (*Entity, error)
	Update(ctx context.Context, id string, attachments []Attachment) (int64, error)
}

func NewRepository(DB *sql.DB,
	Table string,
	columns FieldColumn, toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}) *SqlRepository {
	utype := reflect.TypeOf(Entity{})
	fieldIndex, _ := q.GetColumnIndexes(utype)
	return &SqlRepository{DB: DB, Table: Table, Columns: &columns, toArray: toArray, utype: utype, fieldIndex: fieldIndex}
}

type FieldColumn struct {
	Id   string
	File string
}

type SqlRepository struct {
	DB      *sql.DB
	Table   string
	Columns *FieldColumn
	toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
	utype      reflect.Type
	fieldIndex map[string]int
}

func (s *SqlRepository) Load(ctx context.Context, id string) (*Entity, error) {
	var model []Entity
	query := fmt.Sprintf("select %s, %s from %s where %s= $1", s.Columns.Id, s.Columns.File, s.Table, s.Columns.Id)
	err := q.QueryWithArray(ctx, s.DB, s.fieldIndex, &model, s.toArray, query, id)
	if err != nil {
		return nil, err
	}
	if len(model) > 0 {
		return &model[0], nil
	}
	return nil, err
}

func (s *SqlRepository) Update(ctx context.Context, id string, attachments []Attachment) (int64, error) {
	query := fmt.Sprintf("update %s set %s = $1 where %s =$2", s.Table, s.Columns.File, s.Columns.Id)
	stmt, er0 := s.DB.Prepare(query)
	if er0 != nil {
		return -1, er0
	}
	res, err := stmt.ExecContext(ctx, s.toArray(attachments), id)
	if err != nil {
		return -1, err
	}
	row, er2 := res.RowsAffected()

	if row < 0 {
		return -1, er2
	}
	return row, er2
}
