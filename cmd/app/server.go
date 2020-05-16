package app

import (
	"github.com/spf13/cobra"
	"github.com/tiandi111/raft/config"
	"github.com/tiandi111/raft/pkg/raft"
	graft "github.com/tiandi111/raft/pkg/rpc/raft"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	cfg     raft.Config
	node    *raft.Node
	cfgfile string
	nodeId  int32
	Command = &cobra.Command{
		Use:  "raft",
		Long: "my raft implementation",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	cobra.OnInitialize(initconfig)

	cmd.PersistentFlags().StringVar(&cfgfile, "config", "", "init config")

	cmd.PersistentFlags().Int32Var(&nodeId, "node_id", -1, "assign node id")
}

func initconfig() {
	if cfgfile == "" {
		log.Printf("empty config file path, exit 1")
		os.Exit(1)
	}
	if nodeId < 0 {
		log.Printf("invalid node id, exit 1")
		os.Exit(1)
	}

	nlist, err := config.ParseConfig(cfgfile)
	if err != nil {
		panic(err)
	}

	cfg := &raft.Config{
		ID:     int32(nodeId),
		Others: make(map[int32]string),
	}

	for _, ncfg := range nlist {
		if ncfg.ID == nodeId {
			if ncfg.Addr == "" {
				panic("empty node address")
			}
			cfg.Addr = ncfg.Addr
		} else {
			cfg.Others[ncfg.ID] = ncfg.Addr
		}
	}
}

func run() {
	defer func() {
	}()

	node = raft.NewNode(&cfg)

	err := node.InitClient()
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", node.Config.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	graft.RegisterRaftServer(s, node)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM)

	errc := make(chan error, 1)

	go func() {
		if err := s.Serve(lis); err != nil {
			errc <- err
		}
	}()

	go node.HeartbeatMonitor()

	go node.LeaderHeartbeater()

	select {
	case <-sigc:
		for _, client := range node.Clients {
			err := client.Close()
			if err != nil {
				log.Printf("close client failed, err : %s", err)
			}
		}
	case serr := <-errc:
		log.Fatalf("server err:%s", serr)
	}
}
