package GumboSql

import "database/sql"

type PaginateModel struct {
	TotalPages int
	CurrentPage int
	ResultsPerPage int
	rows *sql.Rows
}
