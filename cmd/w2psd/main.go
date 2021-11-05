package main

import (
	"net"
	"sync"

	"google.golang.org/grpc"

	"github.com/fsnotify/fsnotify"

	"github.com/liupzmin/weewoe/log"
	pb "github.com/liupzmin/weewoe/proto"
	"github.com/liupzmin/weewoe/ssh"
	"github.com/spf13/viper"
)

const (
	UseCache = 0
	Reload   = 1
)

const (
	Bad  = 0
	Good = 1
)

type Target struct {
	Conn     *ssh.Connection
	BootTime int64
}

var (
	instances   map[string]Target
	processInfo Config
	imux, pmux  sync.RWMutex
)

type Process struct {
	Name    string
	Path    string
	Ports   []int
	Group   string
	Host    string
	PIDFile string
	Flag    string
}

type ProcessState struct {
	Process
	State         int
	StateDescribe string
	StartTime     int64
	Timestamp     int64
}

type ProcessConfig struct {
	Host    string
	Process []Process
}

type Config struct {
	Processes []ProcessConfig
}

func main() {
	initConfig()
	initConnection(processInfo)

	lis, err := net.Listen("tcp", ":9527")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterStateServer(s, &State{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("grpc serve failed: %v", err)
	}
}

func initConfig() {
	viper.SetConfigName("w2psd")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/w2psd")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/weewoe")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Fatal error config file: %w \n", err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Debugf("Config file changed:", e.Name)
		pmux.Lock()
		defer pmux.Unlock()
		loadProcessInfo()
		initConnection(processInfo)
	})
	viper.WatchConfig()

	loadProcessInfo()

	log.Debugf("the process is : %+v", processInfo)
}

func loadProcessInfo() {
	pro := viper.Sub("data")

	if pro == nil {
		log.Panicf("processes not found in config file")
	}

	err := pro.Unmarshal(&processInfo)
	if err != nil {
		log.Panicf("unable to decode into struct, %v", err)
	}
}

func initConnection(conf Config) {
	imux.Lock()
	defer imux.Unlock()
	if instances == nil {
		instances = make(map[string]Target)
	}
	for _, v := range conf.Processes {
		if _, ok := instances[v.Host]; !ok {
			conn, err := ssh.NewConnection(v.Host+":22", "root", true)
			if err != nil {
				log.Errorf("connect to %s failed: %s", v.Host, err.Error())
				continue
			}
			btime, err := getBootTime(conn)
			if err != nil {
				log.Errorf("get %s boot time failed: %s", v.Host, err.Error())
				continue
			}
			instances[v.Host] = Target{
				Conn:     conn,
				BootTime: btime,
			}
		}
	}
}
