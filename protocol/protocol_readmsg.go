package protocol

type ReadMsgPacket struct {
	Encryption   string
	IMEI         string
	SerialNumber string
	MsgId        string
}

func (p *ReadMsgPacket) Serialize() []byte {
	var result string
	result += p.Encryption
	result += RD
	result += SEP
	result += p.IMEI
	result += SEP
	result += p.SerialNumber
	result += SEP
	result += p.MsgId
	result += SEP
	result += ENDFLAG

	return []byte(result)
}

func ParseReadMsg(buffer []byte) *ReadMsgPacket {
	encryption, values := ParseCommon(buffer)
	return &ReadMsgPacket{
		Encryption:   encryption,
		IMEI:         values[1],
		SerialNumber: values[2],
		MsgId:        values[3],
	}
}
func SetReadMsg(encryption string, imei string, serialnum string, msgId string) *ReadMsgPacket {
	return &ReadMsgPacket{
		Encryption:   encryption,
		IMEI:         imei,
		SerialNumber: serialnum,
		MsgId:        msgId,
	}

}
