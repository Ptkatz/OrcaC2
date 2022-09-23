package common

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func elapsed() {
	defer Elapsed(func(d time.Duration) {
		fmt.Println("use time " + d.String())
	})()

	time.Sleep(time.Second)
}

func Test0001(t *testing.T) {

	a := RandStr(5)
	a1 := RandStr(5)
	fmt.Println(a)
	fmt.Println(a1)

	fmt.Println(GetOutboundIP())

	fmt.Println(GetNowUpdateInSecond())

	d, _ := Rc4("123456", []byte("asdgdsagdsag435t43321dsgesg"))
	fmt.Println(string(d))

	d, _ = Rc4("123456", d)
	fmt.Println(string(d))

	dd := MAKEINT64(12345, 7890)
	fmt.Println(dd)
	fmt.Println(HIINT32(dd))
	fmt.Println(LOINT32(dd))
	ddd := MAKEINT32(12345, 7890)
	fmt.Println(ddd)
	fmt.Println(HIINT16(ddd))
	fmt.Println(LOINT16(ddd))

	fmt.Println(IsInt(3.0002))
	fmt.Println(IsInt(3))
	fmt.Println(strconv.FormatFloat(3.1415, 'E', -1, 64))

	aa := []int{1, 2, 3, 4, 5, 6, 7, 8}
	Shuffle(len(aa), func(i, j int) { aa[i], aa[j] = aa[j], aa[i] })
	fmt.Println(aa)

	fmt.Println(RandInt())
	fmt.Println(RandInt31n(10))

	fmt.Println(WrapString("abc", 10))

	ts := StrTable{}
	ts.AddHeader("a")
	ts.AddHeader("b")
	ts.AddHeader("c")
	tsl := StrTableLine{}
	tsl.AddData("1234")
	tsl.AddData("123421412")
	ts.AddLine(tsl)
	tsl = StrTableLine{}
	tsl.AddData("aaa")
	ts.AddLine(tsl)
	fmt.Println(WrapString("abc", 10))
	fmt.Println(ts.String("\t"))

	elapsed()
}

type TestStruct struct {
	A int
	B int64
	C string
}

func Test0002(t *testing.T) {
	ts := TestStruct{1, 2, "3"}
	st := StrTable{}
	st.AddHeader("AA")
	st.FromStruct(&ts, func(name string) bool {
		return name != "A"
	})
	stl := StrTableLine{}
	stl.AddData("a")
	stl.FromStruct(&st, &ts, func(name string, v interface{}) interface{} {
		if name == "B" {
			return time.Duration(v.(int64)).String()
		}
		return v
	})
	st.AddLine(stl)
	ts = TestStruct{12, 214124, "124123"}
	stl = StrTableLine{}
	stl.AddData("aaa")
	stl.FromStruct(&st, &ts, func(name string, v interface{}) interface{} {
		if name == "B" {
			return time.Duration(v.(int64)).String()
		}
		return v
	})
	st.AddLine(stl)
	fmt.Println(st.String(""))

	SaveJson("test.json", &ts)
	ts1 := TestStruct{}
	err := LoadJson("test.json", &ts1)
	fmt.Println(err)
	fmt.Println(ts1.C)
}

func Test0003(t *testing.T) {
	a := NumToHex(12345745643, LITTLE_LETTERS)
	b := NumToHex(12345745643, FULL_LETTERS)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(Hex2Num(a, LITTLE_LETTERS))
	fmt.Println(Hex2Num(b, FULL_LETTERS))
	aa := NumToHex(37, LITTLE_LETTERS)
	bb := NumToHex(37, FULL_LETTERS)
	fmt.Println(aa)
	fmt.Println(bb)
	cc := Hex2Num("1i39pJZR", FULL_LETTERS)
	fmt.Println(cc)
	fmt.Println(NumToHex(cc, FULL_LETTERS))
	fmt.Println(NumToHex(cc+1, FULL_LETTERS))

	dd := Hex2Num("ZZZZZZZZ", FULL_LETTERS)
	fmt.Println(dd)
	fmt.Println(NumToHex(dd, FULL_LETTERS))
}

