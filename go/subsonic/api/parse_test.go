package api

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/util"
)

func TestParseReponseFormat(t *testing.T) {

	if defaultFormat != xmlFormat {
		t.Error("Default format should be XML")
	}

	format := parseResponseFormat("")
	if *format != defaultFormat {
		t.Error()
	}

	format = parseResponseFormat("xMl")
	if *format != xmlFormat {
		t.Error()
	}

	format = parseResponseFormat("jSon")
	if *format != jsonFormat {
		t.Error()
	}

	format = parseResponseFormat("jSonX")
	if format != nil {
		t.Error()
	}
}

func TestParseUint(t *testing.T) {

	i := util.ParseUint("", 56)
	if i != 56 {
		t.Error()
	}

	i = util.ParseUint("dfhgfdh", 435)
	if i != 435 {
		t.Error()
	}

	i = util.ParseUint("344343", 435)
	if i != 344343 {
		t.Error()
	}

	i = util.ParseUint("-57", 435)
	if i != 435 {
		t.Error()
	}
}

func TestParseBool(t *testing.T) {

	b := util.ParseBool("")
	if b != nil {
		t.Error()
	}

	b = util.ParseBool("faalse")
	if b != nil {
		t.Error()
	}

	b = util.ParseBool("1")
	if !*b {
		t.Error()
	}

	b = util.ParseBool("tRue")
	if !*b {
		t.Error()
	}

	b = util.ParseBool("0")
	if *b {
		t.Error()
	}

	b = util.ParseBool("faLse")
	if *b {
		t.Error()
	}

}

func TestParseParams(t *testing.T) {
	params := url.Values{}
	params.Add("a", "aa")
	params.Add("b", "bb")
	params.Add("c", "c1")
	params.Add("c", "c2")
	params.Add("d", "d1")
	params.Add("d", "d2")

	body := "f=ff&e=e1&e=e2&a=a1&c=c3"

	p := parseParams(params, []byte(body))

	a := p["a"]
	if len(a) != 2 || a[0] != "aa" || a[1] != "a1" {
		t.Error()
	}

	b := p["b"]
	if len(b) != 1 || b[0] != "bb" {
		t.Error()
	}

	c := p["c"]
	if len(c) != 3 || c[0] != "c1" || c[1] != "c2" || c[2] != "c3" {
		t.Error()
	}

	d := p["d"]
	if len(d) != 2 || d[0] != "d1" || d[1] != "d2" {
		t.Error()
	}

	e := p["e"]
	if len(e) != 2 || e[0] != "e1" || e[1] != "e2" {
		t.Error()
	}

	f := p["f"]
	if len(f) != 1 || f[0] != "ff" {
		t.Error()
	}

	g := p["g"]
	if len(g) != 0 {
		t.Error()
	}

}
