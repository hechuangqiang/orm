package orm

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"math"
)

var bigEndian bool

func init() {
	var i int32 = 0x12345678
	var b byte = byte(i)
	if b == 0x12 {
		bigEndian = true
	}
	bigEndian = false
}

type Record map[string]sql.RawBytes

func (r Record) GetString(field string) string {
	return string(r[field])
}

func (r Record) GetBytes(key string) []byte {
	return []byte(r[key])
}

func (r Record) GetInt(key string) int {
	buf := bytes.NewBuffer(r[key])
	var x int32
	if bigEndian {
		binary.Read(buf, binary.BigEndian, &x)
	} else {
		binary.Read(buf, binary.LittleEndian, &x)
	}
	return int(x)
}

func (r Record) GetInt64(key string) int64 {
	buf := bytes.NewBuffer(r[key])
	var x int64
	if bigEndian {
		binary.Read(buf, binary.BigEndian, &x)
	} else {
		binary.Read(buf, binary.LittleEndian, &x)
	}
	return int64(x)
}

func (r Record) GetFloat32(key string) float32 {
	var v uint32
	if bigEndian {
		v = binary.BigEndian.Uint32(r[key])
	} else {
		v = binary.LittleEndian.Uint32(r[key])
	}
	f := math.Float32frombits(v)
	return f
}
