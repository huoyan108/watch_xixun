package watch_xixun

import (
	"github.com/huoyan108/logs"
)

type NsqConsumers struct {
	serversMap        map[string]*NsqConsumer
	nsqConsumerConfig *NsqConsumerConfig
}

var GNsqConsumers *NsqConsumers

func GetConsumers() *NsqConsumers {
	return GNsqConsumers
}
func NewNsqConsumers(nsqConsumerConfig *NsqConsumerConfig) *NsqConsumers {
	GNsqConsumers = &NsqConsumers{
		nsqConsumerConfig: nsqConsumerConfig,
		serversMap:        make(map[string]*NsqConsumer),
	}
	return GNsqConsumers
}
func (s *NsqConsumers) MakeNewConsumer(topic string) bool {
	defer logs.Logger.Flush()
	newConfig := &NsqConsumerConfig{
		Topic:   topic,
		Addr:    s.nsqConsumerConfig.Addr,
		Channel: s.nsqConsumerConfig.Channel,
	}
	server := s.serversMap[topic]
	if server != nil {
		logs.Logger.Info("has ", topic, " cann't make new consumer")
		return false
	}

	server = NewNsqConsumer(newConfig, Consumer_Control)
	if server == nil {
		logs.Logger.Info("make new consumer error,topic:", topic)
		return false
	}
	logs.Logger.Info("make new consumer sucess,topic:", topic)
	server.Start()
	s.serversMap[topic] = server
	logs.Logger.Info("make new consumer sucess,topic:", topic)
	return true
}
func (s *NsqConsumers) DelConsumer(topic string) {
	server := s.serversMap[topic]
	if server != nil {
		server.Stop()
		delete(s.serversMap, topic)
		logs.Logger.Info("stop consumer commplete,topic:", topic)
	}
}
