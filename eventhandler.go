package watch_xixun

import (
	"github.com/golang/protobuf/proto"
	"github.com/huoyan108/gotcp"
	"github.com/huoyan108/logs"
	"github.com/huoyan108/watch_xixun/pbgo"
	"github.com/huoyan108/watch_xixun/protocol"
	"strconv"
	"time"
)

type Callback struct{}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	checkinterval := GetConfiguration().GetServerConnCheckInterval()
	readlimit := GetConfiguration().GetServerReadLimit()
	writelimit := GetConfiguration().GetServerWriteLimit()
	config := &ConnConfig{
		ConnCheckInterval: uint16(checkinterval),
		ReadLimit:         uint16(readlimit),
		WriteLimit:        uint16(writelimit),
	}
	conn := NewConn(c, config)

	c.PutExtraData(conn)

	conn.Do()
	NewConns().Add(conn)

	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	//defer logs.Logger.Flush()
	conn := c.GetExtraData().(*Conn)
	//logs.Logger.Info("Close client all info")
	GetServer().GetNsqConsumers_Control().DelConsumer(conn.IMEI)
	NewConns().Remove(conn)
	conn.Close()
}

func on_login(c *gotcp.Conn, p *ShaPacket) {
	defer logs.Logger.Flush()
	conn := c.GetExtraData().(*Conn)
	loginPkg := p.Packet.(*protocol.LoginPacket)
	conn.IMEI = loginPkg.IMEI
	conn.ID, _ = strconv.ParseUint(loginPkg.IMEI, 10, 64)
	conn.Encryption = loginPkg.Encryption

	NewConns().SetID(conn.ID, conn.index)

	//构造指令内容
	t := time.Now()
	nowtime := t.Format("060102150405")

	req := &Report.ManageProtocol{
		TimeReq:      nowtime,
		SerialNumber: loginPkg.SerialNumber,
		Tid:          loginPkg.IMEI,
		Type:         Report.ManageProtocol_LOGIN,
		TerminalType: "xixun",
		ProtocolType: loginPkg.ProtocolType,
	}

	//判断是否更新
	logs.Logger.Info("login  ", versionName, loginPkg.ProtocolType)
	if versionName != "" && loginPkg.ProtocolType != "" && versionName != loginPkg.ProtocolType {
		if downloadLink != "" {
			//登陆回执
			c.AsyncWritePacket(p, time.Second)
			//更新命令
			time.AfterFunc(2*time.Second, func() {
				pkg := protocol.SetControlCmd(loginPkg.Encryption, loginPkg.IMEI, loginPkg.SerialNumber, "", "up="+downloadLink)
				c.AsyncWritePacket(pkg, time.Second)
			})
		}
		return
	}
	logs.Logger.Info("make login proto", req)
	reqdata, _ := proto.Marshal(req)
	GetServer().GetProducerManager().Send(GetConfiguration().NsqConfig.UpTopicManager, reqdata)
	logs.Logger.Info("Send login to Dcs", GetConfiguration().NsqConfig.UpTopicManager)
}

