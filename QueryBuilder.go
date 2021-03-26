package GumboSql

import (
	"database/sql"
)

type QueryBuilder struct {

	val Variables
}
type Variables struct {
	db *sql.DB
	table          string
	columns        []string
	values         []string
	whereStatement string
	args []interface{}
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
	b.val.columns = append(b.val.columns,column)
	return b
}

func (b QueryBuilder) AddSelect(column string) QueryBuilder {
	b.val.columns = append(b.val.columns,column)
	return b
}

func (b QueryBuilder) SelectColumns(column []string) QueryBuilder {
	b.val.columns = column
	return b
}

func (b QueryBuilder) Where(column string,value string) QueryBuilder {
	if b.val.whereStatement == "" {
		b.val.whereStatement = " WHERE " + column + " = ? "
	}else{
		b.val.whereStatement += " AND " + column + " = ? "
	}
	b.val.args = append(b.val.args,value)
	return b
}

func (b QueryBuilder) OrWhere(column string,value string) QueryBuilder {
		b.val.whereStatement += " OR " + column + " = ? "
	b.val.args = append(b.val.args,value)
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
			query+="? ,"
		}
		if last := len(query) - 1; last >= 0 && query[last] == ',' {
			query = query[:last]
		}
		query +=")"
		return query
		break
	case 1:
		query = "SELECT "
		if b.val.columns == nil{
			query += " * "
		}else{
			for i := range b.val.columns {
				query += b.val.columns[i]+","
			}
			if last := len(query) - 1; last >= 0 && query[last] == ',' {
				query = query[:last]
			}
		}
		query += " FROM " + b.val.table + " "
		break
	}
	if b.val.whereStatement != "" {
		query += " " + b.val.whereStatement
	}
	return query
}
func (b QueryBuilder) Insert(columns []string, values []string) int64 {
	b.val.columns = columns
	for i := range values {
		b.val.args = append(b.val.args,values[i])
	}
	query := b.buildQuery(0)
	res, err := b.val.db.Exec(query, b.val.args...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}else{
		id, err := res.LastInsertId()
		if err == nil {

			return id
		}
	}
	return 0
}

func (b QueryBuilder) ToSql() string {
	return b.buildQuery(1)
}
