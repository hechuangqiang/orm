package orm

import (
	"fmt"
	"reflect"
	"strconv"
)

type Record map[string][]byte

func (r Record) GetString(field string) string {
	return string(r[field])
}

func (r Record) GetBytes(key string) []byte {
	return r[key]
}

func (r Record) GetInt(key string) int {
	var x int
	s := asString(r[key])
	i64, err := strconv.ParseInt(s, 10, reflect.TypeOf(x).Bits())
	if err != nil {
		panic(err.Error())
	}
	x = int(i64)
	return x
}

func (r Record) GetInt64(key string) int64 {
	s := asString(r[key])
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	return i64
}

func (r Record) GetFloat32(key string) float32 {
	s := asString(r[key])
	f64, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err.Error())
	}
	return float32(f64)
}

func (r Record) GetFloat64(key string) float64 {
	s := asString(r[key])
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err.Error())
	}
	return f64
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}
