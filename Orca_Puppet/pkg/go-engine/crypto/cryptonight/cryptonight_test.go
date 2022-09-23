package cryptonight

import (
	"testing"
)

func Test0000(t *testing.T) {
	if !TestSum("cn/0") {
		t.Error("TestSum fail")
	}
}

func Test0001(t *testing.T) {
	if !TestSum("cn/1") {
		t.Error("TestSum fail")
	}
}

func Test0002(t *testing.T) {
	if !TestSum("cn/2") {
		t.Error("TestSum fail")
	}
}

func Test0004(t *testing.T) {
	if !TestSum("cn/r") {
		t.Error("TestSum fail")
	}
}

func Test0005(t *testing.T) {
	if !TestSum("cn/fast") {
		t.Error("TestSum fail")
	}
}

func Test0006(t *testing.T) {
	if !TestSum("cn/half") {
		t.Error("TestSum fail")
	}
}

func Test0007(t *testing.T) {
	if !TestSum("cn/xao") {
		t.Error("TestSum fail")
	}
}

func Test0008(t *testing.T) {
	if !TestSum("cn/rto") {
		t.Error("TestSum fail")
	}
}

func Test0009(t *testing.T) {
	if !TestSum("cn/rwz") {
		t.Error("TestSum fail")
	}
}

func Test00010(t *testing.T) {
	if !TestSum("cn/zls") {
		t.Error("TestSum fail")
	}
}

func Test00011(t *testing.T) {
	if !TestSum("cn/double") {
		t.Error("TestSum fail")
	}
}

func Test00012(t *testing.T) {
	if !TestSum("cn-lite/0") {
		t.Error("TestSum fail")
	}
}

func Test00013(t *testing.T) {
	if !TestSum("cn-lite/1") {
		t.Error("TestSum fail")
	}
}

func Test00014(t *testing.T) {
	if !TestSum("cn-heavy/0") {
		t.Error("TestSum fail")
	}
}

func Test00015(t *testing.T) {
	if !TestSum("cn-heavy/tube") {
		t.Error("TestSum fail")
	}
}

func Test00016(t *testing.T) {
	if !TestSum("cn-heavy/xhv") {
		t.Error("TestSum fail")
	}
}

func Test00017(t *testing.T) {
	if !TestSum("cn-pico") {
		t.Error("TestSum fail")
	}
}

func Test00018(t *testing.T) {
	if !TestSum("cn-pico/tlo") {
		t.Error("TestSum fail")
	}
}
