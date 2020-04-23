package provider_test

import (
	"fmt"
	"testing"

	"github.com/hednowley/sound/provider"
)

func Test1(t *testing.T) {

	p, err := provider.NewBeetsProvider("beets", "../testdata/beetslib.blb")
	if err != nil {
		t.Errorf("Could not make provider: %v", err)
	}

	p.Iterate(func(path string) error {
		fmt.Println(path)
		i, err := p.GetInfo(path)
		if err != nil {
			t.Errorf("Error getting info for %v: %v", path, err)
		}
		fmt.Println(i.Title)

		return nil
	})

}
