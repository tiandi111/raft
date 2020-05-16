package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type NodeConfig struct {
	ID   int32  `yaml:"id"`
	Addr string `yaml:"addr"`
}

type NodeList []*NodeConfig

func ParseConfig(file string) (NodeList, error) {
	cbyte, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var nlist NodeList
	err = yaml.Unmarshal(cbyte, nlist)
	if err != nil {
		return nil, err
	}
	return nlist, nil
}
