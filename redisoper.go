package watch_xixun

import (
	//	"fmt"
	"github.com/huoyan108/logs"
	"gopkg.in/redis.v5"
	"strings"
	"time"
)

var redisOper *RedisOper
var versionName string
var downloadLink string

type RedisOper struct {
	Addr       string
	client     *redis.Client
	closeChan  chan bool
	redisMatch string
}

func NewRedisOper(addr string, redisMatch string) *RedisOper {

	defer logs.Logger.Flush()
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	logs.Logger.Info(pong, err)

	redisOper = &RedisOper{
		client:     client,
		closeChan:  make(chan bool),
		redisMatch: redisMatch,
	}
	return redisOper
}
func (c *RedisOper) Set(key string, value string) bool {
	err := c.client.Set(key, value, 0).Err()
	if err != nil {
		logs.Logger.Info("Set redis error", key)
		return false
	}
	return true
}
func (c *RedisOper) SetbyMatch(key string, value string) bool {
	key += c.redisMatch
	err := c.client.Set(key, value, 0).Err()
	if err != nil {
		logs.Logger.Info("Set redis error", key)
		return false
	}
	return true
}
func (c *RedisOper) Get(key string) string {
	value, err := c.client.Get(key).Result()
	if err != nil {
		logs.Logger.Info("get redis ", key, ":warn", err)
		return ""
	}
	return value
}
func (c *RedisOper) GetbyMatch(key string) string {
	value, err := c.client.Get(key + c.redisMatch).Result()
	if err != nil {
		logs.Logger.Info("get redis ", key, ":warn:", err)
		return ""
	}
	return value
}
func (c *RedisOper) Del(key string) bool {
	err := c.client.Del(key).Err()
	if err != nil {
		logs.Logger.Info("del redis ", key, "warn", err)
		return false
	}
	return true
}
func (c *RedisOper) DelbyMatch(key string) bool {
	err := c.client.Del(key + c.redisMatch).Err()
	if err != nil {
		logs.Logger.Info("del redis ", key, ":warn:", err)
		return false
	}
	return true
}
func (c *RedisOper) SendToTerminal(value string) {

}
func (c *RedisOper) ScanSend() {

	iter := c.client.Scan(0, "*"+c.redisMatch, 100).Iterator()
	delArr := []string{}
	for iter.Next() {
		//	fmt.Println(iter.Val())
		if len(iter.Val()) != 15 {
			continue
		}
		value := c.Get(iter.Val())
		if value == "" {
			continue
		}
		values := strings.Split(value, ",")
		timeReq := values[0]
		tid := values[1]
		serialnum := values[2]
		action := values[3]

		timereq, _ := time.Parse("060102150405", timeReq)
		dis := time.Now().Sub(timereq).Seconds()
		if int(dis) > 10 {
			logs.Logger.Info("Second send Command ", dis, action)
			SendToClientControl(timeReq, tid, serialnum, action, false)
			delArr = append(delArr, iter.Val())
		}
	}
	for _, key := range delArr {
		c.Del(key)
	}

	if err := iter.Err(); err != nil {
		panic(err)
	}
}
func (c *RedisOper) GetNewVersion() {

	versionName = redisOper.Get("versionName")
	downloadLink = redisOper.Get("downloadLink")
}
func (c *RedisOper) Start() {
	defer func() {
		//	c.Client.Close()
	}()

	for {
		select {
		case <-c.closeChan:
			return
		//case <-time.After(time.Duration(5) * time.Second):
		//	c.ScanSend()
		case <-time.After(time.Duration(60) * time.Second):
			c.GetNewVersion()
		}
	}

}
func (c *RedisOper) Stop() {
	c.closeChan <- true
}
