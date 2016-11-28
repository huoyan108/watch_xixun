package protocol

type ChargePacket struct {
	Encryption   string
	IMEI         string
	SerialNumber string
	WarnStyle    string
}

func (p *ChargePacket) Serialize() []byte {
	var result string
	result += p.Encryption
	result += CG
	result += SEP
	result += p.IMEI
	result += SEP
	result += p.SerialNumber
	result += SEP
	result += ENDFLAG

	return []byte(result)
}

func ParseCharge(buffer []byte) *ChargePacket {
	encryption, values := ParseCommon(buffer)
	return &ChargePacket{
		Encryption:   encryption,
		IMEI:         values[1],
		SerialNumber: values[2],
	}
}
func SetCharge(encryption string, imei string, serialnum string) *ChargePacket {
	return &ChargePacket{
		Encryption:   encryption,
		IMEI:         imei,
		SerialNumber: serialnum,
	}

}
