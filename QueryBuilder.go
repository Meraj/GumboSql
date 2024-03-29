package GumboSql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"strconv"
)

type QueryBuilder struct {
	val Variables
}
type Variables struct {
	db             *sql.DB
	table          string
	columns        []string
	values         []string
	setColumns     []string
	whereStatement string
	orderBy        string
	limitOffset    string
	args           []interface{}
}

func (b QueryBuilder) QueryBuilder(db *sql.DB) QueryBuilder {
	b.val.db = db
	return b
}

func (b QueryBuilder) Table(table string) QueryBuilder {
	b.val.table = table
	return b
}

func (b QueryBuilder) SelectColumn(column string) QueryBuilder {
	b.val.columns = nil
	b.val.columns = append(b.val.columns, column)
	return b
}

func (b QueryBuilder) AddSelect(column string) QueryBuilder {
	b.val.columns = append(b.val.columns, column)
	return b
}

func (b QueryBuilder) SelectColumns(column []string) QueryBuilder {
	b.val.columns = column
	return b
}

func (b QueryBuilder) Where(column string, value interface{}) QueryBuilder {
	if b.val.whereStatement == "" {
		b.val.whereStatement = " WHERE " + column + " = ? "
	} else {
		b.val.whereStatement += " AND " + column + " = ? "
	}
	b.val.args = append(b.val.args, value)
	return b
}

func (b QueryBuilder) WhereIsNull(column string) QueryBuilder {
	if b.val.whereStatement == "" {
		b.val.whereStatement = " WHERE " + column + " IS NULL "
	} else {
		b.val.whereStatement += " AND " + column + " IS NULL "
	}
	return b
}

func (b QueryBuilder) OrWhereIsNull(column string) QueryBuilder {
	b.val.whereStatement += " OR " + column + " IS NULL "
	return b
}

func (b QueryBuilder) WhereIsNotNull(column string) QueryBuilder {
	if b.val.whereStatement == "" {
		b.val.whereStatement = " WHERE " + column + " IS NOT NULL "
	} else {
		b.val.whereStatement += " AND " + column + " IS NULL "
	}
	return b
}

func (b QueryBuilder) OrWhereIsNotNull(column string) QueryBuilder {
	b.val.whereStatement += " OR " + column + " IS NOT NULL "
	return b
}


func (b QueryBuilder) OrWhere(column string, value interface{}) QueryBuilder {
	b.val.whereStatement += " OR " + column + " = ? "
	b.val.args = append(b.val.args, value)
	return b
}

func (b QueryBuilder) WhereWithOperation(column string, operation string, value interface{}) QueryBuilder {
	if b.val.whereStatement == "" {
		b.val.whereStatement = " WHERE " + column + " " + operation + " ? "
	} else {
		b.val.whereStatement += " AND " + column + " " + operation + " ? "
	}
	b.val.args = append(b.val.args, value)
	return b
}

func (b QueryBuilder) OrWhereWithOperation(column string, operation string, value interface{}) QueryBuilder {
	b.val.whereStatement += " OR " + column + " " + operation + " ? "
	b.val.args = append(b.val.args, value)
	return b
}

func (b QueryBuilder) OrderBy(column string, orderType string) QueryBuilder {
	b.val.orderBy = "ORDER BY " + column + " " + orderType
	return b
}

func (b QueryBuilder) Limit(limit_int int, offset_int int) QueryBuilder {
	b.val.limitOffset = " LIMIT " + strconv.Itoa(limit_int) + " OFFSET " + strconv.Itoa(offset_int)
	return b
}

func (b QueryBuilder) buildQuery(SqlType int) string {
	query := ""
	switch SqlType {
	case 0:
		query = "INSERT INTO " + b.val.table + " ("
		for i := range b.val.columns {
			query += b.val.columns[i] + ","
		}
		if last := len(query) - 1; last >= 0 && query[last] == ',' {
			query = query[:last]
		}
		query += ") VALUES ("
		for _ = range b.val.columns {
			query += "? ,"
		}
		if last := len(query) - 1; last >= 0 && query[last] == ',' {
			query = query[:last]
		}
		query += ")"
		return query
		break
	case 1:
		query = "SELECT "
		if b.val.columns == nil {
			query += " * "
		} else {
			for i := range b.val.columns {
				query += b.val.columns[i] + ","
			}
			if last := len(query) - 1; last >= 0 && query[last] == ',' {
				query = query[:last]
			}
		}
		query += " FROM " + b.val.table + " "
		break
	case 2:
		query = "UPDATE " + b.val.table + " SET "
		for i := range b.val.setColumns {
			query += b.val.setColumns[i] + " = ?,"
		}
		if last := len(query) - 1; last >= 0 && query[last] == ',' {
			query = query[:last]
		}
		break

	case 3:
		query = "DELETE FROM " + b.val.table + " "
	}
	if b.val.whereStatement != "" {
		query += " " + b.val.whereStatement
	}
	if b.val.orderBy != "" {
		query += " " + b.val.orderBy
	}
	if b.val.limitOffset != "" {
		query += " " + b.val.limitOffset
	}
	return query
}
func (qb QueryBuilder) Insert(theData map[string]interface{}) (interface{}, error) {
	qb.vars.args = nil
	qb.vars.columns = nil
	for key, element := range theData {
		qb.vars.columns = append(qb.vars.columns, key)
		qb.vars.args = append(qb.vars.args, element)
	}
	var id interface{}
	err := qb.DB.QueryRow(qb.buildQuery(0), qb.vars.args...).Scan(&id)
	qb.vars.args = nil
	if err != nil {
		id = 0
	}
	return id, err
}

func (b QueryBuilder) First() *sql.Row {
	row := b.val.db.QueryRow(b.buildQuery(1), b.val.args...)
	return row
}

func (b QueryBuilder) Get() *sql.Rows {
	row, err := b.val.db.Query(b.buildQuery(1), b.val.args...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return row
}

func (b QueryBuilder) Update(columns []string, values ...interface{}) sql.Result {
	b.val.setColumns = columns
	queryValues := b.val.args
	b.val.args = nil
	for i := range values {
		b.val.args = append(b.val.args, values[i])
	}
	for i := range queryValues {
		b.val.args = append(b.val.args, queryValues[i])
	}
	query := b.buildQuery(2)
	res, err := b.val.db.Exec(query, b.val.args...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return res
}

func (b QueryBuilder) Delete() int64 {
	query := b.buildQuery(3)
	res, err := b.val.db.Exec(query, b.val.args...)
	if err != nil{
		print(err.Error())
	}
	count, err := res.RowsAffected()
	if err != nil {
		print(err.Error())
	}
	return count

}
func (b QueryBuilder) Count() int {
	b.val.columns = nil
	b.val.columns = append(b.val.columns, "COUNT(*) AS total")
	var total int
	err := b.First().Scan(&total)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return total
}

func (b QueryBuilder) Paginate(itemsPerPage int, currentPage int) PaginateModel {
	totalPages := int(math.Ceil(float64(b.Count() / itemsPerPage)))
	limitInt := (currentPage - 1) * itemsPerPage
	b.Limit(itemsPerPage, limitInt)
	var paginateModel PaginateModel
	paginateModel.TotalPages = totalPages + 1
	paginateModel.CurrentPage = currentPage
	paginateModel.ResultsPerPage = itemsPerPage
	paginateModel.Rows = b.Get()
	return paginateModel
}

func (b QueryBuilder) Connection()  *sql.DB {
	return b.val.db
}
