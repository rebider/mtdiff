package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type mysqlHandle struct {
	dsn string
	conn *sql.DB
}

func NewMysqlHandle(dsn string) (m *mysqlHandle, err error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	err = db.Ping()
	if err != nil {
		return
	}
	m = &mysqlHandle{
		dsn: dsn,
		conn: db,
	}
	return
}

// mysql> show tables;
func (m *mysqlHandle) ShowTables() (tables []string, err error) {

	rows, err := m.conn.Query("SHOW TABLES;")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return
		}
		tables = append(tables, name)
	}
	return
}

// mysql> desc table
func (m *mysqlHandle) DescTable(table string) (desc []map[string]string, err error) {

	rows, err := m.conn.Query("DESC "+table+";")
	if err != nil {
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return
	}
	vals := make([][]byte, len(columns))
	scan := make([]interface{}, len(columns))
	for i := range scan {
		scan[i] = &vals[i]
	}
	for rows.Next() {
		if rows.Scan(scan...) != nil {
			return
		}
		row := map[string]string{}
		for i, v := range vals {
			key := columns[i]
			row[key] = string(v)
		}
		desc = append(desc, row)
	}
	return
}

// check table is exists
func (m *mysqlHandle) TableIsExists(table string) (bool, error) {
	tables, err := m.ShowTables()
	if err != nil {
		return false, err
	}
	for _, t := range tables {
		if t == table {
			return true, nil
		}
	}
	return false, nil
}

// mysql>show create table table_name
func (m *mysqlHandle) ShowCreateTable(table string) (createTable string, err error) {
	row := m.conn.QueryRow("SHOW CREATE TABLE "+table+";")
	var name string
	err = row.Scan(&name, &createTable)
	if err != nil {
		return
	}
	return createTable, nil
}

// drop table
func (m *mysqlHandle) DropTable(table string) (err error) {
	_, err = m.conn.Exec("DROP TABLE IF EXISTS "+table)
	return
}

// create table
func (m *mysqlHandle) CreateTable(table string, sql string) (err error) {
	_, err = m.conn.Exec(sql+";")
	if err != nil {
		return
	}
	_, err = m.conn.Exec("alter table "+table+" auto_increment = 0;")
	return
}

// close mysql conn
func (m *mysqlHandle) Close() {
	m.conn.Close()
}