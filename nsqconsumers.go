package watch_xixun

import (
	"log"
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
	newConfig := &NsqConsumerConfig{
		Topic:   topic,
		Addr:    s.nsqConsumerConfig.Addr,
		Channel: s.nsqConsumerConfig.Channel,
	}
	server := s.serversMap[topic]
	if server != nil {
		log.Println("MakeNewConsumer has", topic)
		//return false
	}

	server = NewNsqConsumer(newConfig, Consumer_Control)
	if server == nil {
		return false
	}
	server.Start()
	s.serversMap[topic] = server
	return true
}
func (s *NsqConsumers) DelConsumer(topic string) {
	log.Println("delconsumer")
	server := s.serversMap[topic]
	if server != nil {
		server.Stop()
		log.Println("delconsumer1")
		delete(s.serversMap, topic)
		log.Println("delconsumer2")
	}
}
