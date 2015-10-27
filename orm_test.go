package orm

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func Test_FindAll(t *testing.T) {
	NewDatabase("default", "mysql", "root:123456@tcp(localhost:3306)/goa?charset=utf8&parseTime=true&loc=Local")
	m := NewModule("user")
	r := new(Record)
	m.Select("id", "name", "loginName").Filter("name='", "zhangsan", "' and id=", 1).Limit(0, 10).FindAll(r)
}