func on_posup(c *gotcp.Conn, p *ShaPacket) {
	defer logs.Logger.Flush()
	posup_pkg := p.Packet.(*protocol.PosUpPacket)

	fBattery, _ := strconv.ParseFloat(posup_pkg.Battery, 32)
	nBattery := int32(fBattery * 100)
	t := time.Now()
	nowtime := t.Format("060102150405")
	req := &Report.LocationProtocol{
		TimeReq:      nowtime,
		SerialNumber: posup_pkg.SerialNumber,
		Tid:          posup_pkg.IMEI,
		Locations:    []*Report.Location{},
		Mld: &Report.MobileLocationData{

			//Battery: int32(posup_pkg.Battery), //电量
			Battery: nBattery,                //电量
			Charge:  int32(posup_pkg.Charge), //充电
			//Reason:  1,                 //0定位,1频度上报
		},
	}

	if posup_pkg.WifiCount > 0 {
		c.AsyncWritePacket(p, time.Second)
		location :=
			Report.Location{
				Locationtype: Report.Location_EWifi,
				From:         1,
				Wifis: []*Report.Wifi{
					&Report.Wifi{
						Wifixx: posup_pkg.Wifi,
					},
				},
			}
		req.Locations = append(req.Locations, &location)
	} else if posup_pkg.GPSFlag != "" {
		c.AsyncWritePacket(p, time.Second)
		if posup_pkg.Longitude != "" {
			long, _ := strconv.ParseFloat(posup_pkg.Longitude, 32)
			lat, _ := strconv.ParseFloat(posup_pkg.Latitude, 32)
			location :=
				Report.Location{

					Locationtype: Report.Location_EWGS84,
					From:         1,
					Wgs84: &Report.WGS84{
						LngE6:     int32(long * 1000000),
						LatE6:     int32(lat * 1000000),
						Speed:     1,
						Degree:    1,
						Precision: 1,
						Altitude:  1,
					},
				}
			req.Locations = append(req.Locations, &location)
		} else if posup_pkg.WifiCount > 0 {
			location :=
				Report.Location{
					Locationtype: Report.Location_EWifi,
					From:         1,
					Wifis: []*Report.Wifi{
						&Report.Wifi{
							Wifixx: posup_pkg.Wifi,
						},
					},
				}
			req.Locations = append(req.Locations, &location)
		} else if posup_pkg.JzCount > 0 {
			location :=
				Report.Location{
					Locationtype: Report.Location_EMobileCell,
					From:         1,
					Cells: []*Report.MobileCell{
						&Report.MobileCell{
							Jzxx: posup_pkg.Jzs,
						},
					},
				}
			req.Locations = append(req.Locations, &location)

		}
	}
	logs.Logger.Info("make posup proto", req)
	reqdata, _ := proto.Marshal(req)
	GetServer().GetProducerManager().Send(GetConfiguration().NsqConfig.UpTopicLoction, reqdata)
	logs.Logger.Info("Send posup to Dcs", GetConfiguration().NsqConfig.UpTopicLoction)
}

const (
	TK string = "tk"
	TI string = "ti"
	TT string = "tt"
	RT string = "rt"
	OF string = "of"
	IP string = "ip"
	UP string = "up"
	MG string = "mg"
)

