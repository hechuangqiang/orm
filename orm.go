package orm

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

var dbHive map[string]*sql.DB = make(map[string]*sql.DB)

//map record
type Record map[int]map[string]string

//create new database
func NewDatabase(dbname, dbtype, url string) {
	db, err := sql.Open(dbtype, url)
	if err != nil {
		log.Println("open database fail", err)
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Println("database connect fail", err)
		panic(err.Error())
	}
	dbHive[dbname] = db
}

//module
type Module struct {
	columnstr string //select field
	tableName string //table
	filters   string //condition
	orderby   string //orderby
	limit     string //limit
	join      string //join
	dbname    string //dbname
}

//create new Module
func NewModule(tableName string) *Module {
	m := &Module{tableName: tableName, columnstr: "*", dbname: "default"}
	return m
}

//change db
func (m *Module) User(dbname string) *Module {
	m.dbname = dbname
	return m
}

//select fields
func (m *Module) Select(fields ...string) *Module {
	m.columnstr = ""
	for _, f := range fields {
		m.columnstr = m.columnstr + f + ","
	}
	return m
}

//Filter
func (m *Module) Filter(param ...interface{}) *Module {
	for _, p := range param {
		m.filters += fmt.Sprintf("%v", p)
	}
	return m
}

//orderBy
func (m *Module) OrderBy(param string) *Module {
	m.orderby = fmt.Sprintf("ORDER By %v", param)
	return m
}

//limit
func (m *Module) Limit(size ...int) *Module {
	if len(size) > 1 {
		m.limit = fmt.Sprintf("Limit %d,%d", size[0], size[1])
		return m
	} else {
		m.limit = fmt.Sprintf("Limit %d", size[0])
		return m
	}
}

//leftJoin
func (m *Module) LeftJoin(table, condition string) *Module {
	m.join = fmt.Sprintf("LEFT JOIN %v ON %v", table, condition)
	return m
}

//rightJoin
func (m *Module) RightJoin(table, condition string) *Module {
	m.join = fmt.Sprintf("RIGHT JOIN %v ON %v", table, condition)
	return m
}

//join
func (m *Module) Join(table, condition string) *Module {
	m.join = fmt.Sprintf("INNER JOIN %v ON %v", table, condition)
	return m
}

//fulljoin
func (m *Module) FullJoin(table, condition string) *Module {
	m.join = fmt.Sprintf("FULL JOIN %v ON %v", table, condition)
	return m
}

func (m *Module) FindAll(records interface{}) {
	db := dbHive[m.dbname]
	columnstr := m.columnstr
	if l := len(columnstr); l > 1 {
		columnstr = columnstr[:l-1]
	}
	where := m.filters
	where = strings.TrimSpace(where)
	if len(where) > 0 {
		where = "where " + where
	}
	query := fmt.Sprintf("select %v from %v %v %v %v %v", columnstr, m.tableName, m.join, where, m.orderby, m.limit)
	log.Println("query = ", query)
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
}
