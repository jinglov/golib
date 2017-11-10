package utils

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"time"
)

const (
	DATETIME_FORMAT = "2006-01-02 15:04:05"
)

func Byte2Int32(b []byte) int32 {
	buf := bytes.NewBuffer(b)
	var i int32
	binary.Read(buf, binary.BigEndian, &i)
	return i
}

func String2Int64(b string) int64 {
	i, error := strconv.Atoi(b)
	if error != nil {
		return 0
	}
	return int64(i)
}

func String2Int32(b string) int32 {
	i, error := strconv.Atoi(b)
	if error != nil {
		return 0
	}
	return int32(i)
}
func String2Int(b string) int {
	i, error := strconv.Atoi(b)
	if error != nil {
		return 0
	}
	return i
}

func String2Float32(b string) float32 {
	var f float64
	var err error
	f, err = strconv.ParseFloat(b, 32)
	if err != nil {
		return 0
	}
	return float32(f)
}

var DATETEIM = []string{
	"2006-01-02 15:04:05",
	"2006-01-02",
	"2006-1-2 15:4:5",
	"2006-1-2",
	"2006/01/02 15:04:05",
	"2006/01/02",
	"2006/1/2 15:4:5",
	"2006/1/2",
	"2006-01-02T15:04:05Z",
}

func FormatDateType(dateTime string) (tm time.Time, err error) {
	for _, tf := range DATETEIM {
		tm, err = time.Parse(tf, dateTime)
		if err == nil && !tm.IsZero() {
			return
		}
	}
	return
}
