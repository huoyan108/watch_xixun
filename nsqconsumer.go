package watch_xixun

import (
	"github.com/bitly/go-nsq"
	"github.com/huoyan108/logs"
	"log"
	"sync"
)

const (
	Consumer_Manager int32 = 1
	Consumer_Control int32 = 2
)

type NsqConsumerConfig struct {
	Addr    string
	Topic   string
	Channel string
}

type NsqConsumer struct {
	config    *NsqConsumerConfig
	waitGroup *sync.WaitGroup

	consumer *nsq.Consumer
	nsqType  int32
}

func NewNsqConsumer(config *NsqConsumerConfig, nsqType int32) *NsqConsumer {
	return &NsqConsumer{
		config:    config,
		waitGroup: &sync.WaitGroup{},
		nsqType:   nsqType,
	}
}

func (s *NsqConsumer) recvNsq() {
	s.consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		data := message.Body
		Nsq_EventHandler(s.nsqType, s.config.Topic, data)

		return nil
	}))
}

func (s *NsqConsumer) Start() {
	s.waitGroup.Add(1)
	defer func() {
		s.waitGroup.Done()
		err := recover()
		if err != nil {
			logs.Logger.Info("err found")
			s.Stop()
		}

	}()
	defer logs.Logger.Flush()

	config := nsq.NewConfig()

	var errmsg error
	s.consumer, errmsg = nsq.NewConsumer(s.config.Topic, s.config.Channel, config)

	if errmsg != nil {
		logs.Logger.Info("create consumer error -> " + errmsg.Error())
	}
	s.recvNsq()

	err := s.consumer.ConnectToNSQD(s.config.Addr)
	if err != nil {
		panic("Counld not connect to nsq -> " + err.Error())
	}
}

func (s *NsqConsumer) Stop() {
	defer logs.Logger.Flush()
	s.waitGroup.Wait()
	logs.Logger.Info("stop consumer,topic:", s.config.Topic)
	errmsg := s.consumer.DisconnectFromNSQD(s.config.Addr)

	if errmsg != nil {
		log.Printf("stop consumer error ", errmsg.Error())
	}

	//	s.consumer.Stop()
}
