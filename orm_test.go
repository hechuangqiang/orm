package orm

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func Test_FindAll(t *testing.T) {
	NewDatabase("default", "mysql", "root:123456@tcp(localhost:3306)/goa?charset=utf8&parseTime=true&loc=Local")
	m := NewModule("user u")
	r := make([]Record, 0)
	m.Select("u.id", "u.name", "u.loginName", "a.token", "a.id AuthId").Join("auth a", "a.userId=u.id").Filter("u.loginName='", "admin", "'").Limit(0, 10).FindAll(r)
	t.Log(r)
}

func Test_GetRecords(t *testing.T) {
	NewDatabase("default", "mysql", "root:123456@tcp(localhost:3306)/goa?charset=utf8&parseTime=true&loc=Local")
	m := NewModule("user u")
	records, err := m.Select("id", "name", "loginName").Filter("loginName='admin'").GetRecords()
	t.Log(records, err)
}
