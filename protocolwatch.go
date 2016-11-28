package watch_xixun

import (
	"github.com/giskook/gotcp"
	"github.com/huoyan108/watch_xixun/protocol"
	"log"
)

type ShaPacket struct {
	Type   uint16
	Packet gotcp.Packet
}

func (this *ShaPacket) Serialize() []byte {
	switch this.Type {
	case protocol.Login:
		return this.Packet.(*protocol.LoginPacket).Serialize()
	case protocol.HeartBeat:
		return this.Packet.(*protocol.HeartPacket).Serialize()
	case protocol.PosUp:
		return this.Packet.(*protocol.PosUpPacket).Serialize()
	case protocol.Echo:
		return this.Packet.(*protocol.EchoPacket).Serialize()
	case protocol.WarnUp:
		return this.Packet.(*protocol.WarnUpPacket).Serialize()
	case protocol.Charge:
		return this.Packet.(*protocol.ChargePacket).Serialize()
	case protocol.ReadMsg:
		return this.Packet.(*protocol.ReadMsgPacket).Serialize()
	case protocol.ControlCmdRt:
		return this.Packet.(*protocol.ControlPacket).Serialize()

	}

	return nil
}

func NewShaPacket(Type uint16, Packet gotcp.Packet) *ShaPacket {
	return &ShaPacket{
		Type:   Type,
		Packet: Packet,
	}
}

type ShaProtocol struct {
}

func (this *ShaProtocol) ReadPacket(c *gotcp.Conn) (gotcp.Packet, error) {
	smconn := c.GetExtraData().(*Conn)
	smconn.UpdateReadflag()

	buffer := smconn.GetBuffer()

	conn := c.GetRawConn()
	for {
		if smconn.ReadMore {
			data := make([]byte, 2048)
			readLengh, err := conn.Read(data)
			log.Printf("<IN>    %s\n", data[0:readLengh])
			if err != nil {
				return nil, nil
				return nil, err
			}

			if readLengh == 0 {
				log.Println("ReadPacket readLengh==0")
				return nil, gotcp.ErrConnClosing
			}
			buffer.Write(data[0:readLengh])
			smconn.UpdateReadflag()
		}

		cmdid, pkglen := protocol.CheckProtocol(buffer)

		pkgbyte := make([]byte, pkglen)
		buffer.Read(pkgbyte)
		switch cmdid {
		case protocol.Login:
			pkg := protocol.ParseLogin(pkgbyte)
			smconn.ReadMore = false
			return NewShaPacket(protocol.Login, pkg), nil
		case protocol.HeartBeat:
			pkg := protocol.ParseHeart(pkgbyte)
			smconn.ReadMore = false
			return NewShaPacket(protocol.HeartBeat, pkg), nil
		case protocol.PosUp:
			pkg := protocol.ParsePosUp(pkgbyte)
			smconn.ReadMore = false
			return NewShaPacket(protocol.PosUp, pkg), nil
		case protocol.Echo:
			pkg := protocol.ParseEcho(pkgbyte)
			smconn.ReadMore = false
			return NewShaPacket(protocol.Echo, pkg), nil
		case protocol.WarnUp:
			pkg := protocol.ParseWarnUp(pkgbyte)
			smconn.ReadMore = false
			return NewShaPacket(protocol.WarnUp, pkg), nil
		case protocol.Charge:
			pkg := protocol.ParseCharge(pkgbyte)
			smconn.ReadMore = false
			return NewShaPacket(protocol.Charge, pkg), nil
		case protocol.ReadMsg:
			pkg := protocol.ParseReadMsg(pkgbyte)
			smconn.ReadMore = false
			return NewShaPacket(protocol.ReadMsg, pkg), nil
		case protocol.ControlCmdRt:
			pkg := protocol.ParseControlCmdRt(pkgbyte)
			smconn.ReadMore = false
			return NewShaPacket(protocol.ControlCmdRt, pkg), nil
		case protocol.Illegal:
			smconn.ReadMore = true
		case protocol.HalfPack:
			smconn.ReadMore = true
		}
	}

}
