package watch_xixun

import (
	"log"
	"sync"

	"github.com/bitly/go-nsq"
	//	"github.com/huoyan108/watch_xixun"
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
			log.Println("err found")
			s.Stop()
		}

	}()

	config := nsq.NewConfig()

	var errmsg error
	s.consumer, errmsg = nsq.NewConsumer(s.config.Topic, s.config.Channel, config)

	if errmsg != nil {
		//	panic("create consumer error -> " + errmsg.Error())
		log.Println("create consumer error -> " + errmsg.Error())
	}
	s.recvNsq()

	err := s.consumer.ConnectToNSQD(s.config.Addr)
	if err != nil {
		panic("Counld not connect to nsq -> " + err.Error())
	}
}

func (s *NsqConsumer) Stop() {
	log.Println("stop consumer")
	s.waitGroup.Wait()
	log.Println("stop consumer2", s.config.Addr)

	errmsg := s.consumer.DisconnectFromNSQD(s.config.Addr)

	if errmsg != nil {
		log.Printf("stop consumer error ", errmsg.Error())
	}

	//	s.consumer.Stop()
	log.Println("stop consumer3")
}
