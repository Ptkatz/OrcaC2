package common

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"math"
	"reflect"
	"strconv"
	"strings"
)

const gcharset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func IntArrayToString(a []int, delim string) string {
	var ret string
	for _, s := range a {
		ret += strconv.Itoa(s) + delim
	}
	return ret
}

func Int32ArrayToString(a []int32, delim string) string {
	var ret string
	for _, s := range a {
		ret += strconv.Itoa((int)(s)) + delim
	}
	return ret
}

func Int64ArrayToString(a []int64, delim string) string {
	var ret string
	for _, s := range a {
		ret += strconv.Itoa((int)(s)) + delim
	}
	return ret
}

func RandStr(l int) string {
	b := make([]byte, l)
	for i := range b {
		b[i] = gcharset[RandInt31n(len(gcharset))]
	}
	return string(b)
}

type StrTableLine struct {
	cols []string
}

type StrTable struct {
	header []string
	lines  []StrTableLine
}

func (s *StrTable) AddHeader(h string) {
	s.header = append(s.header, h)
}

func (s *StrTable) AddLine(l StrTableLine) {
	s.lines = append(s.lines, l)
}

func (s *StrTableLine) AddData(d string) {
	s.cols = append(s.cols, d)
}

func (s *StrTable) String(prefix string) string {

	if len(s.header) <= 0 {
		return ""
	}

	colmax := make([]int, 0)
	for _, s := range s.header {
		colmax = append(colmax, len(s))
	}

	totalcol := 0
	for i := 0; i < len(colmax); i++ {
		max := colmax[i]
		for _, sl := range s.lines {
			if i < len(sl.cols) {
				max = MaxOfInt(max, len(sl.cols[i]))
			}
		}
		colmax[i] = max
		totalcol += max
	}
	totalcol += len(colmax) + 1

	/*
		-----------
		| a  | b  |
		-----------
		| 1  | 2  |
		-----------
	*/

	ret := prefix
	ret += strings.Repeat("-", totalcol) + "\n" + prefix
	for i, h := range s.header {
		ret += "|" + WrapString(h, colmax[i])
	}
	ret += "|" + "\n" + prefix

	for _, l := range s.lines {
		ret += strings.Repeat("-", totalcol) + "\n" + prefix
		for i, d := range l.cols {
			ret += "|" + WrapString(d, colmax[i])
		}
		for i := len(l.cols); i < len(colmax); i++ {
			ret += "|" + WrapString("", colmax[i])
		}
		ret += "|" + "\n" + prefix
	}

	ret += strings.Repeat("-", totalcol) + "\n"

	return ret
}

func (s *StrTable) FromStruct(v interface{}, use func(name string) bool) {
	ss := reflect.ValueOf(v).Elem()
	typeOfT := ss.Type()

	for i := 0; i < ss.NumField(); i++ {
		name := typeOfT.Field(i).Name
		if use != nil {
			if !use(name) {
				continue
			}
		}
		s.AddHeader(name)
	}
}

func (s *StrTableLine) FromStruct(st *StrTable, v interface{}, trans func(name string, v interface{}) interface{}) {
	ss := reflect.ValueOf(v).Elem()
	typeOfT := ss.Type()

	for i := 0; i < ss.NumField(); i++ {
		f := ss.Field(i)
		name := typeOfT.Field(i).Name

		if !ArrayContainString(st.header, name) {
			continue
		}

		v := f.Interface()
		if trans != nil {
			v = trans(name, f.Interface())
		}
		if v != nil {
			str := fmt.Sprintf("%v", v)
			s.AddData(str)
		} else {
			s.AddData("")
		}
	}
}

func WrapString(s string, n int) string {
	if n <= len(s) {
		return s
	}
	l := (n - len(s)) / 2
	r := (n - len(s)) - l
	return strings.Repeat(" ", l) + s + strings.Repeat(" ", r)
}

func StructToTable(v interface{}) string {
	t := StrTable{}
	tl := StrTableLine{}
	t.FromStruct(v, nil)
	tl.FromStruct(&t, v, nil)
	t.AddLine(tl)
	return t.String("")
}

const num2char string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	LITTLE_LETTERS = 36
	FULL_LETTERS   = 62
)

func NumToHex(num, n int) string {
	num_str := ""
	for num != 0 {
		yu := num % n
		num_str = string(num2char[yu]) + num_str
		num = num / n
	}
	return num_str
}

func Hex2Num(str string, n int) int {
	v := 0.0
	length := len(str)
	for i := 0; i < length; i++ {
		s := string(str[i])
		index := strings.Index(num2char, s)
		v += float64(index) * math.Pow(float64(n), float64(length-1-i))
	}
	return int(v)
}

func GzipString(data string) string {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write([]byte(data))
	w.Close()
	return string(b.Bytes())
}

func GzipStringBestCompression(data string) string {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return ""
	}
	w.Write([]byte(data))
	w.Close()
	return string(b.Bytes())
}

func GzipStringBestSpeed(data string) string {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return ""
	}
	w.Write([]byte(data))
	w.Close()
	return string(b.Bytes())
}

func GunzipString(data string) string {
	var b bytes.Buffer
	b.WriteString(data)
	r, err := gzip.NewReader(&b)
	if err != nil {
		return ""
	}
	r.Close()
	ret, _ := ioutil.ReadAll(r)
	return string(ret)
}
