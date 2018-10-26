package util_test

import (
	"testing"

	"github.com/hednowley/sound/util"
)

func TestMin(t *testing.T) {

	m := util.Min(100, 80)
	if m != 80 {
		t.Error()
	}

	m = util.Min(-10, -20)
	if m != -20 {
		t.Error()
	}

	m = util.Min(-5, -5)
	if m != -5 {
		t.Error()
	}
}

func TestMax(t *testing.T) {

	m := util.Max(100, 80)
	if m != 100 {
		t.Error()
	}

	m = util.Max(-10, -20)
	if m != -10 {
		t.Error()
	}

	m = util.Max(-5, -5)
	if m != -5 {
		t.Error()
	}
}

func TestContains(t *testing.T) {

	s := []uint{}
	if util.Contains(s, 0) {
		t.Error()
	}

	s = []uint{5, 80}
	if util.Contains(s, 10) {
		t.Error()
	}
	if !util.Contains(s, 80) {
		t.Error()
	}
}
