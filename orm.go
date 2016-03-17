package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
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
	groupby   string //groupby
	limit     string //limit
	join      string //join
	pk        string //pk
	dbname    string //dbname
}

//create new Module
func NewModule(tableName string) *Module {
	m := &Module{tableName: tableName, columnstr: "*", dbname: "default", pk: "id"}
	return m
}

func (m *Module) Clean() *Module {
	m.columnstr = "*"
	m.filters = ""
	m.orderby = ""
	m.limit = ""
	m.join = ""
	m.pk = "id"
	return m
}

func (m *Module) GetDB() *sql.DB {
	return dbHive[m.dbname]
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

func (m *Module) GroupBy(param string) *Module {
	m.groupby = fmt.Sprintf("GROUP BY %v", param)
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
	m.join += fmt.Sprintf(" LEFT JOIN %v ON %v", table, condition)
	return m
}

//rightJoin
func (m *Module) RightJoin(table, condition string) *Module {
	m.join += fmt.Sprintf(" RIGHT JOIN %v ON %v", table, condition)
	return m
}

//join
func (m *Module) Join(table, condition string) *Module {
	m.join += fmt.Sprintf(" INNER JOIN %v ON %v", table, condition)
	return m
}

//fulljoin
func (m *Module) FullJoin(table, condition string) *Module {
	m.join += fmt.Sprintf(" FULL JOIN %v ON %v", table, condition)
	return m
}

func (m *Module) getSqlString() string {
	columnstr := m.columnstr
	if l := len(columnstr); l > 1 {
		columnstr = columnstr[:l-1]
	}
	query := m.buildSql(columnstr)
	query += " " + m.limit
	log.Println("sql = ", query)
	return query
}

func (m *Module) buildSql(columnstr string) string {
	where := m.filters
	where = strings.TrimSpace(where)
	if len(where) > 0 {
		where = "where " + where
	}
	query := fmt.Sprintf("select %v from %v %v %v %v %v", columnstr, m.tableName, m.join, where, m.groupby, m.orderby)
	return query
}

func (m *Module) QueryPage(page *Page, callBackFunc func(*sql.Rows)) error {
	db := dbHive[m.dbname]
	m.Limit(page.StartRow(), page.PageSize)
	query := m.buildSql("count(*)")
	log.Println(query)
	row := db.QueryRow(query)
	err := row.Scan(&page.ResultCount)
	if err != nil {
		return err
	}
	rows, err := db.Query(m.getSqlString())
	if err != nil {
		return err
	}
	defer rows.Close()
	callBackFunc(rows)
	return nil
}

func (m *Module) Query(callBackFunc func(*sql.Rows)) error {
	db := dbHive[m.dbname]
	rows, err := db.Query(m.getSqlString())
	if err != nil {
		return err
	}
	defer rows.Close()
	callBackFunc(rows)
	return nil
}
func (m *Module) QueryOne(callBackFunc func(*sql.Row)) {
	db := dbHive[m.dbname]
	row := db.QueryRow(m.getSqlString())
	callBackFunc(row)
}

func (m *Module) IsExist() (bool, error) {
	count, err := m.Count()
	if count > 0 {
		return true, nil
	}
	return false, err
}

func (m *Module) Count() (int, error) {
	db := dbHive[m.dbname]
	query := m.buildSql("count(*)")
	log.Println("sql = ", query)
	row := db.QueryRow(query)
	var count int
	err := row.Scan(&count)
	return count, err
}

func (m *Module) OneRecord() (record Record, err error) {
	rs, err := m.Limit(1).AllRecords()
	if err != nil {
		return record, err
	}
	if len(rs) == 0 {
		return NewRecord(), errors.New("not fond record")
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
		record := NewRecord()
		for i, v := range values {
			record.result[columns[i]] = v
		}
		records = append(records, record)
	}
	return records, nil
}
func (m *Module) SetPK(pk string) *Module {
	m.pk = pk
	return m
}

func (m *Module) FindRecordById(id int) *Module {
	m.Filter(m.pk, "=", id)
	return m
}

func (m *Module) Insert(record Record) (int, error) {
	columns := ""
	values := ""
	for c, v := range record.param {
		columns = columns + c + ","
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool:
			values = values + fmt.Sprintf("%v", v) + ","
		default:
			values = values + fmt.Sprintf("'%v'", v) + ","
		}
	}
	if l := len(columns); l > 0 {
		columns = columns[:l-1]
	}
	if l := len(values); l > 0 {
		values = values[:l-1]
	}
	insertSql := fmt.Sprintf("insert into %v(%v) values(%v)", m.tableName, columns, values)
	fmt.Println(insertSql)
	db := dbHive[m.dbname]
	result, err := db.Exec(insertSql)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (m *Module) Update(record Record) error {
	values := ""
	for c, v := range record.param {
		values = values + c + "="
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool:
			values += fmt.Sprintf("%v", v)
		default:
			values += fmt.Sprintf("'%v'", v)
		}
		values += ","
	}
	if l := len(values); l > 0 {
		values = values[:l-1]
	}
	sql := fmt.Sprintf("update %v set %v where %v", m.tableName, values, m.filters)
	log.Println("sql = ", sql)
	db := dbHive[m.dbname]
	_, err := db.Exec(sql)
	return err
}

func (m *Module) DeleteById(id int) error {
	m.Filter(m.pk, "=", id)
	return m.Delete()
}

func (m *Module) FindById(id int) *Module {
	m.Filter(m.pk, "=", id)
	return m
}

func (m *Module) Delete() error {
	where := m.filters
	where = strings.TrimSpace(where)
	if len(where) > 0 {
		where = "where " + where
	}
	delSql := fmt.Sprintf("delete from %v %v", m.tableName, where)
	fmt.Println(delSql)
	_, err := dbHive[m.dbname].Exec(delSql)
	return err
}
