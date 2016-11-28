package protocol

import (
	"strconv"
)

const (
	FEEDBACK_CMDID string = "ac"
)

type PosUpPacket struct {
	Encryption   string
	IMEI         string
	SerialNumber string
	LocationTime string
	Longitude    string
	Latitude     string
	GPSFlag      string
	Wifi         string
	WifiCount    int
	Jzs          string
	JzCount      int
	Battery      string
	Charge       int
}

func (p *PosUpPacket) Serialize() []byte {
	var result string
	result += p.Encryption
	result += FEEDBACK_CMDID
	result += SEP
	result += p.IMEI
	result += SEP
	result += p.SerialNumber
	result += SEP
	if p.GPSFlag != "" {
		result += "1,"
	} else {
		result += "0,"
	}
	result += p.LocationTime
	result += SEP
	result += ENDFLAG

	return []byte(result)
}

func ReWriteWifi(wifi string) string {
	ret := wifi[0:2] + ":"
	ret += wifi[2:4] + ":"
	ret += wifi[4:6] + ":"
	ret += wifi[6:8] + ":"
	ret += wifi[8:10] + ":"
	ret += wifi[10:12]

	return ret
}

func ParseJz(jzs []string) string {
	item_count := len(jzs)
	var ret string = ""
	var j int = 2

	ret += jzs[0] + ","
	ret += jzs[1] + "|"
	for i := 0; i < item_count/3; i++ {
		ret += jzs[j] + ","
		ret += jzs[j+1] + ","
		ret += jzs[j+2] + "|"
		j += 3
	}

	return ret[0 : len(ret)-1]
}

func ParseWifi(wifis []string) string {
	item_count := len(wifis)
	var ret string = ""
	var j int = 0
	for i := 0; i < item_count/2; i++ {
		ret += ReWriteWifi(wifis[j]) + ","
		ret += wifis[j+1] + ","
		ret += "TP_LINK|"
		j += 2
	}

	return ret[0 : len(ret)-1]
}

func ParsePosUp(buffer []byte) *PosUpPacket {
	encryption, values := ParseCommon(buffer)

	lat := values[5][1:]
	long := values[6][1:]
	wifi_count := values[12]
	count, _ := strconv.Atoi(wifi_count)
	var wifis string = ""
	if count > 0 {
		wifis = ParseWifi(values[13 : 13+count*2])
	}
	var pos int
	pos = 13 + count*2
	jz_count := values[pos+1]
	jzCount, _ := strconv.Atoi(jz_count)
	pos += 2
	var jzs string = ""
	if jzCount > 0 {
		jzs = ParseJz(values[pos : pos+jzCount*3+2])
	}
	pos = pos + jzCount*3 + 2

	battery := values[pos+2][1:]
	charge, _ := strconv.Atoi(values[pos+3][1:])

	return &PosUpPacket{
		Encryption:   encryption,
		IMEI:         values[1],
		SerialNumber: values[2],
		LocationTime: values[3],
		Latitude:     lat,
		Longitude:    long,
		GPSFlag:      values[len(values)-7],
		Wifi:         wifis,
		WifiCount:    count,
		Jzs:          jzs,
		JzCount:      jzCount,
		Battery:      battery,
		Charge:       charge,
	}
}
