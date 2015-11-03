package orm

import (
	"database/sql"
	"fmt"
	"strings"
)

var dbHive map[string]*sql.DB = make(map[string]*sql.DB)

//create new database
func NewDatabase(dbname, dbtype, url string) {
	db, err := sql.Open(dbtype, url)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
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

func (m *Module) getSqlString() string {
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
	return query
}

func (m *Module) FindAll(records interface{}) error {
	db := dbHive[m.dbname]
	rows, err := db.Query(m.getSqlString())
	if err != nil {
		return err
	}
	defer rows.Close()
	if value, ok := records.([]Record); ok {
		if value == nil {
			value = make([]Record, 0)
		}
	}
	return nil
}

func (m *Module) OneRecord() (Record, error) {
	rs, err := m.Limit(1).AllRecords()
	if err != nil {
		return nil, err
	}
	return rs[0], nil
}

func (m *Module) AllRecords() ([]Record, error) {
	db := dbHive[m.dbname]
	rows, err := db.Query(m.getSqlString())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	records := make([]Record, 0)
	columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scanargs := make([]interface{}, len(values))
	for i := range values {
		scanargs[i] = &values[i]
	}
	for rows.Next() {
		err := rows.Scan(scanargs...)
		if err != nil {
			fmt.Println(err)
		}
		record := make(Record)
		for i, v := range values {
			record[columns[i]] = v
		}
		records = append(records, record)
	}
	return records, nil
}
