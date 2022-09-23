package loggo

import (
	"testing"
	"time"
)

func Test0001(t *testing.T) {
	Ini(Config{
		Level:  LEVEL_INFO,
		Prefix: "test",
		MaxDay: 3,
	})

	for i := 0; i < 10; i++ {
		Info("test %d", i)
		time.Sleep(time.Second)
	}

}
