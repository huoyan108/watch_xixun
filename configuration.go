package watch_xixun

import (
	"encoding/json"
	"os"
)

type NsqConfiguration struct {
	Addr               string
	UpTopicManager     string
	UpTopicLoction     string
	UpTopicControl     string
	DownTopicManager   string
	DownchannelManager string
}

type ServerConfiguration struct {
	BindPort          string
	ReadLimit         uint32
	WriteLimit        uint16
	ConnTimeout       uint16
	ConnCheckInterval uint16
	ServerStatistics  uint16
}

type RedisConfiguration struct {
	ServerInfo string
	RedisMatch string
}

type Configuration struct {
	NsqConfig    *NsqConfiguration
	ServerConfig *ServerConfiguration
	RedisConfig  *RedisConfiguration
}

func ReadConfig(confpath string) (*Configuration, error) {
	file, _ := os.Open(confpath)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	return &configuration, err
}

func (conf *Configuration) GetServerReadLimit() uint32 {
	return conf.ServerConfig.ReadLimit
}

func (conf *Configuration) GetServerWriteLimit() uint16 {
	return conf.ServerConfig.WriteLimit
}

func (conf *Configuration) GetServerConnCheckInterval() uint16 {
	return conf.ServerConfig.ConnCheckInterval
}

func (conf *Configuration) GetServerStatistics() uint16 {
	return conf.ServerConfig.ServerStatistics
}

var Config *Configuration

func SetConfiguration(config *Configuration) {
	Config = config
}

func GetConfiguration() *Configuration {
	return Config
}
