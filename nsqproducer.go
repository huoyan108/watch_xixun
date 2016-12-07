package watch_xixun

import (
	"github.com/bitly/go-nsq"
	"github.com/huoyan108/logs"
)

type NsqProducerConfig struct {
	Addr  string
	Topic string
}

type NsqProducer struct {
	config *NsqProducerConfig
	//waitGroup *sync.WaitGroup

	producer *nsq.Producer
}

func NewNsqProducer(config *NsqProducerConfig) *NsqProducer {
	return &NsqProducer{
		config: config,
		//waitGroup: &sync.WaitGroup{},
	}
}

func (s *NsqProducer) Send(topic string, value []byte) error {
	defer logs.Logger.Flush()
	logs.Logger.Info("<OUT_NSQ> topic:", topic)
	err := s.producer.PublishAsync(topic, value, nil, nil)

	if err != nil {
		logs.Logger.Info(topic)
	}
	return err
}

func (s *NsqProducer) Start() {
	//s.waitGroup.Add(1)
	defer func() {
		err := recover()
		//s.waitGroup.Done()
		if err != nil {
			logs.Logger.Info("err found")
			logs.Logger.Flush()
		}

	}()
	config := nsq.NewConfig()

	var errmsg error
	s.producer, errmsg = nsq.NewProducer(s.config.Addr, config)

	if errmsg != nil {
		//	log.Printf("create producer error" + errmsg.Error())
		//panic("create producer error " + errmsg.Error())
		logs.Logger.Info("create producer error" + errmsg.Error())
	}
}

func (s *NsqProducer) GetTopic() string {
	return s.config.Topic
}

func (s *NsqProducer) Stop() {
	//	s.waitGroup.Done()
	//	s.waitGroup.Wait()

	s.producer.Stop()
}
