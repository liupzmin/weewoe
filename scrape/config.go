package scrape

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/liupzmin/weewoe/log"
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

const TimeLayout = "2006-01-02 15:04:05"

type Target struct {
	Conn *ssh.Connection
	// used to compute process start time
	BootTime int64
}

func (t Target) Close() {
	t.Conn.Close()
}

var (
	instances   map[string]Target
	processInfo Config
	imux, pmux  sync.RWMutex
	pro         *viper.Viper
)

type Process struct {
	Name    string
	Path    string
	Ports   []int64
	Group   string
	Host    string
	PIDFile string
	Flag    string
	Suspend bool
}

type ProcessState struct {
	Process
	State         int64
	StateDescribe string
	StartTime     int64
	Timestamp     int64
}

type PortState struct {
	Process
	States    []*Port
	Timestamp int64
}

type Port struct {
	Number string
	State  int64
}

type ProcessConfig struct {
	Host    string
	Process []Process
}

type Config struct {
	Processes []ProcessConfig
}

func Init() {
	initConfig()
	initConnection(processInfo)
	takeOff()
}

func initConfig() {
	viper.SetConfigName("w2psd")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/weewoe")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/weewoe")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Fatal error config file: %s \n", err.Error())
	}

	log.Infof("read the config file: %s", viper.ConfigFileUsed())

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Debugf("Config file changed:%s, op: %d", e.Name, e.Op)
		// log.Debugf("ALL KEY:%+v", viper.AllKeys())
		if len(viper.AllKeys()) == 0 {
			return
		}
		pmux.Lock()
		defer pmux.Unlock()
		loadProcessInfo()
		initConnection(processInfo)
	})
	viper.WatchConfig()

	loadProcessInfo()
}

func loadProcessInfo() {
	pro = viper.Sub("data")

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
			conn, err := ssh.NewConnection(v.Host+":22", "root")
			if err != nil {
				log.Errorf("connect to %s failed: %s", v.Host, err.Error())
				continue
			}
			instances[v.Host] = Target{
				Conn: conn,
			}
		}
	}
}
