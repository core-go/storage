package upload

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
)

type StorageRepository interface {
	Load(ctx context.Context, id string, column string) (UploadModel, error)
	// Update(ctx context.Context, id string) (UploadModel, error)
	Update(ctx context.Context, item UploadModel) (int64, error)
	// Upload(ctx context.Context, directory string, filename string, data []byte, contentType string) (string, error)
	// Delete(ctx context.Context, id string) (bool, error)
}

func CreateStorageRepository(DB *sql.DB,
	Table string,
	columns UploadFieldColumn, toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}) StorageRepository {
	return &storageRepository{DB: DB, Table: Table, Columns: &columns, toArray: toArray}
}

type UploadFieldColumn struct {
	Cover   *string
	Image   *string
	Gallery *string
	Id      string
}

type storageRepository struct {
	DB      *sql.DB
	Table   string
	Columns *UploadFieldColumn
	toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
	// CoverUrl string
}

func (s *storageRepository) Load(ctx context.Context, id string, column string) (UploadModel, error) {
	var data UploadModel
	query := BuildFindById(s.Table, s.Columns.Id, *s.Columns, column)
	rows, err := s.DB.QueryContext(ctx, query, id)
	if err != nil {
		return data, err
	}
	defer rows.Close()
	rows.Next()
	// if !rows.Next() {
	// 	return data, errors.New("data not found")
	// }
	var err2 error
	if *s.Columns.Cover == column {
		err2 = rows.Scan(&data.UserId, &data.CoverURL)
	} else if *s.Columns.Image == column {
		err2 = rows.Scan(&data.UserId, &data.ImageURL)
	} else if *s.Columns.Gallery == column {
		err2 = rows.Scan(&data.UserId, s.toArray(&data.Gallery))
	}
	if err2 != nil {
		return data, err2
	}
	return data, nil
}

func (s *storageRepository) Update(ctx context.Context, item UploadModel) (int64, error) {
	// query := fmt.Sprintf("update %s set %s = $1 where %s =$2", s.Table, *s.Columns.Cover, s.IdCol)
	query, value := BuildUpdate(s.Table, *s.Columns, item, s.toArray)
	stmt, er0 := s.DB.Prepare(query)
	if er0 != nil {
		return -1, nil
	}
	res, err := stmt.ExecContext(ctx, value...)

	row, er2 := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	if row < 0 {
		return -1, er2
	}
	return row, er2
}

func BuildFindById(table string, id string, columns UploadFieldColumn, column string) string {
	var where = ""
	var selectQuery = []string{}
	where = fmt.Sprintf("where %s = $1", id)
	selectQuery = append(selectQuery, id)
	if *columns.Image != "" && *columns.Image == column {
		selectQuery = append(selectQuery, *columns.Image)
	}
	if *columns.Cover != "" && *columns.Cover == column {
		selectQuery = append(selectQuery, *columns.Cover)
	}
	if *columns.Gallery != "" && *columns.Gallery == column {
		selectQuery = append(selectQuery, *columns.Gallery)
	}
	selectFields := strings.Join(selectQuery, ",")
	return fmt.Sprintf("select %s from %v %v", selectFields, table, where)
}

func BuildUpdate(table string, columns UploadFieldColumn, item UploadModel, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) (string, []interface{}) {
	var setQuery = []string{}
	var value []interface{}
	var index = 1
	if *columns.Image != "" && item.ImageURL != nil {
		setData := fmt.Sprintf("%s = $%d", *columns.Image, index)
		value = append(value, item.ImageURL)
		setQuery = append(setQuery, setData)
		index++
	}
	if *columns.Cover != "" && item.CoverURL != nil {
		setData := fmt.Sprintf("%s = $%d", *columns.Cover, index)
		value = append(value, item.CoverURL)
		setQuery = append(setQuery, setData)
		index++
	}
	if *columns.Gallery != "" && item.Gallery != nil {
		var interfaceSlice []interface{} = make([]interface{}, len(item.Gallery))
		for i, d := range item.Gallery {
			interfaceSlice[i] = d
		}
		setData := fmt.Sprintf("%s = $%d", *columns.Gallery, index)

		value = append(value, toArray(item.Gallery))
		setQuery = append(setQuery, setData)
		index++
	}
	where := fmt.Sprintf("where %s = $%d", columns.Id, index)
	value = append(value, item.UserId)
	sets := "set " + strings.Join(setQuery, ",")
	query := fmt.Sprintf("update %s %s %s", table, sets, where)
	return query, value
}
