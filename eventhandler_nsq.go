package watch_xixun

import (
	"github.com/golang/protobuf/proto"
	"github.com/huoyan108/watch_xixun/pbgo"
	"github.com/huoyan108/watch_xixun/protocol"
	"log"
	"strconv"
	"time"
)

func Nsq_EventHandler(nsqType int32, sTopic string, data []byte) {

	if nsqType == Consumer_Manager {

		//登陆结果
		command := &Report.ManageProtocol{}
		err := proto.Unmarshal(data, command)
		if err != nil {
			log.Println("unmarshal manager error")
		} else {
			log.Println("<IN_NSQ> topic:", sTopic, command)
			tid := command.Tid
			serialnum := command.SerialNumber
			res := command.Result

			if res == Report.ManageProtocol_SUCCESS {
				//向终端发送登陆回执
				ntid, _ := strconv.Atoi(tid)
				c := NewConns().GetConn(uint64(ntid))
				if c != nil {
					loginrt := protocol.SetLoginRt(c.Encryption, tid, serialnum)
					c.SendToClient(loginrt)
					//更改连接状态
					c.Status = ConnSuccess
					//创建独立的控制指令订阅服务
					if GetServer().GetNsqConsumers_Control().MakeNewConsumer(tid) == false {
						log.Println("make new nsqconsumer error.Topic", tid)
					}
				} else {

					log.Println("cann't find tid:", tid)
				}

			} else {
			}
		}

	} else if nsqType == Consumer_Control {
		//控制指令
		command := &Report.ControlProtocol{}
		err := proto.Unmarshal(data, command)
		if err != nil {
			log.Println("unmarshal control error")
		} else {
			log.Println("<IN_NSQ> topic:", sTopic, command)
			timeReq := command.TimeReq
			timereq, _ := time.Parse("060102150405", timeReq)
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
				dis := time.Now().Sub(timereq).Minutes()
				if int(dis) > 3 {
					log.Println("RT_OFF command timeout", dis)
					//return
				}
				action := ""
				action += command.Command.Paras[0].Strpara

				ntid, _ := strconv.Atoi(tid)
				c := NewConns().GetConn(uint64(ntid))
				if c != nil {
					rt := protocol.SetControlCmd(c.Encryption, tid, serialnum, "", action)
					c.SendToClient(rt)
				}
			case Report.Command_CMT_POSP:
				dis := time.Now().Sub(timereq).Minutes()
				if int(dis) > 3 {
					log.Println("POSP command timeout", dis)
					return
				}
				action := "tk="
				action += command.Command.Paras[0].Strpara
				ntid, _ := strconv.Atoi(tid)
				c := NewConns().GetConn(uint64(ntid))
				if c != nil {
					rt := protocol.SetControlCmd(c.Encryption, tid, serialnum, "", action)
					c.SendToClient(rt)
				}
			case Report.Command_CMT_POS_INTERVAL:
				dis := time.Now().Sub(timereq).Minutes()
				if int(dis) > 3 {
					log.Println("INTERVAL command timeout", dis)
					return
				}
				action := "ti="
				action += command.Command.Paras[0].Strpara
				ntid, _ := strconv.Atoi(tid)
				c := NewConns().GetConn(uint64(ntid))
				if c != nil {
					rt := protocol.SetControlCmd(c.Encryption, tid, serialnum, "", action)
					c.SendToClient(rt)
				}
			case Report.Command_CMT_SERVER_SET:
				dis := time.Now().Sub(timereq).Minutes()
				if int(dis) > 3 {
					log.Println("SERVER SET command timeout", dis)
					return
				}
				action := "ip="
				action += command.Command.Paras[0].Strpara
				ntid, _ := strconv.Atoi(tid)
				c := NewConns().GetConn(uint64(ntid))
				if c != nil {
					rt := protocol.SetControlCmd(c.Encryption, tid, serialnum, "", action)
					c.SendToClient(rt)
				}
			case Report.Command_CMT_UPDATE:
				dis := time.Now().Sub(timereq).Minutes()
				if int(dis) > 3 {
					log.Println("UPDATE command timeout", dis)
					//return
				}
				action := "up="
				action += command.Command.Paras[0].Strpara
				ntid, _ := strconv.Atoi(tid)
				c := NewConns().GetConn(uint64(ntid))
				if c != nil {
					rt := protocol.SetControlCmd(c.Encryption, tid, serialnum, "", action)
					c.SendToClient(rt)
				}
			case Report.Command_CMT_SENDMSG:
				dis := time.Now().Sub(timereq).Minutes()
				if int(dis) > 3 {
					log.Println("SEND_MSG command timeout", dis)
					//return
				}
				action := "mg="
				action += command.Command.Paras[0].Strpara
				ntid, _ := strconv.Atoi(tid)
				c := NewConns().GetConn(uint64(ntid))
				if c != nil {
					rt := protocol.SetControlCmd(c.Encryption, tid, serialnum, "", action)
					c.SendToClient(rt)
				}
			case Report.Command_CMT_INIT:
				dis := time.Now().Sub(timereq).Minutes()
				if int(dis) > 3 {
					log.Println("INIT command timeout", dis)
					return
				}
				action := "th"
				ntid, _ := strconv.Atoi(tid)
				c := NewConns().GetConn(uint64(ntid))
				if c != nil {
					rt := protocol.SetControlCmd(c.Encryption, tid, serialnum, "", action)
					c.SendToClient(rt)
				}
			case Report.Command_CMT_SET_IMEI:
				dis := time.Now().Sub(timereq).Minutes()
				if int(dis) > 3 {
					log.Println("SET_IMEI command timeout", dis)
					return
				}
				action := "si="
				action += command.Command.Paras[0].Strpara
				ntid, _ := strconv.Atoi(tid)
				c := NewConns().GetConn(uint64(ntid))
				if c != nil {
					rt := protocol.SetControlCmd(c.Encryption, tid, serialnum, "", action)
					c.SendToClient(rt)
				}

			}
		}

	} else {
	}
}
