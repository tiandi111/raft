package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tiandi111/raft/config"
	"github.com/tiandi111/raft/pkg/raft"
	graft "github.com/tiandi111/raft/pkg/rpc/raft"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	cfg     *raft.Config
	node    *raft.Node
	cfgfile string
	logpath string
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
	rand.Seed(time.Now().Unix())

	cobra.OnInitialize(initlog, initconfig)

	Command.PersistentFlags().StringVar(&cfgfile, "config", `.`, "init config")

	Command.PersistentFlags().StringVar(&logpath, "log", `.`, "log path")

	Command.PersistentFlags().Int32Var(&nodeId, "id", 0, "assign node id")
}

func initlog() {
	file := fmt.Sprintf(`%slog_%d`, logpath, nodeId)
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
}

func initconfig() {
	if cfgfile == "" {
		log.Printf("empty config file path, exit 1")
		os.Exit(1)
	}
	if nodeId <= 0 {
		log.Printf("invalid node id %d, exit 1", nodeId)
		os.Exit(1)
	}

	nlcfg, err := config.ParseConfig(cfgfile)
	if err != nil {
		panic(err)
	}

	cfg = &raft.Config{
		ID:     int32(nodeId),
		Others: make(map[int32]string),
	}

	for _, ncfg := range nlcfg {
		if ncfg.ID == nodeId {
			if ncfg.Addr == "" {
				panic("empty node address")
			}
			cfg.Addr = ncfg.Addr
			cfg.HeartbeatInterval = time.Duration(ncfg.HeartbeatInterval)
			cfg.HeartbeatCheckInterval = time.Duration(ncfg.HeartbeatCheckInterval)
			cfg.MaxElectionTimeout = time.Duration(ncfg.MaxElectionTimeout)
		} else {
			cfg.Others[ncfg.ID] = ncfg.Addr
		}
	}
}

func run() {
	defer func() {
	}()

	node = raft.NewNode(cfg)

	// init server
	errc := make(chan error, 1)
	lis, err := net.Listen("tcp", node.Config.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	graft.RegisterRaftServer(s, node)
	go func() {
		if err := s.Serve(lis); err != nil {
			errc <- err
		}
	}()

	// init client
	err = node.InitClient()
	if err != nil {
		panic(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM)

	go node.HeartbeatMonitor()

	go node.LeaderHeartbeater()

	for {
		select {
		case <-sigc:
			for _, client := range node.Clients {
				err := client.Close()
				if err != nil {
					log.Printf("close client failed, err : %s", err)
				}
			}
			node.DoneC <- struct{}{}
			return
		case serr := <-errc:
			log.Fatalf("server err:%s", serr)
			return
		default:
			node.Report()
			time.Sleep(time.Second)
		}
	}
}
