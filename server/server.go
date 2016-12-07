package main

import (
	"fmt"
	"github.com/huoyan108/gotcp"
	"github.com/huoyan108/logs"
	"github.com/huoyan108/watch_xixun"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//init log
	logs.SetLogConf("./seelog.xml")
	defer logs.Logger.Flush()

	// read configuration
	configuration, err := watch_xixun.ReadConfig("./conf.json")
	watch_xixun.SetConfiguration(configuration)
	checkError(err)

	//init redis client
	redisoper := watch_xixun.NewRedisOper(configuration.RedisConfig.ServerInfo, configuration.RedisConfig.RedisMatch)
	go redisoper.Start()

	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+configuration.ServerConfig.BindPort)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a tcp server
	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &watch_xixun.Callback{}, &watch_xixun.ShaProtocol{})

	//create nsq client
	// creates a nsqproducer server DAS TO DCS
	//manager
	nsqpconfigManager := &watch_xixun.NsqProducerConfig{
		Addr:  configuration.NsqConfig.Addr,
		Topic: configuration.NsqConfig.UpTopicManager,
	}
	logs.Logger.Info(configuration.NsqConfig)
	nsqpserverManager := watch_xixun.NewNsqProducer(nsqpconfigManager)
	//Loction
	nsqpconfigLoctionLoction := &watch_xixun.NsqProducerConfig{
		Addr:  configuration.NsqConfig.Addr,
		Topic: configuration.NsqConfig.UpTopicLoction,
	}
	nsqpserverLoction := watch_xixun.NewNsqProducer(nsqpconfigLoctionLoction)

	//controlcb
	nsqpconfigControlcb := &watch_xixun.NsqProducerConfig{
		Addr:  configuration.NsqConfig.Addr,
		Topic: configuration.NsqConfig.UpTopicControl,
	}
	nsqpserverControlcb := watch_xixun.NewNsqProducer(nsqpconfigControlcb)

	// creates a nsqconsumer server   DCS TO DAS
	//downManagercb
	nsqcconfigManager := &watch_xixun.NsqConsumerConfig{
		Addr:    configuration.NsqConfig.Addr,
		Topic:   configuration.NsqConfig.DownTopicManager,
		Channel: configuration.NsqConfig.DownchannelManager,
	}
	nsqcserverManager := watch_xixun.NewNsqConsumer(nsqcconfigManager, watch_xixun.Consumer_Manager)

	// create watch_xixun server
	watch_xixunserverconfig := &watch_xixun.ServerConfig{
		Listener:      listener,
		AcceptTimeout: time.Duration(configuration.ServerConfig.ConnTimeout) * time.Second,
	}
	watch_xixunserver := watch_xixun.NewServer(srv, watch_xixunserverconfig, nsqpserverManager,
		nsqpserverLoction, nsqpserverControlcb, nsqcserverManager, nsqcconfigManager)
	watch_xixun.SetServer(watch_xixunserver)
	// starts service
	fmt.Println("listening:", listener.Addr())
	watch_xixunserver.Start()

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	redisoper.Stop()
	watch_xixunserver.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
