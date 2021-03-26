package main

import (
	"GumboSql"
	"database/sql"
	"log"
)
func main(){
	var dbCreator GumboSql.DatabaseCreator
	connection := "root:@tcp(127.0.0.1:3306)/test_db"
	db, err := sql.Open("mysql", connection)
	if err != nil{
		log.Print(err.Error())
	}
	dbCreator = dbCreator.DatabaseCreator(db)
	dbCreator = dbCreator.Table("new_table").Column("id","int (255) NOT NULL AUTO_INCREMENT").Column("test_column","VARCHAR (255)").Column("other_one","INT (255)").Column("","PRIMARY KEY (id)")
	dbCreator = dbCreator.Table("table_two").Column("id","int (255) NOT NULL AUTO_INCREMENT").Column("name","VARCHAR (255)").Column("user_code","INT (255)").Column("","PRIMARY KEY (id)")
	dbCreator.Init()

	// query builder
	var queryBuilder GumboSql.QueryBuilder
	queryBuilder = queryBuilder.QueryBuilder(db)
	queryBuilder = queryBuilder.Table("table_two")
	queryBuilder.Insert([]string{"name", "user_code"},[]string{"jafar", "5"})
	queryBuilder.Insert([]string{"name", "user_code"},[]string{"mamad", "2"})
	queryBuilder.Insert([]string{"name", "user_code"},[]string{"reza", "3"})
	queryBuilder.Where("name","jafar").First() // return *sql.Row
	queryBuilder.Limit(5,0).OrderBy("name", "DESC").Get() // return *sql.Rows
	queryBuilder.Where("name","jafar").Update([]string{"name"},[]string{"moji"}) // update
	queryBuilder.Where("name","jafar").OrWhere("name","moji").First()
}
