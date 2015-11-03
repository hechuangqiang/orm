package orm

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	NewDatabase("default", "mysql", "root:123456@tcp(localhost:3306)/goa?charset=utf8&parseTime=true&loc=Local")
}

func Test_FindAll(t *testing.T) {
	m := NewModule("user u")
	r := make([]Record, 0)
	m.Select("u.id", "u.name", "u.loginName", "a.token", "a.id AuthId").Join("auth a", "a.userId=u.id").Filter("u.loginName='", "admin", "'").Limit(0, 10).FindAll(r)
	t.Log(r)
}

func Test_GetRecords(t *testing.T) {
	m := NewModule("user u")
	records, _ := m.Select("id", "name", "loginName").Filter("1=1").OrderBy("id desc").AllRecords()
	if len(records) != 2 {
		t.Fail()
	}
}

func Test_OneRecord(t *testing.T) {
	m := NewModule("user")
	record, _ := m.Select("id,name,loginName").Filter("loginName='zs'").OneRecord()
	if record == nil {
		t.Fail()
	}
	if record.GetInt("id") != 2 {
		t.Fail()
	}
	if record.GetString("name") != "zhangsan" {
		t.Fail()
	}
}