type TestStruct1 struct {
	TestStruct
	D int64
}

func Test0004(t *testing.T) {
	ts := TestStruct{1, 2, "3"}
	ts1 := TestStruct1{ts, 3}
	fmt.Println(StructToTable(&ts1))
}

func Test0005(t *testing.T) {
	fmt.Println(GetXXHashString("1"))
	fmt.Println(GetXXHashString("2"))
	fmt.Println(GetXXHashString("asfaf"))
	fmt.Println(GetXXHashString("dffd43321"))
}

func Test0006(t *testing.T) {
	src := "safa3232sgsgd343q421dsdgsddsgsarwdsddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddgdsgsgewrngfxgcjfhrsgcbxgfhreu658545ghuj,hgfdtsz nsdtzbntshjtwg,tu523jlikr[]iwsfffffds23525ewfsu45632rqwfsgrxy4353rsfzshrey4324fASffdjftui4e22=-"
	fmt.Println(len(src))
	a := GzipString(src)
	fmt.Println(len(a))
	b := GzipStringBestCompression(src)
	fmt.Println(len(b))
	c := GzipStringBestSpeed(src)
	fmt.Println(len(c))

	if src != GunzipString(a) {
		t.Error("fail")
	}
	if src != GunzipString(b) {
		t.Error("fail")
	}
	if src != GunzipString(c) {
		t.Error("fail")
	}
}

func Test0007(t *testing.T) {
	fmt.Println(GetCrc32String(""))
	fmt.Println(GetCrc32String("1"))
	fmt.Println(GetCrc32String("2"))
	fmt.Println(GetCrc32String("asfsadgewwe"))
}

func Test0008(t *testing.T) {
	c := NewChannel(10)
	c.Write(1)
	i := <-c.Ch()
	fmt.Println(i)
	c.Close()
	c.Close()
	c.Write(1)
	c.Write(1)
}

func Test0009(t *testing.T) {
	c := NewChannel(1)
	c.Write(1)
	fmt.Println(c.WriteTimeout(1, 1000))
	fmt.Println(c.WriteTimeout(1, 1000))
	i := <-c.ch
	fmt.Println(i)
	fmt.Println(c.WriteTimeout(1, 1000))
	time.Sleep(time.Second)
}

func Test0010(t *testing.T) {
	a := make([]int, 3)
	a[0] = 1
	a[1] = 111
	a[2] = 1111
	fmt.Println(HasInt(a, 1))
	fmt.Println(HasInt(a, 12))
}

func Test0011(t *testing.T) {
	a := make([]string, 3)
	a[0] = "1"
	a[1] = "111"
	a[2] = "1111"
	fmt.Println(HasString(a, "1"))
	fmt.Println(HasString(a, "12"))
}

func Test0012(t *testing.T) {
	Copy("common.go", "common.go.1")
	fmt.Println(FileExists("common.go.1"))
	fmt.Println(FileMd5("common.go.1"))
	fmt.Println(FileReplace("common.go.1", "func", "fuck"))
	fmt.Println(IsSymlink("common.go.1"))
	fmt.Println(FileFind("common.go.1", "fuck ini()"))
	fmt.Println(FileLineCount("common.go.1"))
}

func Test0013(t *testing.T) {
	Walk("./", func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	})
}

func Test0014(t *testing.T) {
	fmt.Println(NearlyEqual(1, 10))
	fmt.Println(NearlyEqual(8, 10))
	fmt.Println(NearlyEqual(9, 10))
	fmt.Println(NearlyEqual(99, 100))
	fmt.Println(NearlyEqual(90000, 100000))
	fmt.Println(NearlyEqual(80000, 100000))
}

func Test0015(t *testing.T) {
	fmt.Println("start")
	Sleep(3)
	fmt.Println("end")
}
