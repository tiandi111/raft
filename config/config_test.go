package config

import (
	"fmt"
	"github.com/tiandi111/ds/test"
	"testing"
)

func TestParseConfig(t *testing.T) {
	nlist, err := ParseConfig("config.yaml")
	test.AssertNil(t, err)
	fmt.Println(nlist)
}
