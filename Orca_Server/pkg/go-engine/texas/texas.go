package texas

import (
	"bufio"
	"Orca_Server/pkg/go-engine/common"
	"Orca_Server/pkg/go-engine/loggo"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Load() {

	loggo.Info("texas Load start")

	err := LoadNormalColor(false)
	if err != nil {
		panic(err)
	}
	err = LoadProbility(false)
	if err != nil {
		panic(err)
	}
}

func LoadLocal() {

	loggo.Info("texas LoadLocal start")

	err := LoadNormalColor(true)
	if err != nil {
		panic(err)
	}
	err = LoadProbility(true)
	if err != nil {
		panic(err)
	}
}

func LoadNormalColor(local bool) error {

	var path string
	if !local {
		path = common.GetDataDir() + "/texas/"
	}

	err := LoadColor(path + "texas_data_color.txt")
	if err != nil {
		return err
	}
	err = LoadNormal(path + "texas_data_normal.txt")
	if err != nil {
		return err
	}
	err = LoadColor(path + "texas_data_extra_color_6.txt")
	if err != nil {
		return err
	}
	err = LoadNormal(path + "texas_data_extra_normal_6.txt")
	if err != nil {
		return err
	}
	err = LoadColor(path + "texas_data_extra_color_5.txt")
	if err != nil {
		return err
	}
	err = LoadNormal(path + "texas_data_extra_normal_5.txt")
	if err != nil {
		return err
	}
	return nil
}
func LoadProbility(local bool) error {

	var path string
	if !local {
		path = common.GetDataDir() + "/texas/"
	}

	for i := 6; i >= 2; i-- {
		err := loadProbility(i, path+"texas_data_opt_"+strconv.Itoa(i)+".txt")
		if err != nil {
			return err
		}
	}
	return nil
}

type KeyData struct {
	index   int
	postion int
	max     int64
	ty      int
}

var colorMap map[int64]KeyData
var normalMap map[int64]KeyData

type ProbilityData struct {
	avg float32
	min float32
	max float32
}

var probilityMap [7]map[int64]ProbilityData
var optprobilityMap [7]map[int64]ProbilityData

func loadProbility(i int, file string) error {

	file = filepath.Clean(file)
	file = filepath.ToSlash(file)

	loggo.Info("start loadProbility %v", file)

	if probilityMap[i] == nil {
		probilityMap[i] = make(map[int64]ProbilityData)
	}
	if optprobilityMap[i] == nil {
		optprobilityMap[i] = make(map[int64]ProbilityData)
	}

	f, err := os.Open(file)
	if err != nil {
		loggo.Error("loadProbility %v %s", file, err)
		return err
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}

		line = strings.TrimSpace(line)
		params := strings.Split(line, " ")

		key, _ := strconv.ParseInt(params[0], 10, 64)
		ty, _ := strconv.ParseInt(params[1], 10, 32)
		probility, _ := strconv.ParseFloat(params[2], 32)
		min, _ := strconv.ParseFloat(params[3], 32)
		max, _ := strconv.ParseFloat(params[4], 32)

		if ty == 0 {
			probilityMap[i][key] = ProbilityData{float32(probility), float32(min), float32(max)}
		} else {
			optprobilityMap[i][key] = ProbilityData{float32(probility), float32(min), float32(max)}
		}
	}

	loggo.Info("end loadProbility %v", file)

	return nil
}

func LoadNormal(file string) error {

	file = filepath.Clean(file)
	file = filepath.ToSlash(file)

	loggo.Info("start LoadNormal %v", file)

	if normalMap == nil {
		normalMap = make(map[int64]KeyData)
	}

	f, err := os.Open(file)
	if err != nil {
		loggo.Error("LoadNormal %v %s", file, err)
		return err
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}

		line = strings.TrimSpace(line)
		params := strings.Split(line, " ")

		key, _ := strconv.ParseInt(params[0], 10, 64)
		i, _ := strconv.ParseInt(params[1], 10, 32)
		index, _ := strconv.ParseInt(params[2], 10, 32)
		max, _ := strconv.ParseInt(params[5], 10, 64)
		ty, _ := strconv.ParseInt(params[7], 10, 32)

		keyData := KeyData{int(index), int(i), max, int(ty)}
		normalMap[key] = keyData
	}

	loggo.Info("end LoadNormal %v", file)

	return nil
}

