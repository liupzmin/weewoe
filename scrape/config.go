package scrape

import (
	"strings"
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

var (
	instances = &Instances{
		Set: make(map[string]Target),
	}
	processInfo Config
	pmux        sync.RWMutex
	pro         *viper.Viper
)

type Instances struct {
	sync.RWMutex
	Set map[string]Target
}

func (i *Instances) GetTarget(key string) (*Target, bool) {
	i.RLock()
	defer i.RUnlock()
	var (
		t  Target
		ok bool
	)
	if t, ok = i.Set[key]; ok {
		return &t, ok
	}
	return nil, ok
}

func (i *Instances) AddConn(key string, conn *ssh.Connection) {
	i.Lock()
	defer i.Unlock()
	i.Set[key] = Target{
		Conn: conn,
	}
}

func (i *Instances) Clear() {
	i.Lock()
	defer i.Unlock()
	for _, v := range i.Set {
		v.Close()
	}
	i.Set = make(map[string]Target)
}

type Target struct {
	Conn *ssh.Connection
	// used to compute process start time
	BootTime int64
}

func (t Target) Close() {
	t.Conn.Close()
}

type Process struct {
	OSUser  string
	Name    string
	Path    string
	Ports   []int64
	Group   string
	Host    string
	PIDFile string
	Flag    string
	Suspend bool
}

func (p *Process) GetConnectionKey(host string) string {
	return host + p.OSUser + p.Name
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

	var pi Config
	err := pro.Unmarshal(&pi)
	if err != nil {
		log.Panicf("unable to decode into struct, %v", err)
	}

	processInfo = pi
	log.Debugf("the process info has been loaded:%+v", processInfo)
}

func initConnection(conf Config) {
	log.Info("begin to create ssh connection!")
	defer func() {
		log.Info("create ssh connection done!")
	}()

	instances.Clear()
	for _, v := range conf.Processes {
		h := strings.Split(v.Host, ":")[0]
		for _, p := range v.Process {
			if _, ok := instances.GetTarget(p.GetConnectionKey(h)); ok {
				continue
			}

			conn, err := ssh.NewConnection(v.Host, p.OSUser)
			if err != nil {
				log.Errorf("connect to %s failed: %s", v.Host, err.Error())
				continue
			}
			instances.AddConn(p.GetConnectionKey(h), conn)
		}
	}
}
