package protocol

import ()

type ControlPacket struct {
	Encryption   string
	IMEI         string
	SerialNumber string
	Action       string
	//Style        Report.Command_CommandType
	Style string
}

func (p *ControlPacket) Serialize() []byte {
	var result string
	result += p.Encryption
	result += p.IMEI
	result += SEP
	result += p.SerialNumber
	result += ",123456cmd,"
	result += p.Action
	result += ENDFLAG

	return []byte(result)
}

func SetControlCmd(encryption string, imei string, serialnum string, style string, action string) *ControlPacket {
	return &ControlPacket{
		Encryption:   encryption,
		IMEI:         imei,
		SerialNumber: serialnum,
		Action:       action,
		Style:        style,
	}
}
func ParseControlCmdRt(buffer []byte) *ControlPacket {
	encryption, values := ParseCommon(buffer)
	style := string(values[3][0:2])
	return &ControlPacket{
		Encryption:   encryption,
		IMEI:         values[0],
		SerialNumber: values[1],
		Action:       values[3],
		Style:        style,
	}
}
