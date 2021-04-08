package GumboSql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DatabaseCreator struct {
	table_str string
	query_str string
	queries []string
	connection *sql.DB
}
func (dc DatabaseCreator) DatabaseCreator(db *sql.DB) DatabaseCreator {
	dc.connection = db
	return dc
}
func (dc DatabaseCreator) Table(table string) DatabaseCreator {
	if dc.query_str != "" {
		if last := len(dc.query_str) - 1; last >= 0 && dc.query_str[last] == ',' {
			dc.query_str = dc.query_str[:last]
		}
		dc.queries = append(dc.queries,dc.query_str+")")
		dc.query_str = ""
	}
	dc.query_str += "CREATE TABLE IF NOT EXISTS "+table+" ("
	return dc
}

func (dc DatabaseCreator) Column(name string, dataType string) DatabaseCreator {
dc.query_str += name + " " + dataType +","
	return dc
}

func (dc DatabaseCreator) ID() DatabaseCreator {
	dc.query_str += "id INT (255) NOT NULL AUTO_INCREMENT,PRIMARY KEY (id),"
	return dc
}
func (dc DatabaseCreator) BigID() DatabaseCreator {
	dc.query_str += "id BIGINT (20) NOT NULL AUTO_INCREMENT,PRIMARY KEY (id),"
	return dc
}

func (dc DatabaseCreator) Init(){
	if last := len(dc.query_str) - 1; last >= 0 && dc.query_str[last] == ',' {
		dc.query_str = dc.query_str[:last]
	}
	dc.queries = append(dc.queries,dc.query_str+")")
	for i := range dc.queries {
		_, err:= dc.connection.Query(dc.queries[i])
if err != nil{
	log.Print(err.Error())
}
	}
}