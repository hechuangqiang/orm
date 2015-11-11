package orm

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	NewDatabase("default", "mysql", "root:123456@tcp(localhost:3306)/goa?charset=utf8&parseTime=true&loc=Local")
}

func Test_FindAll(t *testing.T) {
	m := NewModule("user u")
	r := make([]interface{}, 0)
	m.Select("u.id", "u.name", "u.loginName", "a.token", "a.id AuthId").Join("auth a", "a.userId=u.id").Filter("u.loginName='", "admin", "'").Limit(0, 10).QueryOne(func(row *sql.Row) {
		row.Scan(r...)
	})
	t.Log(r)
}

func Test_GetRecords(t *testing.T) {
	m := NewModule("user u")
	records, _ := m.Select("id", "name", "loginName").Filter("1=1").OrderBy("id desc").AllRecords()
	t.Log(len(records))
}

func Test_OneRecord(t *testing.T) {
	m := NewModule("user")
	record, _ := m.Select("id,name,loginName").Filter("loginName='zs'").OneRecord()
	if record.GetInt("id") != 2 {
		t.Fail()
	}
	if record.GetString("name") != "zhangsan" {
		t.Fail()
	}
}
func Test_Insert(t *testing.T) {
	record := NewRecord()
	record.Set("name", "lisi")
	record.Set("loginName", "lisi")
	record.Set("pwd", "123321")
	record.Set("createTime", time.Now())
	m := NewModule("user")
	err := m.Insert(record)
	t.Log("err = ", err)
}
