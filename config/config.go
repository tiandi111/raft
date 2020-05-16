package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type NodeConfig struct {
	ID                     int32  `yaml:"id"`
	Addr                   string `yaml:"addr"`
	HeartbeatInterval      int64  `yaml:"heartbeat_interval"`
	HeartbeatCheckInterval int64  `yaml:"heartbeat_check_interval"`
	MaxElectionTimeout     int64  `yaml:"max_election_timeout"`
}

type NodeListConfig []NodeConfig

func ParseConfig(file string) (NodeListConfig, error) {
	cbyte, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	nlist := []NodeConfig{}
	err = yaml.Unmarshal(cbyte, &nlist)
	if err != nil {
		return nil, err
	}
	return nlist, nil
}