func on_ControlCmdrt(c *gotcp.Conn, p *ShaPacket) {
	defer logs.Logger.Flush()
	controlRt_pkg := p.Packet.(*protocol.ControlPacket)
	//构造指令内容
	t := time.Now()
	nowtime := t.Format("060102150405")

	style := Report.Command_CMT_REP
	switch controlRt_pkg.Style {
	case TK:
		style = Report.Command_CMT_POSP
	case TI:
		style = Report.Command_CMT_POS_INTERVAL
	case TT:
		style = Report.Command_CMT_POS_INTERVAL
	case RT:
		style = Report.Command_CMT_RT_OFF
	case OF:
		style = Report.Command_CMT_RT_OFF
	case IP:
		style = Report.Command_CMT_SERVER_SET
	case UP:
		style = Report.Command_CMT_UPDATE
	case MG:
		style = Report.Command_CMT_SENDMSG
	default:
		return
		//style = Report.Command_CMT_REP
	}
	req := &Report.ControlProtocol{
		TimeReq:      nowtime,
		SerialNumber: controlRt_pkg.SerialNumber,
		Tid:          controlRt_pkg.IMEI,
		Command: &Report.Command{
			Type: Report.Command_CMT_REP,
			Paras: []*Report.Command_Param{
				&Report.Command_Param{
					Type:    Report.Command_Param_STRING,
					Strpara: style.String(),
				},
			},
		},
	}
	logs.Logger.Info("make controlrt proto", req)
	reqdata, _ := proto.Marshal(req)
	GetServer().GetProducerManager().Send(GetConfiguration().NsqConfig.UpTopicControl, reqdata)
	logs.Logger.Info("Send controlrt to Dcs", GetConfiguration().NsqConfig.UpTopicControl, controlRt_pkg.Action)
	redisOper.DelbyMatch(controlRt_pkg.IMEI)
}
func on_warnup(c *gotcp.Conn, p *ShaPacket) {
	defer logs.Logger.Flush()
	c.AsyncWritePacket(p, time.Second)

	warnup_pkg := p.Packet.(*protocol.WarnUpPacket)

	warnstyle, _ := strconv.Atoi(warnup_pkg.WarnStyle)
	//构造指令内容
	t := time.Now()
	nowtime := t.Format("060102150405")

	req := &Report.ControlProtocol{
		TimeReq:      nowtime,
		SerialNumber: warnup_pkg.SerialNumber,
		Tid:          warnup_pkg.IMEI,
		Command: &Report.Command{
			Type: Report.Command_CMT_WARNUP,
			Paras: []*Report.Command_Param{
				&Report.Command_Param{
					Type:  Report.Command_Param_INT8,
					Npara: int64(warnstyle),
				},
			},
		},
	}
	logs.Logger.Info("make warnup proto", req)
	reqdata, _ := proto.Marshal(req)
	GetServer().GetProducerManager().Send(GetConfiguration().NsqConfig.UpTopicControl, reqdata)
	logs.Logger.Info("Send warnup to Dcs", GetConfiguration().NsqConfig.UpTopicControl)

}
func on_ReadMsg(c *gotcp.Conn, p *ShaPacket) {
	//c.AsyncWritePacket(p, time.Second)
	defer logs.Logger.Flush()

	readMsg_pkg := p.Packet.(*protocol.ReadMsgPacket)

	msgId := readMsg_pkg.MsgId
	//构造指令内容
	t := time.Now()
	nowtime := t.Format("060102150405")

	req := &Report.ControlProtocol{
		TimeReq:      nowtime,
		SerialNumber: readMsg_pkg.SerialNumber,
		Tid:          readMsg_pkg.IMEI,
		Command: &Report.Command{
			Type: Report.Command_CMT_WARNUP,
			Paras: []*Report.Command_Param{
				&Report.Command_Param{
					Type:    Report.Command_Param_STRING,
					Strpara: msgId,
				},
			},
		},
	}
	logs.Logger.Info("make readmsg proto", req)
	reqdata, _ := proto.Marshal(req)
	GetServer().GetProducerManager().Send(GetConfiguration().NsqConfig.UpTopicControl, reqdata)
	logs.Logger.Info("Send readmsg to Dcs", GetConfiguration().NsqConfig.UpTopicControl)

}
func on_Charge(c *gotcp.Conn, p *ShaPacket) {
	defer logs.Logger.Flush()
	//c.AsyncWritePacket(p, time.Second)

	charge_pkg := p.Packet.(*protocol.ChargePacket)

	//构造指令内容
	t := time.Now()
	nowtime := t.Format("060102150405")

	req := &Report.ControlProtocol{
		TimeReq:      nowtime,
		SerialNumber: charge_pkg.SerialNumber,
		Tid:          charge_pkg.IMEI,
		Command: &Report.Command{
			Type: Report.Command_CMT_CONN_CHARGER,
		},
	}
	logs.Logger.Info("make conn_charge proto", req)
	reqdata, _ := proto.Marshal(req)
	GetServer().GetProducerManager().Send(GetConfiguration().NsqConfig.UpTopicControl, reqdata)
	logs.Logger.Info("Send charge to Dcs", GetConfiguration().NsqConfig.UpTopicControl)

}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	shaPacket := p.(*ShaPacket)
	switch shaPacket.Type {
	case protocol.Login:
		on_login(c, shaPacket)
	case protocol.HeartBeat:
		c.AsyncWritePacket(shaPacket, time.Second)
	case protocol.PosUp:
		on_posup(c, shaPacket)
	case protocol.Echo:
		c.AsyncWritePacket(shaPacket, time.Second)
	case protocol.WarnUp:
		on_warnup(c, shaPacket)
	case protocol.ControlCmdRt:
		on_ControlCmdrt(c, shaPacket)
	case protocol.ReadMsg:
		on_ReadMsg(c, shaPacket)
	case protocol.Charge:
		on_Charge(c, shaPacket)
	}

	return true
}