func LoadColor(file string) error {

	file = filepath.Clean(file)
	file = filepath.ToSlash(file)

	loggo.Info("start LoadColor %v", file)

	if colorMap == nil {
		colorMap = make(map[int64]KeyData)
	}

	f, err := os.Open(file)
	if err != nil {
		loggo.Error("LoadColor %v %s", file, err)
		return err
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}

		line = strings.TrimSpace(line)
		params := strings.Split(line, " ")

		key, _ := strconv.ParseInt(params[0], 10, 64)
		i, _ := strconv.ParseInt(params[1], 10, 32)
		index, _ := strconv.ParseInt(params[2], 10, 32)
		max, _ := strconv.ParseInt(params[5], 10, 64)
		ty, _ := strconv.ParseInt(params[7], 10, 32)

		keyData := KeyData{int(index), int(i), max, int(ty)}
		colorMap[key] = keyData
	}

	loggo.Info("end LoadColor %v", file)

	return nil
}

type Poke struct {
	color int8
	value int8
}

var huaseName = []string{"方", "梅", "红", "黑"}
var valueName = []string{"", "", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

func (p *Poke) ToByte() int8 {
	return (int8)(p.color<<4 | p.value)
}

func (p *Poke) ToString() string {
	if p.value == PokeValue_GUI && p.color == PokeColor_GUI {
		return "鬼"
	}
	return huaseName[p.color] + valueName[p.value]
}

func NewPoke(byteValue int8) Poke {
	return Poke{(int8)(byteValue >> 4), (int8)(byteValue % 16)}
}

const (
	PokeColor_FANG = 0
	PokeColor_MEI  = 1
	PokeColor_HONG = 2
	PokeColor_HEI  = 3
	PokeColor_GUI  = 5
)

const (
	PokeValue_2   = 2
	PokeValue_3   = 3
	PokeValue_4   = 4
	PokeValue_5   = 5
	PokeValue_6   = 6
	PokeValue_7   = 7
	PokeValue_8   = 8
	PokeValue_9   = 9
	PokeValue_10  = 10
	PokeValue_J   = 11
	PokeValue_Q   = 12
	PokeValue_K   = 13
	PokeValue_A   = 14
	PokeValue_GUI = 8
)

const (
	TEXAS_CARD_TYPE_GAOPAI          = 1  //高牌
	TEXAS_CARD_TYPE_DUIZI           = 2  //对子
	TEXAS_CARD_TYPE_LIANGDUI        = 3  //两对
	TEXAS_CARD_TYPE_SANTIAO         = 4  //三条
	TEXAS_CARD_TYPE_SHUNZI          = 5  //顺子
	TEXAS_CARD_TYPE_TONGHUA         = 6  //同花
	TEXAS_CARD_TYPE_HULU            = 7  //葫芦
	TEXAS_CARD_TYPE_SITIAO          = 8  //四条
	TEXAS_CARD_TYPE_TONGHUASHUN     = 9  //同花顺
	TEXAS_CARD_TYPE_KINGTONGHUASHUN = 10 //皇家同花顺
)

var winName = []string{"无",
	"高牌",
	"对子",
	"两对",
	"三条",
	"顺子",
	"同花",
	"葫芦",
	"四条",
	"同花顺",
	"皇家同花顺",
	"MAX"}

var GUI = Poke{PokeColor_GUI, PokeValue_GUI}

const (
	GUINUM = 2
	GENNUM = 52 + GUINUM
)

var allCards = genAllCards()

func genAllCards() []int8 {
	var ret []int8
	for i := 0; i < 4; i++ {
		for j := 0; j < 52/4; j++ {
			p := Poke{int8(i), int8(j + PokeValue_2)}
			ret = append(ret, p.ToByte())
		}
	}
	return ret
}

func GetWinNameId(n string) int {
	for i, p := range winName {
		if p == n {
			return i
		}
	}
	return 0
}

func GetWinName(id int) string {
	return winName[id]
}
