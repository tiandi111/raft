#!/bin/bash

path=$GOPATH/src/github.com/tiandi111/raft
cfg_path=$path/config/config.yaml
log_path=$path/output/

mkdir $log_path

export GOBIN=$GOPATH/bin

go install "$path"/cmd/raft.go

start_cluster () {

  nohup raft --config="$cfg_path" --id=1 --log="$log_path"

  nohup raft --config="$cfg_path" --id=2 --log="$log_path"

  nohup raft --config="$cfg_path" --id=3 --log="$log_path"

}


raft --id=3 --config=$GOPATH/src/github.com/tiandi111/raft/config/config.yaml --log=$GOPATH/src/github.com/tiandi111/raft/output/