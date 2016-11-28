package watch_xixun

import (
	"github.com/giskook/gotcp"
	//"github.com/huoyan108/watch_xixun/nsq"
	//"github.com/huoyan108/watch_xixun"
	"log"
	"net"
	"time"
)

type ServerConfig struct {
	Listener      *net.TCPListener
	AcceptTimeout time.Duration
}

type Server struct {
	config           *ServerConfig
	srv              *gotcp.Server
	checkconnsticker *time.Ticker
	//	Dbsrv            *DbServer
	nsqpserverManager   *NsqProducer
	nsqpserverLoction   *NsqProducer
	nsqpserverControlcb *NsqProducer
	nsqcserverManager   *NsqConsumer
	//nsqcserverControl   *[]NsqConsumer
	nsqConsumerConfig *NsqConsumerConfig
	nsqConsumers      *NsqConsumers
}

var Gserver *Server

func SetServer(server *Server) {
	Gserver = server
}

func GetServer() *Server {
	return Gserver
}
func NewServer(srv *gotcp.Server, config *ServerConfig, nsqpserverManager *NsqProducer, nsqpserverLoction *NsqProducer, nsqpserverControlcb *NsqProducer, nsqcserverManager *NsqConsumer,
	nsqConsumerConfig *NsqConsumerConfig) *Server {
	serverstatistics := GetConfiguration().GetServerStatistics()
	return &Server{
		config:           config,
		srv:              srv,
		checkconnsticker: time.NewTicker(time.Duration(serverstatistics) * time.Second),
		//	Dbsrv:            dbsrv,
		nsqpserverManager:   nsqpserverManager,
		nsqpserverLoction:   nsqpserverLoction,
		nsqpserverControlcb: nsqpserverControlcb,
		nsqcserverManager:   nsqcserverManager,
		nsqConsumerConfig:   nsqConsumerConfig,
		nsqConsumers:        NewNsqConsumers(nsqConsumerConfig),
	}
}
func (s *Server) GetNsqConsumers_Control() *NsqConsumers {
	return s.nsqConsumers
}
func (s *Server) GetProducerManager() *NsqProducer {
	return s.nsqpserverManager
}
func (s *Server) GetProducerLocation() *NsqProducer {
	return s.nsqpserverLoction
}
func (s *Server) GetProducerControl() *NsqProducer {
	return s.nsqpserverControlcb
}

func (s *Server) GetConsumerManager() *NsqConsumer {
	return s.nsqcserverManager
}

//func (s *Server) GetProducerManagerTopic() string {
//	return s.config.UptopicManager
//}
//func (s *Server) GetProducerLocationTopic() string {
//	return s.config.UptopicLocation
//}
//func (s *Server) GetProducerControlTopic() string {
//	return s.config.UptopicControl
//}
func (s *Server) Start() {
	go s.checkStatistics()
	s.nsqpserverManager.Start()
	s.nsqpserverLoction.Start()
	s.nsqpserverControlcb.Start()
	s.nsqcserverManager.Start()
	s.srv.Start(s.config.Listener, s.config.AcceptTimeout)
	//	s.nsqcserverManager.Start()
	//	s.nsqcserverManager.Stop()
	//
}

func (s *Server) Stop() {
	s.srv.Stop()
	s.nsqpserverManager.Stop()
	s.nsqpserverLoction.Stop()
	s.nsqpserverControlcb.Stop()
	s.nsqcserverManager.Stop()
	s.checkconnsticker.Stop()
}

func (s *Server) checkStatistics() {
	for {
		<-s.checkconnsticker.C
		log.Printf("---------------------Totol Connections : %d---------------------\n", NewConns().GetCount())
	}
}
