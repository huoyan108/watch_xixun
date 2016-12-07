package watch_xixun

import (
	"github.com/golang/protobuf/proto"
	"github.com/huoyan108/logs"
	"github.com/huoyan108/watch_xixun/pbgo"
	"github.com/huoyan108/watch_xixun/protocol"
	"strconv"
	"strings"
	"time"
)

func Nsq_EventHandler(nsqType int32, sTopic string, data []byte) {
	defer logs.Logger.Flush()

	if nsqType == Consumer_Manager {

		//登陆结果
		command := &Report.ManageProtocol{}
		err := proto.Unmarshal(data, command)
		if err != nil {
			logs.Logger.Info("unmarshal manager error")
		} else {
			logs.Logger.Info("<IN_NSQ> topic:", sTopic, command)
			tid := command.Tid
			serialnum := command.SerialNumber
			res := command.Result

			if res == Report.ManageProtocol_SUCCESS {
				ntid, _ := strconv.Atoi(tid)
				c := NewConns().GetConn(uint64(ntid))
				if c != nil {
					//向终端发送登陆回执
					loginrt := protocol.SetLoginRt(c.Encryption, tid, serialnum)
					c.SendToClient(loginrt)
					//更改连接状态
					c.Status = ConnSuccess
					//创建独立的控制指令订阅服务
					if GetServer().GetNsqConsumers_Control().MakeNewConsumer(tid) == false {
						logs.Logger.Info("make new nsqconsumer error.Topic", tid)
					}
					////判断是否更新
					//if versionName != "" && protocolType != "" && versionName != protocolType {
					//	if downloadLink != "" {
					//		time.AfterFunc(2*time.Second, func() {
					//			pkg := protocol.SetControlCmd(c.Encryption, tid, serialnum, "", "up="+downloadLink)
					//			c.SendToClient(pkg)
					//		})
					//	}
					//}
					//发送redis缓存的数据
					value := redisOper.GetbyMatch(tid)
					if value == "" {
						return
					}
					values := strings.Split(value, ",")
					timeReq := values[0]
					tid := values[1]
					serialnum := values[2]
					action := values[3]

					logs.Logger.Info("Second send Command ", action)
					SendToClientControl(timeReq, tid, serialnum, action, false)
					redisOper.DelbyMatch(tid)

				} else {

					logs.Logger.Info("cann't find tid:", tid)
				}

			} else {
			}
		}

	} else if nsqType == Consumer_Control {
		command := &Report.ControlProtocol{}
		err := proto.Unmarshal(data, command)
		if err != nil {
			logs.Logger.Info("unmarshal control error")
		} else {
			logs.Logger.Info("<IN_NSQ> topic:", sTopic, command)
			timeReq := command.TimeReq
			tid := command.Tid
			serialnum := command.SerialNumber
			style := command.Command.Type
			switch style {
			case Report.Command_CMT_REP:
				commandStyle := command.Command.Paras[0].Strpara
				//sCommandStyletyle := strconv.Itoa(int(commandStyle))
				if commandStyle == "Command_CMT_READMSG" {
					warnstyle := command.Command.Paras[0].Npara
					swarnstyle := strconv.Itoa(int(warnstyle))
					ntid, _ := strconv.Atoi(tid)
					c := NewConns().GetConn(uint64(ntid))
					if c != nil {
						rt := protocol.SetWarnUpRt(c.Encryption, tid, serialnum, swarnstyle)
						c.SendToClient(rt)
					}

				} else if commandStyle == "Command_CMT_READMSG " {
					msgId := command.Command.Paras[0].Strpara
					ntid, _ := strconv.Atoi(tid)
					c := NewConns().GetConn(uint64(ntid))
					if c != nil {
						rt := protocol.SetReadMsg(c.Encryption, tid, serialnum, msgId)
						c.SendToClient(rt)
					}

				} else if commandStyle == "Command_CMT_CONN_CHARGER " {
					ntid, _ := strconv.Atoi(tid)
					c := NewConns().GetConn(uint64(ntid))
					if c != nil {
						rt := protocol.SetCharge(c.Encryption, tid, serialnum)
						c.SendToClient(rt)
					}
				}
			case Report.Command_CMT_RT_OFF:
				action := ""
				action += command.Command.Paras[0].Strpara
				SendToClientControl(timeReq, tid, serialnum, action, true)
			case Report.Command_CMT_POSP:
				action := "tk="
				action += command.Command.Paras[0].Strpara
				SendToClientControl(timeReq, tid, serialnum, action, true)
			case Report.Command_CMT_POS_INTERVAL:
				action := "ti="
				action += command.Command.Paras[0].Strpara
				SendToClientControl(timeReq, tid, serialnum, action, true)
			case Report.Command_CMT_SERVER_SET:
				action := "ip="
				action += command.Command.Paras[0].Strpara
				SendToClientControl(timeReq, tid, serialnum, action, true)
			case Report.Command_CMT_UPDATE:
				action := "up="
				action += command.Command.Paras[0].Strpara
				SendToClientControl(timeReq, tid, serialnum, action, true)
			case Report.Command_CMT_SENDMSG:
				action := "mg="
				action += command.Command.Paras[0].Strpara
				SendToClientControl(timeReq, tid, serialnum, action, true)
			case Report.Command_CMT_INIT:
				action := "th"
				SendToClientControl(timeReq, tid, serialnum, action, true)
			case Report.Command_CMT_SET_IMEI:
				action := "si="
				action += command.Command.Paras[0].Strpara
				SendToClientControl(timeReq, tid, serialnum, action, true)
			}
		}

	} else {
	}
}
func SendToClientControl(timeReq string, tid string, serialnum string, action string, bAutoSend bool) bool {
	timereq, _ := time.Parse("060102150405", timeReq)
	dis := time.Now().Sub(timereq).Minutes()
	if int(dis) > 30 {
		logs.Logger.Info("Command timeout", dis)
		return false
	}

	ntid, _ := strconv.Atoi(tid)
	c := NewConns().GetConn(uint64(ntid))
	if c != nil {
		rt := protocol.SetControlCmd(c.Encryption, tid, serialnum, "", action)
		err := c.SendToClient(rt)
		if err != nil {
			return true
		} else {
			if bAutoSend != false {
				logs.Logger.Info("Send client error,", err, "set redis", tid)
				value := timeReq + "," + tid + "," + serialnum + "," + action
				redisOper.SetbyMatch(tid, value)
			}
		}
	} else {
		if bAutoSend != false {
			logs.Logger.Info("Send client error,can't find client,set redis", tid)
			value := timeReq + "," + tid + "," + serialnum + "," + action
			redisOper.SetbyMatch(tid, value)
		}
	}
	return false
}
