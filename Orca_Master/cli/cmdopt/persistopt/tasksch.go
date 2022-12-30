package persistopt

import "time"

func stringTotime(timemee string) (t time.Time) {
	timetemp := "2006-01-02 1a5:04:05"
	tm, err := time.ParseInLocation(timetemp, timemee, time.Local)
	if err != nil {
		panic(err)
	}
	return tm
}

func CheckTime(str string) (t time.Time) {
	if str != "" {
		t := stringTotime(str)
		return t
	}
	tm := time.Now().Add(time.Minute * +1)
	return tm
}
